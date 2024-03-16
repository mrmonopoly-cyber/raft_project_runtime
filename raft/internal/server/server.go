package server

import (
	"log"
	"net"
	messages "raft/internal/messages"
	append "raft/internal/messages/AppendEntryRPC"
	res_ap "raft/internal/messages/AppendEntryResponse"
	cop_st "raft/internal/messages/CopyStateRPC"
	req_vo "raft/internal/messages/RequestVoteRPC"
	res_vo "raft/internal/messages/RequestVoteResponse"
	//"raft/internal/node/address"
	p "raft/pkg/protobuf"
	"sync"
	"raft/internal/node"
	state "raft/internal/raftstate"
)

type Server struct {
	_state         state.State
	messageChannel chan messages.Rpc
	otherNodes     *sync.Map
	listener       net.Listener
	wg             sync.WaitGroup
}

func NewServer(term uint64, id string, role state.Role, serversId []string) *Server {
	listener, err := net.Listen("tcp", "localhost:"+id)

	if err != nil {
		log.Fatalf("Failed to listen on port %s: %s", id, err)
	}

	var server = &Server{
		_state:         state.NewState(term, id, role, serversId),
		otherNodes:     &sync.Map{},
		messageChannel: make(chan messages.Rpc),
		listener:       listener,
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
	for i := 0; i < len(s._state.GetServersID()); i++ {
		var serverId = s._state.GetServersID()[i]
		if serverId != s._state.GetID() {
			var conn, err = net.Dial("tcp", "localhost:"+serverId)

			for err != nil {
				log.Println("Dial error: ", err)
				log.Println("retrying...")
				conn, err = net.Dial("tcp", "localhost:"+serverId)
			}

			if conn != nil {
				newNode, err := node.NewNode(conn.RemoteAddr().String())
				if err != nil {
					log.Println("Error in conversion: ", err)
					return
				}
				newNode.AddConnOut(&conn)

				value, ok := s.otherNodes.Load(newNode.GetIp())
				if ok {
					value.(node.Node).AddConnOut(&conn)
					s.otherNodes.Swap(newNode.GetIp(), value)
				} else {
					s.otherNodes.Store(newNode.GetIp(), newNode)
				}
			}
		}
	}
}

func (s *Server) acceptIncomingConn() {
	defer s.wg.Done()
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Println("Failed on accept: ", err)
			continue
		}

		newNode, err := node.NewNode(conn.RemoteAddr().String())
		if err != nil {
			log.Println("Error during conversion: ", err)
		}
		newNode.AddConnIn(&conn)

		value, ok := s.otherNodes.Load(newNode.GetIp())
		if ok {
			value.(node.Node).AddConnIn(&conn)
			s.otherNodes.Swap(newNode.GetIp(), value)
		} else {
			s.otherNodes.Store(newNode.GetIp(), newNode)
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
			message, _ := conn.(node.Node).ReadRpc()

			// for every []byte: decode it to RPC
			var mess *messages.Message = messages.NewMessage([]byte(message))
			var recRpc messages.Rpc
			switch mess.Mex_type {
			case messages.APPEND_ENTRY:
				var rpc append.AppendEntryRPC
				recRpc = &rpc
			case messages.APPEND_ENTRY_RESPONSE:
				var rpc res_ap.AppendEntryResponse
				recRpc = &rpc
			case messages.REQUEST_VOTE:
				var rpc req_vo.RequestVoteRPC
				recRpc = &rpc
			case messages.REQUEST_VOTE_RESPONSE:
				var rpc res_vo.RequestVoteResponse
				recRpc = &rpc
			case messages.COPY_STATE:
				var rpc cop_st.CopyStateRPC
				recRpc = &rpc
			default:
				panic("impossible case")
			}
			recRpc.Decode(mess.Payload)
			s.messageChannel <- recRpc

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
