package server

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	messages "raft/internal/messages"
	"raft/internal/messages/AppendEntryRPC"
	"raft/internal/messages/RequestVoteRPC"
	"raft/internal/node"
	cutom_mex "raft/internal/node/message"
	"raft/internal/raftstate"
	state "raft/internal/raftstate"
	p "raft/pkg/protobuf"
	"reflect"
	"sync"
)

type pairMex struct{
    payload messages.Rpc
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
	listener, err := net.Listen("tcp", ip_addr+":"+port)

	if err != nil {
		log.Fatalf("Failed to listen on port %s: %s", port, err)
	}

	var server = &Server{
		_state:         state.NewState(term, ip_addr, state.FOLLOWER),
		otherNodes:     &sync.Map{},
		messageChannel: make(chan pairMex),
		listener:       listener,
	}

	for i := 0; i < len(serversIp)-2; i++ {
		var new_node node.Node
		var err error
		new_node, err = node.NewNode(serversIp[i], port)
		if err != nil {
			println("error in creating new node: %v", err)
			continue
		}
		server.otherNodes.Store(generateID(serversIp[i]), new_node)
	}
	return server
}

func (s *Server) Start() {
	s.wg.Add(3)

	go s.acceptIncomingConn()

	s.connectToServers()

	s._state.StartElectionTimeout()

	if s._state.Leader() {
		s._state.StartHearthbeatTimeout()
	}

	go s.run()

	go s.handleResponse()

	s.wg.Wait()
}

/*
 * Create a connection between this server to all the others and populate the map containing these connections
 */
func (s *Server) connectToServers() {
	s.otherNodes.Range(func(key any, value interface{}) bool {
		var nodeEle node.Node
		var errEl bool
		nodeEle, errEl = value.(node.Node)
		if !errEl {
			log.Println("invalid object in otherNodes map: ", reflect.TypeOf(nodeEle))
			return false
		}
		var ipAddr string = nodeEle.GetIp()
		var port string = nodeEle.GetPort()
		var conn, err = net.Dial("tcp", ipAddr+":"+port)
		for err != nil {
			log.Println("Dial error: ", err)
			return false
		}
		if conn != nil {
			nodeEle.AddConnOut(&conn)
		}
		return true
	})
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
		var value, found = s.otherNodes.Load(id_node)

		if found {
			var connectedNode node.Node = value.(node.Node)
			connectedNode.AddConnIn(&conn)
		} else {
			var new_node, _ = node.NewNode(newConncetionIp, newConncetionPort)
			new_node.AddConnIn(&conn)
			s.otherNodes.Store(id_node, new_node)
		}
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
            var node node.Node = conn.(node.Node)
			var message string
			var errMes error
			message, errMes = node.Recv()
			if errMes != nil {
				fmt.Printf("error in reading from node %v with error %v",
					node.GetIp(), errMes)
				return false
			}
			s.messageChannel <- 
                pairMex{*cutom_mex.NewMessage([]byte(message)).ToRpc(),node.GetIp()}
			return true
		})
	}
}

func (s *Server) sendAll(rpc *messages.Rpc){
    s.otherNodes.Range(func(key, value any) bool {
        var node node.Node = value.(node.Node)
        var mex cutom_mex.Message
        var raw_mex []byte

        mex = cutom_mex.FromRpc(*rpc)
        raw_mex = mex.ToByte()
        node.Send(raw_mex)
        return true
    })
}

func (s *Server) run() {
	defer s.wg.Done()
	for {
		var mess pairMex

        select {
        case mess = <-s.messageChannel:
            var rpcCall messages.Rpc
            var sender string
            var oldRole raftstate.Role
            var newRole raftstate.Role
            var resp *messages.Rpc
            var mex cutom_mex.Message

            oldRole = s._state.GetRole()
            rpcCall = mess.payload
            sender = mess.sender
            resp = rpcCall.Execute(&s._state)
            newRole = s._state.GetRole()

            if resp != nil {
                mex = cutom_mex.FromRpc(*resp)
                var f any
                var ok bool
                f, ok = s.otherNodes.Load(generateID(sender))

                if !ok {
                    log.Printf("Node %s not found", sender)
                    continue
                }

                f.(node.Node).Send(mex.ToByte())
            }

            if newRole == state.LEADER && oldRole != newRole {
                go s.leaderHearthBit()
            }

        case <-s._state.ElectionTimeout().C:
            s.startNewElection()
        }
	}
}

func (s *Server) startNewElection(){
    var entries []p.Entry
    var len_ent int
    var voteRequest messages.Rpc

    s._state.IncrementTerm()
    len_ent = len(entries) -1
    log.Println("len entries:", len_ent)
    log.Println("entries: ", entries)
    entries = s._state.GetEntries()

    voteRequest = RequestVoteRPC.NewRequestVoteRPC(
        s._state.GetTerm(),
        s._state.GetId(),
        0,
        0)
        // uint64(len_ent),
        // entries[len_ent].GetTerm())

    s.sendAll(&voteRequest)
}

//TODO: finish creation of hearthBit Rpc, the current field are wrong and only for testing purpose
func (s *Server) leaderHearthBit(){
    for s._state.GetRole() == state.LEADER{
        select{
        case <- s._state.HeartbeatTimeout().C:
            var hearthBit messages.Rpc = AppendEntryRPC.NewAppendEntryRPC(
                s._state.GetTerm(),
                s._state.GetId(),
                0,
                0,
                nil,
                0,
            )
            s.sendAll(&hearthBit)
        }
    }
}
