package server

import (
	"fmt"
	"log"
	"net"
	genericmessage "raft/internal/genericMessage"
	"raft/internal/node"
	"raft/internal/raftstate"
	state "raft/internal/raftstate"
	"raft/internal/rpcs"
	"raft/internal/rpcs/AppendEntryRpc"
	"raft/internal/rpcs/ClientReq"
	"raft/internal/rpcs/RequestVoteRPC"
	"raft/internal/rpcs/UpdateNode"
	p "raft/pkg/raft-rpcProtobuf-messages/rpcEncoding/out/protobuf"
	"reflect"
	"strings"
	"sync"

	"google.golang.org/protobuf/proto"
)

type pairMex struct{
    payload *rpcs.Rpc
    sender string
}

type server struct {
    wg             sync.WaitGroup
	_state         state.State
	messageChannel chan pairMex
	unstableNodes   *sync.Map
    clientNodes    *sync.Map
	listener       net.Listener
}

func (s *server) Start() {
    log.Println("Start accepting connections")
    s._state.StartElectionTimeout()

    s.wg.Add(1)
    go s.acceptIncomingConn()
    go func(){
        s.wg.Add(1)
        defer s.wg.Done()
        s.run()
    }()

    s.wg.Wait()
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
        go (*this).handleResponseSingleNode(new_node)
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
            log.Println("reconnecting to an already known node: ",newConncetionIp)
        }

        log.Printf("node with ip %v not found", newConncetionIp)
        var new_node node.Node = node.NewNode(newConncetionIp, newConncetionPort,conn, s._state.GetStatePool())
        s.unstableNodes.Store(new_node.GetIp(),new_node)
        go func ()  {
            s.wg.Add(1)
            defer s.wg.Done()
            s.handleConnection(new_node)
        }()
	}
}

func (s* server) handleConnection(workingNode node.Node){
    var nodeIp string = workingNode.GetIp()

    if strings.Contains(nodeIp, "10.0.0") {
        s.handleResponseSingleNode(workingNode)
        return
    }
    s.handleNewClientConnection(workingNode)
}

func (s* server) handleNewClientConnection(client node.Node){
    log.Println("new client request to cluster")
    if s._state.Leader(){
        var clientReq rpcs.Rpc = &ClientReq.ClientReq{}
        var ok = "ok"
        var leaderIp p.PublicIp = p.PublicIp{IP: ok,}
        var mex,err = proto.Marshal(&leaderIp)
        var resp rpcs.Rpc = nil

        if err != nil {
            log.Panicln("error encoding confirmation leader public ip for client:",err)
        }
    
        if err != nil {
            log.Panicln(err)
        }

        client.Send(mex)

        mex,err = client.Recv()
        if err != nil {
            fmt.Printf("error in reading from node %v with error %v\n",client.GetIp(), err)
            client.CloseConnection()
            return
        }

        err = clientReq.Decode(mex)
        if err != nil {
            fmt.Printf("error in decoding client request from node %v with error %v\n",client.GetIp(), err)
            client.CloseConnection()
            return
        }
        log.Println("managing client Request: ", clientReq.ToString())
        resp = clientReq.Execute(s._state,client)
        
        if resp != nil{
            s.encodeAndSend(resp,client)
        }

    }else{
        var leaderIp p.PublicIp = p.PublicIp{IP: s._state.GetLeaderIpPublic(),}
        var mex,err = proto.Marshal(&leaderIp)
        log.Printf("sending public ip of leader: %v\n", s._state.GetLeaderIpPublic())

        if err != nil {
            log.Panicln("error encoding leader public ip for client:",err)
        }

        client.Send(mex)
    }
    client.CloseConnection()
}

func (s *server) handleResponseSingleNode(workingNode node.Node) {
    var nodeIp = workingNode.GetIp()
    var message []byte
    var errMes error
    var notifyChan chan int
    var newConfDelete p.LogEntry = p.LogEntry{
        OpType: p.Operation_JOIN_CONF_DEL,
        Term: s._state.GetTerm(),
        Payload: []byte(nodeIp),
        Description: "removed node " + nodeIp + " to configuration: ",
    }

    if s._state.Leader() {
        log.Println("i'm leader, joining conf")
        go s.joinConf(workingNode)
    }

    for{
        message, errMes = workingNode.Recv()
        if errMes != nil {
            fmt.Printf("error in reading from node %v with error %v",nodeIp, errMes)
            if !s._state.Leader() {
                s._state.StartElectionTimeout()
            }

            if s._state.Leader() || s._state.GetNumberNodesInCurrentConf() == 2{
                var chans =  s._state.AppendEntries([]*p.LogEntry{&newConfDelete})
                notifyChan = chans[len(chans)-1]
            }

            <- notifyChan
            log.Println("safe to remove")
            workingNode.CloseConnection()
            s.unstableNodes.Delete(nodeIp); 
            s._state.GetStatePool().RemNode(nodeIp) 
            break
        }
        if message != nil {
            s.messageChannel <- pairMex{genericmessage.Decode(message),workingNode.GetIp()}
        }
    }

}

func (s *server) joinConf(workingNode node.Node){
    var nodeIp = workingNode.GetIp()
    var chans []chan int
    var notifyChan chan int
    var newConfEntry p.LogEntry = p.LogEntry{
        OpType: p.Operation_JOIN_CONF_ADD,
        Term: s._state.GetTerm(),
        Payload: []byte(nodeIp),
        Description: "added new node " + nodeIp + " to configuration: ",
    }
    var commitConf p.LogEntry = p.LogEntry{
        OpType: p.Operation_COMMIT_CONFIG,
        Term: s._state.GetTerm(),
        Payload: []byte(nodeIp),
        Description: "committing config add of node " + nodeIp,
    }

    var commitedEntries []*p.LogEntry = s._state.GetCommittedEntries()
    var appendEntryRpc rpcs.Rpc = s.nodeAppendEntryPayload(workingNode,nil)

    chans = s._state.AppendEntries([]*p.LogEntry{&newConfEntry})
    notifyChan = chans[len(chans)-1]

    log.Println("updating node: ",workingNode.GetIp())
    s.encodeAndSend(UpdateNode.ChangeVoteRightNode(false),workingNode)
    
    log.Println("sending appendEntry mex udpated: ", appendEntryRpc.ToString())
    s.encodeAndSend(appendEntryRpc,workingNode)

    log.Println("waiting that matchIndex is: ", len(commitedEntries)-1)
    for  workingNode.GetMatchIndex() < len(commitedEntries)-1 {
        //HACK: WAIT POLLING
    }
    s.encodeAndSend(UpdateNode.ChangeVoteRightNode(true),workingNode)
    workingNode.NodeUpdated()
    log.Println("done updating node: ",workingNode.GetIp())

    <- notifyChan
    log.Println("commit config")
    s._state.AppendEntries([]*p.LogEntry{&commitConf})
}

func (s *server) run() {
    for {
        var mess pairMex
        var leaderCommitEntry int64

        select {
        case mess = <-s.messageChannel:
            var rpcCall rpcs.Rpc
            var sender string = mess.sender
            var oldRole raftstate.Role
            var resp rpcs.Rpc
            var byEnc []byte
            var errEn error
            var f any
            var ok bool
            var senderNode node.Node
            var newConf []string
            var failedConn []string

            f, ok = s.unstableNodes.Load(sender)
            if !ok {
                log.Printf("Node %s not found for mex:%v\n", sender, (*mess.payload).ToString())
                continue
            }

            senderNode = f.(node.Node)
            oldRole = s._state.GetRole()
            rpcCall = *mess.payload
            resp = rpcCall.Execute(s._state, senderNode)

            if s._state.ConfChanged() {
                log.Printf("configuration changed, adding the new nodes\n")
                newConf = s._state.GetConfig()
                for _, v := range newConf {
                    var _,found = s.unstableNodes.Load(v)
                    if !found {
                        failedConn,errEn = s.connectToNodes([]string{v},"8080") //WARN: hard encoding port
                        if errEn != nil {
                            for _, v := range failedConn {
                                log.Println("failed connecting to this node: " + v)
                            }
                            log.Panic("devel debug") //WARN: to remove in the future
                        }
                    }
                }
            }

            if resp != nil {
                byEnc, errEn = genericmessage.Encode(&resp)
                if errEn != nil{
                    log.Panicln("error encoding this rpc: ", resp.ToString())
                }
                senderNode.Send(byEnc)
            }

            if s._state.Leader() && oldRole != state.LEADER {
                log.Printf("init commonMatch pool with lastLogIndex: %v\n",s._state.LastLogIndex())
                s._state.GetStatePool().InitCommonMatch(s._state.LastLogIndex()+1)
                s.applyOnFollowers(func(n node.Node) {
                    n.NodeUpdated()
                })
                go s.leaderHearthBit()
            }
        case <-s._state.ElectionTimeout().C:
            if !s._state.Leader() {
                s.startNewElection()
            }
        case leaderCommitEntry = <-(*s._state.GetLeaderEntryChannel()):
            //TODO: check that at least the majority of the followers has a commit index
            // >= than this, if not send him an AppendEntryRpc with the entry,
            //Search only through stable nodes because may be possible that a new node is still 
            //updating while you check his commitIndex
            var err error
            var entryToCommit *p.LogEntry 

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
                    return 
                }

                AppendEntry = s.nodeAppendEntryPayload(n,[]*p.LogEntry{entryToCommit})

                rawMex,err = genericmessage.Encode(&AppendEntry)
                if err != nil {
                    log.Panicln("error encoding AppendEntry: ",AppendEntry.ToString())
                }

                log.Printf("sending appendEntry to: %v, %v",n.GetIp(), AppendEntry.ToString())
                n.Send(rawMex)
            })
        }
    }
}

//utility
func (s *server) startNewElection(){
    log.Println("started new election");

    s._state.IncrementTerm()
    s._state.IncreaseSupporters()

    if s._state.GetNumberNodesInCurrentConf() == 1 {
        log.Println("became leader: ",s._state.GetRole())
        s._state.SetRole(raftstate.LEADER)
        s._state.SetLeaderIpPrivate(s._state.GetIdPrivate())
        s._state.SetLeaderIpPublic(s._state.GetIdPublic())
        s._state.ResetElection()
        go s.leaderHearthBit()
        return
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
                s._state.GetIdPrivate(),
                int64(lastLogIndex),
                uint64(lastLogTerm))

                raw_mex,err = genericmessage.Encode(&voteRequest)
                if err != nil {
                    log.Panicln("error in Encoding this voteRequest: ",(voteRequest).ToString())
                }
                log.Printf("sending election request %v to %v\n", voteRequest.ToString(),n.GetIp())
                n.Send(raw_mex)
            })
}

func (s *server) leaderHearthBit(){
    //log.Println("start sending hearthbit")
    for s._state.Leader(){
        <- s._state.HeartbeatTimeout().C

        log.Println("start broadcast")
        s.applyOnFollowers(func(n node.Node) {
            var hearthBit rpcs.Rpc = s.nodeAppendEntryPayload(n,nil)
            s.encodeAndSend(hearthBit,n)
        })
        log.Println("end broadcast")
        s._state.StartHearthbeatTimeout()
    }
    s._state.GetStatePool().InitCommonMatch(s._state.LastLogIndex())
    log.Println("no longer LEADER, stop sending hearthbit")
}

func (s* server) applyOnFollowers(fn func(n node.Node)){
    var currentConf []string = s._state.GetConfig()
    for _, v := range currentConf {
        var nNode node.Node 
        var value any
        var found bool

        if v == s._state.GetIdPrivate() {
            continue
        }

        value, found= s.unstableNodes.Load(v)
        if !found {
            log.Panicln("node in conf not saved in unstablequeue ", v)
        }
        nNode = value.(node.Node)
        if !found {
            var s = reflect.TypeOf(value)
            log.Panicln("failed conversion type node, type is: ", s)
        }
        go fn(nNode)
    }
}

func (s *server) nodeAppendEntryPayload(n node.Node, toAppend []*p.LogEntry) rpcs.Rpc{
    var nodeNextIndex = n.GetNextIndex()
    var prevLogTerm uint64 =0
    var hearthBit rpcs.Rpc
    var entryPayload []*p.LogEntry = s._state.GetCommittedEntriesRange(int(nodeNextIndex))
    var prevLogEntry *p.LogEntry
    var err error
    var prevLogIndex int64 = -1

    if toAppend != nil{
        entryPayload = append(entryPayload, toAppend...)
    }

    if nodeNextIndex <= s._state.LastLogIndex() || toAppend != nil {

        if nodeNextIndex > 0 {
            prevLogEntry,err = s._state.GetEntriAt(int64(nodeNextIndex)-1)
            if err != nil{
                log.Panicln("error retrieving entry at index :", nodeNextIndex-1)
            }
            prevLogTerm = prevLogEntry.Term
            prevLogIndex = int64(nodeNextIndex) -1
        }
        
        hearthBit = AppendEntryRpc.NewAppendEntryRPC(
            s._state, prevLogIndex, prevLogTerm,entryPayload)

    }else {
        hearthBit = AppendEntryRpc.GenerateHearthbeat(s._state)  
        log.Printf("sending hearthbit: %v\n", hearthBit.ToString())
    }
    return hearthBit
}

func (s *server) encodeAndSend(rpcMex rpcs.Rpc, n node.Node){
    var err error
    var rawMex []byte

    rawMex,err = genericmessage.Encode(&rpcMex)
    if err != nil {
        log.Panicln("error encoding rpc: ", rpcMex.ToString())
    }

    err = n.Send(rawMex)

    if err != nil {
        log.Println("error sending rpcRawMex: ", err)
    }

}
