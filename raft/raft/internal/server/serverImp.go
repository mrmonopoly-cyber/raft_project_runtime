package server

import (
	"fmt"
	"log"
	"net"
	genericmessage "raft/internal/genericMessage"
	"raft/internal/node"
	"raft/internal/node/nodeState"
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
        var nodeState nodeState.VolatileNodeState

        if err != nil {
            log.Panicln("error encoding confirmation leader public ip for client:",err)
        }
    
        nodeState,err = client.GetNodeState()
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
        clientReq.Execute(&s._state,&nodeState)
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
            workingNode.CloseConnection()
            s.unstableNodes.Delete(nodeIp);
            if s._state.Leader(){
                s._state.AppendEntries([]*p.LogEntry{&newConfDelete})
            }
            break
        }
        if message != nil {
            s.messageChannel <- 
            pairMex{genericmessage.Decode(message),workingNode.GetIp()}
        }
    }

}

func (s *server) joinConf(workingNode node.Node){
    var nodeIp = workingNode.GetIp()

    log.Println("debug, joinConf : ,", s._state.GetConfig())

    var newConfEntry p.LogEntry = p.LogEntry{
        OpType: p.Operation_JOIN_CONF_ADD,
        Term: s._state.GetTerm(),
        Payload: []byte(nodeIp),
        Description: "added new node " + nodeIp + " to configuration: ",
    }

    s._state.AppendEntries([]*p.LogEntry{&newConfEntry})
    s.updateNewNode(workingNode)              
}

func (s *server) updateNewNode(workingNode node.Node){
    var volatileState nodeState.VolatileNodeState
    var err error

    volatileState,err = workingNode.GetNodeState()
    if err != nil {
        log.Panicln(err)
    }
    log.Printf("updating node %v\n", workingNode.GetIp())

    log.Printf("\nupdating, list of entries to send: %v\n\n",s._state.GetCommittedEntries())

    for  i,e := range s._state.GetCommittedEntries() {
        log.Printf("sending update mex to %v with data %v\n",workingNode.GetIp(), e)
        err = s.generateUpdateRequest(workingNode,false,e)
        if err != nil {
            log.Printf("error generata UpdateRequest : %v\n", err)
            return 
        }
        for  volatileState.GetMatchIndex() < i {
            //WARN: WAIT
        }
    }
    err = s.generateUpdateRequest(workingNode,true,nil)
    if err != nil {
        log.Printf("error generata UpdateRequest : %v\n", err)
        return 
    }
    log.Printf("node %v updated\n",workingNode.GetIp())
    volatileState.NodeUpdated()
}

func (this *server) generateUpdateRequest(workingNode node.Node, voting bool, entry *p.LogEntry) error{
    var updateReq rpcs.Rpc 
    var mex []byte
    var err error
    
    updateReq = UpdateNode.NewUpdateNodeRPC(voting, entry)
    mex,err = genericmessage.Encode(&updateReq)
    if err != nil {
        log.Panic("error encoding UpdateNode rpc")
    }
    return workingNode.Send(mex)
}

func (s *server) sendAll(rpc *rpcs.Rpc){
   log.Println("start broadcast")
   s.applyOnFollowers(func(n node.Node) {
       var raw_mex []byte
       var err error

       raw_mex,err = genericmessage.Encode(rpc)
       if err != nil {
           log.Panicln("error in Encoding this rpc: ",(*rpc).ToString())
       }
       log.Printf("sending to %v with key %v\n", n.GetIp(),n.GetIp())
       err = n.Send(raw_mex)
       if err != nil{
           log.Panicln("error sending message to node ",err)
       }


   })
   log.Println("end broadcast")
}

func (s *server) run() {
    for {
        var mess pairMex
        var leaderCommitEntry int64

        select {
        case mess = <-s.messageChannel:
            var rpcCall *rpcs.Rpc
            var sender string = mess.sender
            var oldRole raftstate.Role
            var resp *rpcs.Rpc
            var byEnc []byte
            var errEn error
            var f any
            var ok bool
            var senderState nodeState.VolatileNodeState
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
            rpcCall = mess.payload
            senderState,errEn = senderNode.GetNodeState()
            if errEn != nil {
                log.Panicln(errEn)
            }
            resp = (*rpcCall).Execute(&s._state, &senderState)

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
                //          log.Println("reponse to send to: ", sender)

                //log.Println("sending mex to: ",sender)
                byEnc, errEn = genericmessage.Encode(resp)
                if errEn != nil{
                    log.Panicln("error encoding this rpc: ", (*resp).ToString())
                }
                senderNode.Send(byEnc)
            }

            if s._state.Leader() && oldRole != state.LEADER {
                s.setVolState() 
                go s.leaderHearthBit()
            }
            //         log.Println("rpc processed")
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
            var prevEntryTerm uint64 = 0
            var leaderCommit = s._state.GetCommitIndex()
            var numStableNodes uint = 0

            log.Println("new log entry to propagate")
            if leaderCommitEntry-1 >= 0{
                var prevEntry,err = s._state.GetEntriAt(leaderCommitEntry-1)
                if err != nil {
                    log.Panic("invalid index entry: ", leaderCommitEntry)
                }
                prevEntryTerm = prevEntry.Term
            }

            entryToCommit,err = s._state.GetEntriAt(leaderCommitEntry)
            if err != nil {
                log.Panic("invalid index entry: ", leaderCommitEntry)
            }

            s.applyOnFollowers(func(n node.Node) {
                var nodeState nodeState.VolatileNodeState
                var AppendEntry rpcs.Rpc
                var rawMex []byte

                numStableNodes++
                nodeState,err = n.GetNodeState()
                if err!=nil {
                    log.Panicln(err)
                }

                if nodeState.GetMatchIndex() >= int(leaderCommitEntry) || !nodeState.Updated(){
                    return 
                }

                AppendEntry = AppendEntryRpc.NewAppendEntryRPC(
                    s._state,leaderCommitEntry-1,prevEntryTerm,[]*p.LogEntry{entryToCommit},leaderCommit)
                    rawMex,err = genericmessage.Encode(&AppendEntry)
                    if err != nil {
                        log.Panicln("error encoding AppendEntry: ",AppendEntry.ToString())
                    }

                    n.Send(rawMex)
                })
        }
    }
}

func (s *server) startNewElection(){
    log.Println("started new election");
    var entryTerm uint64 = s._state.GetTerm()

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
            var candidateLastLogIndex int
            var nodeState nodeState.VolatileNodeState

            nodeState,err = n.GetNodeState()
            if err != nil {
                log.Panicln(err)
            }

            candidateLastLogIndex = nodeState.GetMatchIndex()
            RequestVoteRPC.NewRequestVoteRPC(
                s._state.GetTerm(),
                s._state.GetIdPrivate(),
                int64(candidateLastLogIndex),
                entryTerm)

                raw_mex,err = genericmessage.Encode(&voteRequest)
                if err != nil {
                    log.Panicln("error in Encoding this voteRequest: ",(voteRequest).ToString())
                }
                log.Printf("sending to %v with key %v\n", n.GetIp(),n.GetIp())
                n.Send(raw_mex)
            })
}

func (s *server) leaderHearthBit(){
    //log.Println("start sending hearthbit")
    for s._state.Leader(){
        select{
        case <- s._state.HeartbeatTimeout().C:
            var hearthBit rpcs.Rpc

            hearthBit = AppendEntryRpc.GenerateHearthbeat(s._state)  
            log.Println("sending hearthbit")
            s.sendAll(&hearthBit)
            s._state.StartHearthbeatTimeout()
        }
    }
    s.setVolState()
    log.Println("no longer LEADER, stop sending hearthbit")
}

func (s *server) setVolState() {
    s.unstableNodes.Range(func(key, value any) bool {
            var nNode node.Node
            var nodeState nodeState.VolatileNodeState
            var found bool
            var err error

            nNode,found = value.(node.Node)
            if !found{
                panic("error type is not a node.Node")
            }

            nodeState,err = nNode.GetNodeState()
            if err != nil {
                log.Panicln(err)
            }

            nodeState.InitVolatileState(s._state.LastLogIndex())
        return true;
    })
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
        fn(nNode)
    }
}
