package server


import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"sync"
	genericmessage "raft/internal/genericMessage"
	"raft/internal/node"
	"raft/internal/node/nodeState"
	"raft/internal/raftstate"
	state "raft/internal/raftstate"
	"raft/internal/rpcs"
	"raft/internal/rpcs/AppendEntryRpc"
	"raft/internal/rpcs/ClientReq"
	"raft/internal/rpcs/RequestVoteRPC"
	p "raft/pkg/raft-rpcProtobuf-messages/rpcEncoding/out/protobuf"
	"reflect"
	"strings"

	"google.golang.org/protobuf/proto"
)

type pairMex struct{
    payload *rpcs.Rpc
    sender string
}

type server struct {
	_state         state.State
	messageChannel chan pairMex
	otherNodes     *sync.Map
    clientNodes    *sync.Map
	listener       net.Listener
	wg             sync.WaitGroup
}

func (s *server) Start() {
    s.wg.Add(2)
    log.Println("Start accepting connections")
    go s.acceptIncomingConn()
    s._state.StartElectionTimeout()
    go s.run()

    s.wg.Wait()
}

//utility
func generateID(input string) string {
	hasher := sha256.New()
	hasher.Write([]byte(input))
	hashedBytes := hasher.Sum(nil)
	id := hex.EncodeToString(hashedBytes)
	return id
}

func (s *server) acceptIncomingConn() {
	defer s.wg.Done()
	for {
  //      log.Println("waiting new connection")
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
        var id_node string = generateID(newConncetionIp)
        var found bool
		_, found = s.otherNodes.Load(id_node)
        log.Printf("new connec: %v\n", newConncetionIp)

        if found {
            log.Println("adding a new node who is already in the cluster, probably a bug:",newConncetionIp)
            continue
        }
        log.Printf("node with ip %v not found", newConncetionIp)
        
        var newNode node.Node
        var errNode error
        newNode, errNode = node.NewNode(newConncetionIp, newConncetionPort,conn)
        if errNode != nil {
            log.Println("Error: ", errNode)
        }
        go s.handleConnection(id_node,&newNode)
	}
}

func (s* server) handleConnection(idNode string, workingNode *node.Node){
    s.wg.Add(1)
    defer s.wg.Done()

    var nodeIp string = (*workingNode).GetIp()

    if strings.Contains(nodeIp, "10.0.0") {
        s.otherNodes.Store(idNode, *workingNode)
        s._state.IncreaseNodeInCluster()
        s.handleResponseSingleNode(idNode,workingNode)
    }else{
        s.handleNewClientConnection(workingNode)
    }
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
            fmt.Printf("error in reading from node %v with error %v",(*client).GetIp(), err)
            (*client).CloseConnection()
            return
        }

        err = clientReq.Decode(mex)
        if err != nil {
            fmt.Printf("error in decoding client request from node %v with error %v",(*client).GetIp(), err)
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
            s.otherNodes.Delete(id_node);
            s._state.DecreaseNodeInCluster()
            break
        }
        if message != nil {
            s.messageChannel <- 
            pairMex{genericmessage.Decode(message),(*workingNode).GetIp()}
        }
    }

}

func (s *server) sendAll(rpc *rpcs.Rpc){
   log.Println("start broadcast")
    s.otherNodes.Range(func(key, value any) bool {
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
	defer s.wg.Done()
	for {
		var mess pairMex
    var errUpdt error 
    /* To keep LastApplied and Leader's commitIndex always up to dated  */
    errUpdt = s._state.UpdateLastApplied()
    if errUpdt != nil {
      log.Println("Error: ", errUpdt)
    } 
    if s._state.Leader() {
        s._state.CheckCommitIndex(s.getMatchIndexes())
        var lastSent int = s._state.CheckLastSent()
        if lastSent != -1 {
            s.sendAppendEntry()
        }
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

            f, ok = s.otherNodes.Load(generateID(sender))
            if !ok {
                log.Printf("Node %s not found", sender)
                continue
            }

            senderNode = f.(node.Node)
            oldRole = s._state.GetRole()
            rpcCall = mess.payload
            senderState = senderNode.GetNodeState()
            resp = (*rpcCall).Execute(&s._state, senderState)

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
    s.wg.Add(1)
    defer s.wg.Done()
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

func (s *server) sendAppendEntry() {
  var term uint64 = s._state.GetTerm()
  var privateId string = s._state.GetIdPrivate()
  var publicId string = s._state.GetIdPublic()
  var commitIndex = s._state.GetCommitIndex()
  var entries = s._state.GetEntries()

  s._state.UpdateLastSent()
  s.otherNodes.Range(func (key, value any) bool {
    var nNode node.Node
    var err bool
    var byEnc []byte
    var errEn error

    nNode,err = value.(node.Node)
    if !err {
       panic("error type is not a node.Node")
    }

    var appendEntry = AppendEntryRpc.NewAppendEntryRPC(term, privateId, publicId, (*nNode.GetNodeState()).GetNextIndex(), entries, commitIndex) 
    byEnc, errEn = genericmessage.Encode(&appendEntry)
    if errEn != nil{
       log.Panicln("error encoding this rpc: ", appendEntry.ToString())
    }
    nNode.Send(byEnc)
    return true
  })
}

func (s *server) getMatchIndexes() []int {
  var idxList = make([]int, 0) 
  s.otherNodes.Range(func(key, value any) bool {
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
    s.otherNodes.Range(func(key, value any) bool {
            var nNode node.Node
            var err bool

            nNode,err = value.(node.Node)
            if !err {
                panic("error type is not a node.Node")
            }
            nNode.ResetState(s._state.GetLastLogIndex())
        return true;
    })
}
