package server

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	messages "raft/internal/messages"
	append "raft/internal/messages/AppendEntryRPC"
	res_ap "raft/internal/messages/AppendEntryResponse"
	cop_st "raft/internal/messages/CopyStateRPC"
	req_vo "raft/internal/messages/RequestVoteRPC"
	res_vo "raft/internal/messages/RequestVoteResponse"
	"reflect"

	//"raft/internal/node/address"
	"raft/internal/node"
	state "raft/internal/raftstate"
	p "raft/pkg/protobuf"
	"sync"
)

type Server struct {
	_state         state.State
	messageChannel chan messages.Rpc
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


func NewServer(term uint64,ip_addr string, port string, role state.Role, serversIp []string) *Server {
	listener, err := net.Listen("tcp", ip_addr + ":" + port)

	if err != nil {
		log.Fatalf("Failed to listen on port %s: %s", port, err)
	}

	var server = &Server{
		_state:         state.NewState(term, port, role),
		otherNodes:     &sync.Map{},
		messageChannel: make(chan messages.Rpc),
		listener:       listener,
	}

    for i :=0; i< len(serversIp) -2; i++{
        var new_node,_= node.NewNode(serversIp[i],port)
        server.otherNodes.Store(generateID(serversIp[i]),new_node)
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
    s.otherNodes.Range(func (key any, value interface{}) bool{
        var nodeEle,errEl = value.(node.Node)
        if errEl == false{
            log.Println("invalid object in otherNodes map: ", reflect.TypeOf(nodeEle))
            return false
        }
        var ipAddr string = nodeEle.GetIp()
        var port string = nodeEle.GetPort()
		var conn, err = net.Dial("tcp",ipAddr + ":" + port)
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
        var value , found  = s.otherNodes.Load(id_node)

        if found {
            var connectedNode node.Node = value.(node.Node)
            connectedNode.AddConnIn(&conn)
        }else {
            var new_node,_= node.NewNode(newConncetionIp, newConncetionPort)        
            new_node.AddConnIn(&conn)
            s.otherNodes.Store(id_node,new_node)
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
            var message *messages.Rpc
            var errMes error
			message, errMes = conn.(node.Node).ReadRpc()
            if errMes != nil {
                fmt.Printf("error in reading from node %v with error %v", 
                    conn.(node.Node).GetIp(), errMes)
            }
			s.messageChannel <- *message
			return true
		})
	}
}

/*
 *  Send to every connected server a bunch of bytes
 */
func (s *Server) sendAll(mex messages.Rpc) {
	s.otherNodes.Range(func(k, conn interface{}) bool {
		log.Println(mex.ToString())
		conn.(node.Node).SendRpc(mex)
		return true
	})
}

/*
 *  Send AppendEntry message
 */
func (s *Server) sendAppendEntryRPC() {

	prevLogIndex := len(s._state.GetEntries()) - 2
	prevLogTerm := s._state.GetEntries()[prevLogIndex].GetTerm()
	appendEntry := append.NewAppendEntryRPC(
		s._state.GetTerm(),
		s._state.GetID(),
		uint64(prevLogIndex),
		prevLogTerm,
		make([]*p.Entry, 0),
		s._state.GetCommitIndex())

	s.sendAll(appendEntry)
}

/*
 *  Send RequestVote messages
 */
func (s *Server) sendRequestVoteRPC() {
	lastLogIndex := len(s._state.GetEntries()) - 1
	lastLogTerm := s._state.GetEntries()[lastLogIndex].GetTerm()
	requestVote := req_vo.NewRequestVoteRPC(
		s._state.GetTerm(),
		s._state.GetID(),
		uint64(lastLogIndex),
		lastLogTerm)

	s.sendAll(requestVote)
}

func (s *Server) startElection() {
	log.Println("/////////////////////////////////////Starting new election...", s._state.GetID())
	s._state.IncrementTerm()
	s._state.VoteFor(s._state.GetID())
	//s._state.ElectionTimeout()
	//s.sendRequestVoteRPC()
}

/*
 * Leader behaviour
 */
func (s *Server) leader() {
	select {
	case mess := <-s.messageChannel:
		switch mess.(type) {
		case *append.AppendEntryRPC:
			mess.ToString()
		//TODO Handle append entry

		case *res_ap.AppendEntryResponse:
			mess.ToString()
			// TODO Handle Response

		default:
			// Do nothing
		}
	case <-s._state.HeartbeatTimeout().C:
		hearthbit := append.NewAppendEntryRPC(
			s._state.GetTerm(),
			s._state.GetID(),
			0,
			0,
			make([]*p.Entry, 0),
			0)

		log.Println("Send heartbeat...", s._state.GetID())
		s._state.StartHearthbeatTimeout()
		s.sendAll(hearthbit)
	}
}

/*
 * Followers behaviour
 */
func (s *Server) follower() {
	select {
	case mess := <-s.messageChannel:
		switch mess.(type) {
		case *append.AppendEntryRPC:
			mess.ToString()
			log.Println("********** Message received ***********", mess.ToString())
			s._state.StartElectionTimeout()
			// TODO: Handle append entry

		case *req_vo.RequestVoteRPC:
			mess.ToString()
			// TODO: Handle request vote

		case *cop_st.CopyStateRPC:
			mess.ToString()
		// TODO: Handle copy state message

		default:
			// default behaviour: do nothing
		}
	case <-s._state.ElectionTimeout().C:
		if s._state.GetRole() != state.LEADER {
			s._state.SetRole(state.CANDIDATE)
			s.startElection()
		}
	}
}

/*
 * Handle candidates
 */
func (s *Server) candidate() {
	select {
	case mess := <-s.messageChannel:
		switch mess.(type) {
		case *append.AppendEntryRPC:
			mess.ToString()
			// TODO Handle response

		case *req_vo.RequestVoteRPC:
			mess.ToString()
			// TODO Handle request vote

		case *res_vo.RequestVoteResponse:
			mess.ToString()
			// TODO Handle Response vote

		case *cop_st.CopyStateRPC:
			mess.ToString()
			// TODO Handle copy state

		default:
			// default behaviour: do nothing
		}
	case <-s._state.ElectionTimeout().C:
		// TODO Handle election timeout during an election phase
	}
}

func (s *Server) run() {
	defer s.wg.Done()
	for {
		switch s._state.GetRole() {
		case state.LEADER:
			s.leader()
		case state.FOLLOWER:
			s.follower()
		case state.CANDIDATE:
			s.candidate()
		}
	}
}
