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
	"raft/internal/rpcs/NewConfiguration"
	"raft/internal/rpcs/RequestVoteRPC"
	"raft/internal/rpcs/UpdateNode"
	"raft/pkg/raft-rpcProtobuf-messages/rpcEncoding/out/protobuf"
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
	_state         state.State
	messageChannel chan pairMex
	unstableNodes   *sync.Map
	stableNodes     *sync.Map
    clientNodes    *sync.Map
	listener       net.Listener
	wg             sync.WaitGroup
}

func (s *server) Start() {
    log.Println("Start accepting connections")
    s._state.StartElectionTimeout()

    s.wg.Add(1)
    go s.acceptIncomingConn()
    s.wg.Add(1)
    go s.run()

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
        new_node = node.NewNode(serversIp[i], port, nodeConn)
        (*this).unstableNodes.Store(new_node.GetIp(), new_node)
        go (*this).handleResponseSingleNode(new_node.GetIp(), &new_node)

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
        var new_node node.Node = node.NewNode(newConncetionIp, newConncetionPort,conn)
        go func ()  {
            s.wg.Add(1)
            defer s.wg.Done()
            s.handleConnection(newConncetionIp,&new_node)
        }()
	}
}

func (s* server) handleConnection(idNode string, workingNode *node.Node){
    var nodeIp string = (*workingNode).GetIp()

    if strings.Contains(nodeIp, "10.0.0") {
        s.handleResponseSingleNode(idNode,workingNode)
        return
    }
    s.handleNewClientConnection(workingNode)
}

func (s* server) handleNewClientConnection(client *node.Node){
    log.Println("new client request to cluster")
    if s._state.Leader(){
        var clientReq rpcs.Rpc = &ClientReq.ClientReq{}
        var ok = "ok"
        var leaderIp p.PublicIp = p.PublicIp{IP: ok,}
        var mex,err = proto.Marshal(&leaderIp)

        if err != nil {
            log.Panicln("error encoding confirmation leader public ip for client:",err)
        }

        (*client).Send(mex)

        mex,err = (*client).Recv()
        if err != nil {
            fmt.Printf("error in reading from node %v with error %v\n",(*client).GetIp(), err)
            (*client).CloseConnection()
            return
        }

        err = clientReq.Decode(mex)
        if err != nil {
            fmt.Printf("error in decoding client request from node %v with error %v\n",(*client).GetIp(), err)
            (*client).CloseConnection()
            return
        }
        log.Println("managing client Request: ", clientReq.ToString())
        clientReq.Execute(&s._state,(*client).GetNodeState())
    }else{
        var leaderIp p.PublicIp = p.PublicIp{IP: s._state.GetLeaderIpPublic(),}
        var mex,err = proto.Marshal(&leaderIp)
        log.Printf("sending public ip of leader: %v\n", s._state.GetLeaderIpPublic())

        if err != nil {
            log.Panicln("error encoding leader public ip for client:",err)
        }

        (*client).Send(mex)
    }
    (*client).CloseConnection()
}

func (s *server) handleResponseSingleNode(id_node string, workingNode *node.Node) {
    if s._state.Leader() {
        go func (){
            s.wg.Add(1)
            defer s.wg.Done()
            s.joinConf(id_node,workingNode)
        }()
    }
    for{
        var message []byte
        var errMes error
        message, errMes = (*workingNode).Recv()
        if errMes != nil {
            fmt.Printf("error in reading from node %v with error %v",(*workingNode).GetIp(), errMes)
            if !s._state.Leader() {
                s._state.StartElectionTimeout()
            }
            (*workingNode).CloseConnection()
            s.stableNodes.Delete(id_node);
            s.unstableNodes.Delete(id_node);
            s._state.DecreaseNodeInCluster()
            break
        }
        if message != nil {
            s.messageChannel <- 
            pairMex{genericmessage.Decode(message),(*workingNode).GetIp()}
        }
    }

}

func (s *server) joinConf(id_node string, workingNode *node.Node){
    var newConfRequest rpcs.Rpc = NewConfiguration.NewNewConfigurationRPC(append(s._state.GetConfig(),id_node))
    var newConfEntry p.LogEntry = p.LogEntry{
        OpType: p.Operation_JOIN_CONF,
        Term: s._state.GetTerm(),
        Payload: nil,
        Description: "added new node " + id_node + " to configuration: ",
    }

    log.Printf("adding node %v to the stable queue\n", (*workingNode).GetIp())
    s.stableNodes.Store((*workingNode).GetIp(), *workingNode)
    s._state.AppendEntries([]*p.LogEntry{&newConfEntry},(*s)._state.LastLogIndex()+1)
    if s._state.Leader() {
        s._state.UpdateConfiguration([]string{(*workingNode).GetIp()})
        s._state.IncreaseNodeInCluster()
        s.sendAll(&newConfRequest)
        s.updateNewNode(id_node,workingNode)              
    }
}

func (s *server) updateNewNode(id_node string,workingNode *node.Node){
    var volatileState *nodeState.VolatileNodeState = (*workingNode).GetNodeState()
    var index = 0
    for  i,e := range s._state.GetEntries() {
        generateUpdateRequest(workingNode,false,e)
        for  (*volatileState).GetMatchIndex() < i {
            //WARN: WAIT
        }
        index = i
    }
    (*s)._state.IncreaseNodeInCluster()
    s.stableNodes.Store(id_node,*workingNode)
    generateUpdateRequest(workingNode,true,nil)
    for  (*volatileState).GetMatchIndex() < index+1 {
        //WARN: WAIT
    }
    s._state.CommitConfig()
}

func generateUpdateRequest(workingNode *node.Node, voting bool, entry *protobuf.LogEntry){
    var updateReq rpcs.Rpc 
    var mex []byte
    var err error

    updateReq = UpdateNode.NewUpdateNodeRPC(voting, entry)
    mex,err = genericmessage.Encode(&updateReq)
    if err != nil {
        log.Panic("error encoding UpdateNode rpc")
    }
    (*workingNode).Send(mex)
}

func (s *server) sendAll(rpc *rpcs.Rpc){
   log.Println("start broadcast")
    s.stableNodes.Range(func(key, value any) bool {
        var nNode node.Node 
        var found bool 
        nNode, found = value.(node.Node)
        if !found {
            var s = reflect.TypeOf(value)
            log.Panicln("failed conversion type node, type is: ", s)
        }
        var raw_mex []byte
        var err error

        raw_mex,err = genericmessage.Encode(rpc)
        if err != nil {
            log.Panicln("error in Encoding this rpc: ",(*rpc).ToString())
        }
    //    log.Printf("sending: %v to %v", (*rpc).ToString(), (nNode).GetIp() )
        log.Printf("sending to %v\n", nNode.GetIp())
        nNode.Send(raw_mex)
        return true
    })
   log.Println("end broadcast")
}

func (s *server) run() {
    for {
        var mess pairMex
        /* To keep LastApplied and Leader's commitIndex always up to dated  */
        s._state.UpdateLastApplied()
        if s._state.Leader() {
            s._state.CheckCommitIndex(s.getMatchIndexes())
        }

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
            var senderState *nodeState.VolatileNodeState
            var senderNode node.Node
            var newConf []string
            var failedConn []string

            f, ok = s.unstableNodes.Load(sender)
            if !ok {
                log.Printf("Node %s not found", sender)
                continue
            }

            senderNode = f.(node.Node)
            oldRole = s._state.GetRole()
            rpcCall = mess.payload
            senderState = senderNode.GetNodeState()
            resp = (*rpcCall).Execute(&s._state, senderState)

            if !s._state.ConfStatus() {
                newConf = s._state.GetConfig()
                for _, v := range newConf {
                    var _,found = s.stableNodes.Load(v)
                    if !found {
                        var e,found = s.unstableNodes.Load(v)
                        var newNode node.Node = e.(node.Node)
                        if !found {
                            failedConn,errEn = s.connectToNodes([]string{v},"8080") //WARN: hard encoding port
                            if errEn != nil {
                                for _, v := range failedConn {
                                    log.Println("failed connecting to this node: " + v)
                                }
                                log.Panic("devel debug") //WARN: to remove in the future
                            }
                        }
                        s.stableNodes.Store(v,newNode)
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
                go func (){
                    s.wg.Add(1)
                    defer s.wg.Done()
                    s.leaderHearthBit()
                }()
            }
            //         log.Println("rpc processed")
        case <-s._state.ElectionTimeout().C:
            if !s._state.Leader() {
                s.startNewElection()
            }
        }
    }
}

func (s *server) startNewElection(){
    log.Println("started new election");
    var entries []*p.LogEntry
    var len_ent int
    var voteRequest rpcs.Rpc
    var entryTerm uint64 = 0

    s._state.IncrementTerm()

    entries = s._state.GetEntries()
    len_ent = len(entries)
    if len_ent-1 > 0{
        entryTerm = entries[len_ent-1].GetTerm()
    }

    voteRequest = RequestVoteRPC.NewRequestVoteRPC(
        s._state.GetTerm(),
        s._state.GetIdPrivate(),
        int64(len_ent),
        entryTerm)

    s._state.IncreaseSupporters()
    //log.Println("node in cluster: ",s._state.GetNumNodeInCluster())
    if s._state.GetNumNodeInCluster() == 1 {
        log.Println("became leader: ",s._state.GetRole())
        s._state.SetRole(raftstate.LEADER)
        s._state.SetLeaderIpPrivate(s._state.GetIdPrivate())
        s._state.SetLeaderIpPublic(s._state.GetIdPublic())
        s._state.ResetElection()
        go s.leaderHearthBit()
    }else {
     //   log.Println("sending to everybody request vote :" + voteRequest.ToString())
        s.sendAll(&voteRequest)
    }
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

func (s *server) getMatchIndexes() []int {
  var idxList = make([]int, 0) 
  s.stableNodes.Range(func(key, value any) bool {
    var nNode node.Node
    var err bool

    nNode,err = value.(node.Node)
    if !err {
      panic("error type is not a node.Node")
    }

    idxList = append(idxList, (*nNode.GetNodeState()).GetMatchIndex())
    return true
  })
  return idxList
}

func (s *server) setVolState() {
    s.stableNodes.Range(func(key, value any) bool {
            var nNode node.Node
            var err bool

            nNode,err = value.(node.Node)
            if !err {
                panic("error type is not a node.Node")
            }
            nNode.ResetState(s._state.LastLogIndex())
        return true;
    })
}
