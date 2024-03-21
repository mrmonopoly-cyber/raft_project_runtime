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
	custom_mex "raft/internal/node/message"
	"raft/internal/raftstate"
	state "raft/internal/raftstate"
	p "raft/pkg/protobuf"
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
	listener, err := net.Listen("tcp", ":"+port)

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


    log.Println("Start accepting connections")
	// go s.acceptIncomingConn()

    log.Println("connect To other Servers")
	// s.connectToServers()

    log.Println("Start election Timeout")
	// s._state.StartElectionTimeout()

    // log.Println("if")
	// if s._state.Leader() {
	// 	s._state.StartHearthbeatTimeout()
	// }

    log.Println("start main run")
	// go s.run()

    log.Println("start handle response")
	// go s.handleResponse()

    log.Println("wait to finish")
	s.wg.Wait()
}

/*
 * Create a connection between this server to all the others and populate the map containing these connections
 */
func (s *Server) connectToServers() {
    log.Println("connecting to list of nodes")
    if s.otherNodes == nil {
        panic("Map of Node not allocated")
    }
	s.otherNodes.Range(func(key any, value interface{}) bool {
        log.Println("connecting to a node")
		var nodeEle node.Node
		var errEl bool
		nodeEle, errEl = value.(node.Node)
		if !errEl {
			log.Println("invalid object in otherNodes map: ", reflect.TypeOf(nodeEle))
			return false
		}
		var ipAddr string = nodeEle.GetIp()
		var port string = nodeEle.GetPort()
        log.Println("connecting to: " + ipAddr + ":" + "8080")
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
        
        log.Println("enstablish connection with node: ", id_node)

		if found {
			var connectedNode node.Node = value.(node.Node)
			connectedNode.AddConnIn(&conn)
		} else {
			var new_node, _ = node.NewNode(newConncetionIp, newConncetionPort)
			new_node.AddConnIn(&conn)
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
            var node node.Node = conn.(node.Node)
			var message string
			var errMes error
			message, errMes = node.Recv()
			if errMes != nil {
				fmt.Printf("error in reading from node %v with error %v",
					node.GetIp(), errMes)
				return false
			}
            if message != "" {
                s.messageChannel <- 
                pairMex{custom_mex.NewMessage([]byte(message)).ToRpc(),node.GetIp()}
            }
			return true
		})
	}
}

func (s *Server) sendAll(rpc *messages.Rpc){
    s.otherNodes.Range(func(key, value any) bool {
        var node node.Node = value.(node.Node)
        var mex custom_mex.Message
        var raw_mex []byte

        mex = custom_mex.FromRpc(*rpc)
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
            var rpcCall *messages.Rpc
            var sender string
            var oldRole raftstate.Role
            var newRole raftstate.Role
            var resp *messages.Rpc
            var mex custom_mex.Message

            oldRole = s._state.GetRole()
            rpcCall = mess.payload
            sender = mess.sender
            resp = (*rpcCall).Execute(&s._state)
            newRole = s._state.GetRole()

            if resp != nil {
                mex = custom_mex.FromRpc(*resp)
                var f any
                var ok bool
                f, ok = s.otherNodes.Load(generateID(sender))

                if !ok {
                    log.Printf("Node %s not found", sender)
                    continue
                }

                f.(node.Node).Send(mex.ToByte())
            }

            if newRole == state.LEADER && oldRole != state.LEADER{
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
        s.sendAll(&voteRequest)
    }
}

func (s *Server) leaderHearthBit(){
    log.Println("start sending hearthbit")
    for s._state.Leader(){
        select{
        case <- s._state.HeartbeatTimeout().C:
            var hearthBit messages.Rpc

            hearthBit = AppendEntryRPC.GenerateHearthbeat(s._state)
            log.Println("sending hearthbit")
            s.sendAll(&hearthBit)
            s._state.StartHearthbeatTimeout()
        }
    }
    log.Println("no longer LEADER, stop sending hearthbit")
}
