package server

import (
	"fmt"
	"log"
	"net"
    "time"
	genericmessage "raft/internal/genericMessage"
	"raft/internal/node"
	"raft/internal/raft_log"
	"raft/internal/raftstate"
	state "raft/internal/raftstate"
	"raft/internal/rpcs"
	"raft/internal/rpcs/AppendEntryRpc"
	"raft/internal/rpcs/RequestVoteRPC"
	"raft/internal/rpcs/UpdateNode"
	"raft/internal/rpcs/redirection"
	p "raft/pkg/raft-rpcProtobuf-messages/rpcEncoding/out/protobuf"
	"reflect"
	"strings"
	"sync"
)

type pairMex struct{
    payload rpcs.Rpc
    sender string
    workdone chan int
}

type server struct {
    wg             sync.WaitGroup
	_state         state.State
	messageChannel chan pairMex
	unstableNodes   *sync.Map
    numNodes uint
    clientNodes    *sync.Map
	listener       net.Listener
}

func (s *server) Start() {
    log.Println("Start accepting connections")

    s.wg.Add(1)
    go s.acceptIncomingConn()
    s.wg.Add(1)
    go s.run()

    s.wg.Wait()
    log.Println("run finished")
}

//utility
func (this *server) connectToNodes(serversIp []string, port string) ([]string,error){
    var failedConn []string = make([]string, 0)
    var err error

	for i := 0; i < len(serversIp)-1; i++ {
		var new_node node.Node
        var nodeConn net.Conn
        var erroConn error

        nodeConn,erroConn = net.Dial("tcp",serversIp[i]+":"+port)
        if erroConn != nil {
            log.Println("Failed to connect to node: ", serversIp[i])
            failedConn = append(failedConn, serversIp[i])
            err = erroConn
            continue
        }
        new_node = node.NewNode(serversIp[i], port, nodeConn,this._state.GetStatePool())
        log.Printf("connected to new node, storing it: %v\n", new_node.GetIp())
        (*this).unstableNodes.Store(new_node.GetIp(), new_node)
        (*this).numNodes++
        go (*this).internalNodeConnection(new_node)
	}

    return failedConn,err
}

func (s *server) acceptIncomingConn() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Println("Failed on accept: ", err)
			continue
		}

		tcpAddr, ok := conn.RemoteAddr().(*net.TCPAddr)
		if !ok {
			fmt.Println("Connection is not using TCP")
			continue
		}

		var newConncetionIp string = tcpAddr.IP.String()
		var newConncetionPort string = string(rune(tcpAddr.Port))
        var found bool
		_, found = s.unstableNodes.Load(newConncetionIp)
        log.Printf("new connec: %v\n", newConncetionIp)

        if found {
            log.Panicln("reconnecting to an already known node: ",newConncetionIp)
        }

        log.Printf("node with ip %v not found", newConncetionIp)
        var newNode node.Node = node.NewNode(newConncetionIp, newConncetionPort,conn, s._state.GetStatePool())
        s.unstableNodes.Store(newNode.GetIp(),newNode)
        s.numNodes++
        go s.handleConnection(newNode)
	}
}

func (s* server) handleConnection(workingNode node.Node){
    if strings.Contains(workingNode.GetIp(), "10.0.0") {
        s.internalNodeConnection(workingNode)
        return
    }
    s.externalAgentConnection(workingNode)
}

func (s *server) externalAgentConnection(agent node.Node){
    var leaderIp = s._state.GetLeaderIp(raftstate.PUB)
    var rawMex []byte
    var err error
    var resp rpcs.Rpc 
    var inputMex rpcs.Rpc 
    var clientReq pairMex = pairMex{}

    //INFO: send to the client who you think is the leader, (blank means you are not in the conf)
    resp = Redirection.NewredirectionRPC(leaderIp)
    s.encodeAndSend(resp,agent)

    for{
        rawMex,err = agent.Recv()
        if err != nil || rawMex == nil {
            fmt.Printf("error in reading from node %v with error %v\n",agent.GetIp(), rawMex)
            break
        }
        inputMex,err = genericmessage.Decode(rawMex)
        if err != nil{
            log.Println(err)
            break
        }
        clientReq.payload = inputMex
        clientReq.sender = agent.GetIp()
        clientReq.workdone = make(chan int)
        s.messageChannel <- clientReq
        log.Println("waiting agent completion on chann: ",clientReq.workdone)
        <- clientReq.workdone
    }

    log.Println("Done serving client: ", agent.GetIp())
    agent.CloseConnection()
    s.unstableNodes.Delete(agent.GetIp())
    s.numNodes--
}

func (s *server) internalNodeConnection(workingNode node.Node) {
    var nodeIp = workingNode.GetIp()
    var message []byte
    var errMes error
    var rpcMex rpcs.Rpc
    var nodeReq pairMex = pairMex{}


    for{
        message, errMes = workingNode.Recv()
        if errMes != nil {
            log.Printf("error in reading from node %v with error %v\n",nodeIp, errMes)
            break
        }else if message != nil {
            rpcMex,errMes = genericmessage.Decode(message)
            if errMes != nil {
                log.Println(errMes)
                continue
            }
            nodeReq.payload = rpcMex
            nodeReq.sender = workingNode.GetIp()
            nodeReq.workdone = make(chan int)
            s.messageChannel <- nodeReq
            <- nodeReq.workdone
        }
    }
    workingNode.CloseConnection()
    s.unstableNodes.Delete(nodeIp)
    s.numNodes--
}

func (s *server) run() {
    var mess pairMex
    var leaderCommitEntry int64
    var nodeToUpdate raftstate.NewNodeToUpdateInfo
    var timeoutElection,err = s._state.GetTimeoutNotifycationChan(raftstate.TIMER_ELECTION)
    var entryToPropagateChann = *s._state.GetLeaderEntryChannel()

    if err != nil{
        log.Panicln(err)
    }


    for {
        select {
        case mess = <-s.messageChannel:
            go s.newMessageReceived(mess)
        case <- timeoutElection:
            go s.startNewElection()
        case leaderCommitEntry = <- entryToPropagateChann:
            go s.newEntryToCommit(leaderCommitEntry)
        case nodeToUpdate = <- s._state.GetNewNodeToUpdate():
            for _, v := range nodeToUpdate.NodeList {
                go s.updateNewNode(v,nodeToUpdate.MatchToArrive)
            }
        }
    }
}

func (s *server) updateNewNode(nodeIp string, matchIdx uint64) error{
    var found = false
    var value any
    var newNode node.Node
    var waitTimer *time.Timer = nil
    var commitConf p.LogEntry
    var hearthBit rpcs.Rpc
        
    if nodeIp == s._state.GetMyIp(raftstate.PRI){
        return nil
    }

    value,found = s.unstableNodes.Load(nodeIp)



    if !found{ //HACK: polling wait
        waitTimer = time.NewTimer(time.Duration(5 * time.Second))
        <- waitTimer.C
        return s.updateNewNode(nodeIp,matchIdx)
    }

    newNode = value.(node.Node)

    
    log.Println("updating node: ", nodeIp)
    s.encodeAndSend(UpdateNode.ChangeVoteRightNode(false),newNode)

    hearthBit = s.nodeAppendEntryPayload(newNode,nil)
    s.encodeAndSend(hearthBit,newNode)
    
    for newNode.GetMatchIndex() < int(matchIdx){
        //HACK: polling wait
        waitTimer = time.NewTimer(time.Duration(5 * time.Second))
        <- waitTimer.C
    }
    s.encodeAndSend(UpdateNode.ChangeVoteRightNode(true),newNode)
    newNode.NodeUpdated()
    log.Println("done updating node: ", nodeIp)
    
    commitConf = p.LogEntry{
        Term: s._state.GetTerm(),
        OpType: p.Operation_COMMIT_CONFIG_ADD,
        Description: "committin conf adding " + nodeIp,
        Payload: []byte(nodeIp),
    }

    var entryWrapper *raft_log.LogInstance = s._state.NewLogInstance(&commitConf,func() {})
    s._state.AppendEntries([]*raft_log.LogInstance{entryWrapper})

    return nil
}

func (s *server) newEntryToCommit(leaderCommitEntry int64){
    //TODO: check that at least the majority of the followers has a commit index
    // >= than this, if not send him an AppendEntryRpc with the entry,
    //Search only through stable nodes because may be possible that a new node is still 
    //updating while you check his commitIndex
    var err error
    var entryToCommit *raft_log.LogInstance

    entryToCommit,err = s._state.GetEntriAt(leaderCommitEntry)
    log.Println("new log entry to propagate: ", entryToCommit)
    if err != nil {
        log.Panic("invalid index entry: ", leaderCommitEntry)
    }

    s.applyOnFollowers(func(n node.Node) {
        log.Println("propagate log entry to: ",n.GetIp())
        var AppendEntry rpcs.Rpc
        var rawMex []byte

        if err!=nil {
            log.Panicln(err)
        }

        if !n.Updated(){
            log.Println("node not updated: ", n.GetIp())
            return 
        }

        AppendEntry = s.nodeAppendEntryPayload(n,[]raft_log.LogInstance{*entryToCommit})

        rawMex,err = genericmessage.Encode(AppendEntry)
        if err != nil {
            log.Panicln("error encoding AppendEntry: ",AppendEntry.ToString())
        }

        log.Printf("sending appendEntry to: %v, %v",n.GetIp(), AppendEntry.ToString())
        n.Send(rawMex)
    })
}

func (s *server) newMessageReceived(mess pairMex){
            var rpcCall rpcs.Rpc
            var sender string = mess.sender
            var oldRole raftstate.Role
            var resp rpcs.Rpc
            var byEnc []byte
            var errEn error
            var f any
            var ok bool
            var senderNode node.Node

            f, ok = s.unstableNodes.Load(sender)
            if !ok {
                log.Panicf("Node %s not found for mex:%v\n", sender, mess.payload.ToString())
            }

            senderNode = f.(node.Node)
            oldRole = s._state.GetRole()
            rpcCall = mess.payload
            resp = rpcCall.Execute(s._state, senderNode)
            log.Println("finih executing rpc: ",rpcCall.ToString())

            if resp != nil {
                log.Println("sending resp to caller RPC: ", resp.ToString())
                byEnc, errEn = genericmessage.Encode(resp)
                if errEn != nil{
                    log.Panicln("error encoding this rpc: ", resp.ToString())
                }
                senderNode.Send(byEnc)
            }

            if s._state.GetRole() == raftstate.LEADER && oldRole != state.LEADER {
                log.Printf("init commonMatch pool with lastLogIndex: %v\n",s._state.LastLogIndex())
                s._state.GetStatePool().InitCommonMatch(s._state.LastLogIndex()+1)
                s.applyOnFollowers(func(n node.Node) {
                    n.NodeUpdated()
                })
                go s.leaderHearthBit()
            }
            mess.workdone <- 1
}

//utility
func (s *server) startNewElection(){
    log.Println("started new election");

    s._state.IncrementTerm()
    s._state.IncreaseSupporters()

    if s._state.GetNumberNodesInCurrentConf() == 1 {
        if s.numNodes == 0 {
            log.Println("became leader: ",s._state.GetRole())
            s._state.SetRole(raftstate.LEADER)
            go s.leaderHearthBit()
            return
        }
    }

    s.applyOnFollowers(func(n node.Node) {
            var voteRequest rpcs.Rpc 
            var raw_mex []byte
            var err error
            var lastLogIndex int
            var lastLogTerm uint

            lastLogIndex = s._state.LastLogIndex()
            lastLogTerm = s._state.LastLogTerm()
            voteRequest = RequestVoteRPC.NewRequestVoteRPC(
                s._state.GetTerm(),
                s._state.GetMyIp(raftstate.PRI),
                int64(lastLogIndex),
                uint64(lastLogTerm))

                raw_mex,err = genericmessage.Encode(voteRequest)
                if err != nil {
                    log.Panicln("error in Encoding this voteRequest: ",(voteRequest).ToString())
                }
                log.Printf("sending election request %v to %v\n", voteRequest.ToString(),n.GetIp())
                n.Send(raw_mex)
            })
    s._state.RestartTimeout(raftstate.TIMER_ELECTION)
}

func (s *server) leaderHearthBit(){
    var timerHearthbit,err = s._state.GetTimeoutNotifycationChan(raftstate.TIMER_HEARTHBIT)
    if err != nil {
        log.Panicln(err)
    }

    for {
        log.Println("sending hearthbit")
        <- timerHearthbit

        s.applyOnFollowers(func(n node.Node) {
            var hearthBit rpcs.Rpc = s.nodeAppendEntryPayload(n,nil)
            s.encodeAndSend(hearthBit,n)
        })
        s._state.RestartTimeout(raftstate.TIMER_ELECTION)
    }
}

func (s* server) applyOnFollowers(fn func(n node.Node)){
    var currentConf []string = s._state.GetConfig()
    for _, v := range currentConf {
        var nNode node.Node 
        var value any
        var found bool

        if v == s._state.GetMyIp(raftstate.PRI){
            continue
        }

        value, found= s.unstableNodes.Load(v)
        if !found {
            log.Println("node in conf not saved in unstablequeue ", v)
            continue
        }
        nNode = value.(node.Node)
        if !found {
            var s = reflect.TypeOf(value)
            log.Panicln("failed conversion type node, type is: ", s)
        }
        go fn(nNode)
    }
}

func (s *server) nodeAppendEntryPayload(n node.Node, toAppend []raft_log.LogInstance) rpcs.Rpc{
    var nodeNextIndex = n.GetNextIndex()
    var prevLogTerm uint64 =0
    var hearthBit rpcs.Rpc
    var logDataInstance []raft_log.LogInstance= s._state.GetCommittedEntriesRange(int(nodeNextIndex))
    var logEntryPayload []*p.LogEntry = nil
    var prevLogEntry *raft_log.LogInstance
    var err error
    var prevLogIndex int64 = -1

    if toAppend != nil{
        logDataInstance = append(logDataInstance, toAppend...)
    }

    for _, v := range logDataInstance {
        logEntryPayload = append(logEntryPayload, v.Entry)
    }

    if nodeNextIndex <= s._state.LastLogIndex() || toAppend != nil {

        if nodeNextIndex > 0 {
            prevLogEntry,err = s._state.GetEntriAt(int64(nodeNextIndex)-1)
            if err != nil{
                log.Panicln("error retrieving entry at index :", nodeNextIndex-1)
            }
            prevLogTerm = prevLogEntry.Entry.Term
            prevLogIndex = int64(nodeNextIndex) -1
        }
        
        hearthBit = AppendEntryRpc.NewAppendEntryRPC(
            s._state, prevLogIndex, prevLogTerm,logEntryPayload)

    }else {
        hearthBit = AppendEntryRpc.GenerateHearthbeat(s._state)  
        log.Printf("sending hearthbit: %v\n", hearthBit.ToString())
    }
    return hearthBit
}

func (s *server) encodeAndSend(rpcMex rpcs.Rpc, n node.Node){
    var err error
    var rawMex []byte

    rawMex,err = genericmessage.Encode(rpcMex)
    if err != nil {
        log.Panicln("error encoding rpc: ", rpcMex.ToString())
    }

    err = n.Send(rawMex)

    if err != nil {
        log.Println("error sending rpcRawMex: ", err)
    }

}
