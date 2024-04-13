package server

import (
	"crypto/sha256"
	"encoding/hex"
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
	"raft/internal/rpcs/RequestVoteRPC"
	p "raft/pkg/rpcEncoding/out/protobuf"
	"reflect"
	"strings"
	"sync"
)

type pairMex struct{
    payload *rpcs.Rpc
    sender string
}

type Server struct {
	_state         state.State
	messageChannel chan pairMex
	otherNodes     *sync.Map
	listener       net.Listener
	wg             sync.WaitGroup
}

func generateID(input string) string {
	// Create a new SHA256 hasher
	hasher := sha256.New()

	// Write the input string to the hasher
	hasher.Write([]byte(input))

	// Get the hashed bytes
	hashedBytes := hasher.Sum(nil)

	// Convert the hashed bytes to a hexadecimal string
	id := hex.EncodeToString(hashedBytes)

	return id
}

func NewServer(term uint64, ipAddPrivate string, ipAddrPublic string, port string, serversIp []string) *Server {
	listener, err := net.Listen("tcp",":"+port)

	if err != nil {
		log.Fatalf("Failed to listen on port %s: %s", port, err)
	}

    log.Printf("my ip are: %v, %v\n",  ipAddPrivate, ipAddrPublic)

	var server = &Server{
		_state:         state.NewState(term, ipAddPrivate, ipAddrPublic, state.FOLLOWER),
		otherNodes:     &sync.Map{},
		messageChannel: make(chan pairMex),
		listener:       listener,
	}

    log.Println("number of others ip: ", len(serversIp))
	for i := 0; i < len(serversIp)-1; i++ {
		var new_node node.Node
        var nodeConn net.Conn
        var erroConn error
        var nodeId string

        nodeConn,erroConn = net.Dial("tcp",serversIp[i]+":"+port)
        if erroConn != nil {
            log.Println("Failed to connect to node: ", serversIp[i])
            continue
        }
        new_node = node.NewNode(serversIp[i], port, nodeConn)
        nodeId = generateID(serversIp[i])
        server.otherNodes.Store(nodeId, new_node)
        server._state.IncreaseNodeInCluster()
        go server.handleResponseSingleNode(nodeId, &new_node)

	}
	return server
}

func (s *Server) Start() {
    s.wg.Add(2)
    log.Println("Start accepting connections")
    go s.acceptIncomingConn()
    s._state.StartElectionTimeout()
    go s.run()

    s.wg.Wait()
}

func (s *Server) acceptIncomingConn() {
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
        var value any
        var connNode node.Node
		value, found = s.otherNodes.Load(id_node)
        log.Printf("new connec: %v\n", newConncetionIp)

        if strings.Contains(newConncetionIp, "10.0.0") {
            if found {
                connNode = value.(node.Node)
            } else {
                log.Printf("node with ip %v not found", newConncetionIp)
                var new_node node.Node = node.NewNode(newConncetionIp, newConncetionPort,conn)
                s.otherNodes.Store(id_node, new_node)
                s._state.IncreaseNodeInCluster()
                connNode = new_node
            }
            go s.handleResponseSingleNode(id_node,&connNode)
        }else{
            log.Println("new client request to cluster")
            if s._state.Leader(){
                conn.Write([]byte("ok\n"))
            }else{
                conn.Write([]byte(s._state.GetLeaderIpPublic()+"\n"))
            }
        }

	}
}

func (s *Server) handleResponseSingleNode(id_node string, workingNode *node.Node) {
    s.wg.Add(1)
    defer s.wg.Done()

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

func (s *Server) sendAll(rpc *rpcs.Rpc){
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

func (s *Server) run() {
	defer s.wg.Done()
	for {
		var mess pairMex
    /* To keep LastApplied and Leader's commitIndex always up to dated  */
    s._state.UpdateLastApplied()
    if s._state.Leader() {
        s._state.CheckCommitIndex(s.getMatchIndexes())
    }

        select {
        case mess = <-s.messageChannel:
            //log.Println("processing message: ", (*mess.payload).ToString())
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

func (s *Server) startNewElection(){
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

func (s *Server) leaderHearthBit(){
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

func (s *Server) getMatchIndexes() []int {
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

func (s *Server) setVolState() {
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
