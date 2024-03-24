package server

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net"
	"raft/internal/messages"
	"raft/internal/messages/AppendEntryRpc"
	"raft/internal/messages/RequestVoteRPC"
	"raft/internal/node"
	"raft/internal/raftstate"
	state "raft/internal/raftstate"
	p "raft/pkg/rpcEncoding/out/protobuf"
	"reflect"
	"sync"
)

type pairMex struct{
    payload *messages.Rpc
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

func NewServer(term uint64, ip_addr string, port string, serversIp []string) *Server {
	listener, err := net.Listen("tcp",":"+port)

	if err != nil {
		log.Fatalf("Failed to listen on port %s: %s", port, err)
	}

	var server = &Server{
		_state:         state.NewState(term, ip_addr, state.FOLLOWER),
		otherNodes:     &sync.Map{},
		messageChannel: make(chan pairMex),
		listener:       listener,
	}

    log.Println("number of others ip: ", len(serversIp))
	for i := 0; i < len(serversIp)-1; i++ {
		var new_node node.Node
		new_node = node.NewNode(serversIp[i], port)
        log.Printf("connecting to the server: %v\n", serversIp[i])
        var nodeConn net.Conn
        var erroConn error
        var nodeId string

        nodeConn,erroConn = net.Dial("tcp",serversIp[i]+":"+port)
        if erroConn != nil {
            log.Println("Failed to connect to node: ", serversIp[i])
            continue
        }
        new_node.AddConn(nodeConn)
        log.Println("storing new node with ip :", serversIp[i])
        nodeId = generateID(serversIp[i])
        server.otherNodes.Store(nodeId, new_node)
        server._state.IncreaseNodeInCluster()

	}
	return server
}

func (s *Server) Start() {
	s.wg.Add(3)


    log.Println("Start accepting connections")
	go s.acceptIncomingConn()


    log.Println("Start election Timeout")
	s._state.StartElectionTimeout()

    log.Println("start main run")
	go s.run()

    log.Println("start handle response")
	go s.handleResponse()

    log.Println("wait to finish")
	s.wg.Wait()
}

func (s *Server) acceptIncomingConn() {
	defer s.wg.Done()
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
		var id_node string = generateID(newConncetionIp)
        var found bool
		_, found = s.otherNodes.Load(id_node)
        
        log.Println("enstablish connection with node: ", newConncetionIp)

		if found {
            log.Printf("node with ip %v found", newConncetionIp)
            continue
		} else {
            log.Printf("node with ip %v not found", newConncetionIp)
			var new_node node.Node = node.NewNode(newConncetionIp, newConncetionPort)
			new_node.AddConn(conn)
			s.otherNodes.Store(id_node, new_node)
            s._state.IncreaseNodeInCluster()
		}
        log.Println("finish accepting new connections")
	}
}

/*
 * Handle incoming messages and send it to the corresponding channel
 */
func (s *Server) handleResponse() {
	defer s.wg.Done()
	// iterating over the connections map and receive byte message
	for {
		s.otherNodes.Range(func(k, conn interface{}) bool {
            // var node *node.Node = conn.(*node.Node)
            var nNode node.Node
            var err bool

            nNode,err = conn.(node.Node)

            if !err {
                panic("error type is not a node.Node")
            }

			var message string
			var errMes error
			message, errMes = (nNode).Recv()
            if errMes == errors.New("connection not instantiated"){
                return false
            }
			if errMes != nil {
				fmt.Printf("error in reading from node %v with error %v",
					(nNode).GetIp(), errMes)
				return false
			}
            if message != "" {
                log.Println("received message from: " + (nNode).GetIp())
                log.Println("data of message: " + message )
                s.messageChannel <- 
                // pairMex{custom_mex.NewMessage([]byte(message)).ToRpc(),(nNode).GetIp()}
                pairMex{messages.Decode(message),(nNode).GetIp()}
            }
			return true
		})
	}
}

func (s *Server) sendAll(rpc *messages.Rpc){
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

        raw_mex,err = (*rpc).Encode()
        if err != nil {
            log.Panicln("error in Encoding this rpc: ",(*rpc).ToString())
        }
        log.Printf("sending: %v to %v", (*rpc).ToString(), (nNode).GetIp() )
        nNode.Send(raw_mex)
        return true
    })
    log.Println("end broadcast")
}

func (s *Server) run() {
	defer s.wg.Done()
	for {
		var mess pairMex

        select {
        case mess = <-s.messageChannel:
            log.Println("processing message: ", (*mess.payload).ToString())
            var rpcCall *messages.Rpc
            var sender string
            var oldRole raftstate.Role
            var resp *messages.Rpc
            var byEnc []byte
            var errEn error

            oldRole = s._state.GetRole()
            rpcCall = mess.payload
            sender = mess.sender
            resp = (*rpcCall).Execute(&s._state)

            if resp != nil {
                log.Println("reponse to send to: ", sender)
                var f any
                var ok bool
                f, ok = s.otherNodes.Load(generateID(sender))

                if !ok {
                    log.Printf("Node %s not found", sender)
                    continue
                }

                log.Println("sending mex to: ",sender)
                // f.(node.Node).Send(mex.ToByte())
                byEnc, errEn = (*resp).Encode()
                if errEn != nil{
                    log.Panicln("error encoding this rpc: ", (*resp).ToString())
                }
                f.(node.Node).Send(byEnc)
            }

            if s._state.Leader() && oldRole != state.LEADER{
                s.wg.Add(1)
                go s.leaderHearthBit()
            }

            log.Println("rpc processed")
        case <-s._state.ElectionTimeout().C:
            s.startNewElection()
        }
	}
}

func (s *Server) startNewElection(){
    var entries []p.LogEntry
    var len_ent int
    var voteRequest messages.Rpc
    var entryTerm uint64 = 0

    s._state.IncrementTerm()

    entries = s._state.GetEntries()
    len_ent = len(entries)
    if len_ent > 0{
        entryTerm = entries[len_ent].GetTerm()
    }


    voteRequest = RequestVoteRPC.NewRequestVoteRPC(
        s._state.GetTerm(),
        s._state.GetId(),
        uint64(len_ent),
        entryTerm)

    s._state.IncreaseSupporters()
    log.Println("node in cluster: ",s._state.GetNumNodeInCluster())
    if s._state.GetNumNodeInCluster() == 1 {
        log.Println("became leader: ",s._state.GetRole())
        s._state.SetRole(raftstate.LEADER)
        s._state.ResetElection()
        go s.leaderHearthBit()
    }else {
        log.Println("sending to everybody request vote :" + voteRequest.ToString())
        s.sendAll(&voteRequest)
    }
}

func (s *Server) leaderHearthBit(){
    defer s.wg.Done()
    log.Println("start sending hearthbit")
    for s._state.Leader(){
        select{
        case <- s._state.HeartbeatTimeout().C:
            var hearthBit messages.Rpc

            // hearthBit = AppendEntryRPC.GenerateHearthbeat(s._state)
            hearthBit = AppendEntryRpc.GenerateHearthbeat(s._state)  
            log.Println("sending hearthbit")
            s.sendAll(&hearthBit)
            s._state.StartHearthbeatTimeout()
        }
    }
    log.Println("no longer LEADER, stop sending hearthbit")
}