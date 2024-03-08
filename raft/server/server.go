package server

import (
	"log"
	"net"
	"net/http"
	"sync"
  m "raft/messages"
)

type Server struct {
	_state raftStateImpl
	connections        *sync.Map
	appendEntryChannel chan []byte
  requestVoteChannel chan m.RequestVoteRPC
  appendEntryResponseChannel chan m.AppendEntryResponse
  requesVoteRespChannel chan m.RequestVoteResponse
  copyStateChannel chan m.CopyStateRPC
}


func NewServer(term int, id string, role Role, serversId []string) *Server {
	var server = &Server{
    *NewState(term, id, role, serversId),
		&sync.Map{},
		make(chan []byte),
    make(chan m.RequestVoteRPC),
    make(chan m.AppendEntryResponse),
    make(chan m.RequestVoteResponse),
    make(chan m.CopyStateRPC),
	}

  return server
}

func (s *Server) Start() {

  s.startListen()

	s.ConnectToServers()

  go s.handleResponse(&s.appendEntryChannel)

	s._state.StartElectionTimeout()

	if s._state.Leader() {
		s._state.StartHearthbeatTimeout()
	}

	go s.run()
}

/*
 * Create a connection between this server to all the others and populate the map containing these connections 
*/
func (s *Server) ConnectToServers() {
	//go acceptIncomingConn(list, &s.appendEntryChannel)

	for i := 0; i < len(s._state.GetServersID()); i++ {
		var serverId = s._state.GetServersID()[i]
		if serverId != s._state.id {
			conn, err := net.Dial("tcp", "localhost:"+serverId)

			if err != nil {
        log.Println("Dial error: ", err)
			}
			log.Println("Server ", s._state.id)

      if conn != nil {
			  s.connections.Store(serverId, conn)
      }
		}
	}
}

func (s *Server) startListen() {
  listener, err := net.Listen("tcp", "localhost:" + s._state.GetID())
  
  if err != nil {
    log.Println("Error during Listen: ", err)
  }

  go http.Serve(listener, nil)
  log.Println("Listener up on port: " + s._state.GetID())
}

/*
 * Handle incoming messages and send it to the corresponding channel
*/
func (s *Server) handleResponse(channel *chan []byte) {
	buffer := make([]byte, 2048)

	for {
    s.connections.Range(func(k, conn interface{}) bool {
		  _, err := conn.(net.Conn).Read(buffer)
		  if err != nil {
			  log.Println("Read error: ", err)
		  }
      *channel <- buffer
      
      return true
    })
	}

}

func (s *Server) sendAll() {

	data1 := []byte("m")
	s.connections.Range(func(k, conn interface{}) bool {
		conn.(net.Conn).Write(data1)
		return true
	})
}

func (s *Server) sendHeartbeat() {
	if s._state.Leader() {
		//var appendEntry = r.NewAppendEntry(s._state.GetTerm(), s._state.GetID(), 0, 0, make([]l.Entry, 0), 0)
		log.Println("Send heartbeat...", s._state.id)
		s._state.StartHearthbeatTimeout()
		s.sendAll()
	}
}

func (s *Server) startElection() {
	log.Println("/////////////////////////////////////Starting new election...", s._state.id)
}

/*
 * Handle leader
*/
func (s *Server) leader() {
  log.Println("enter in leader")
	select {
	case <-s._state.HeartbeatTimeout().C:
		s.sendHeartbeat()
	}
	
}

/*
 * Handle followers
*/
func (s *Server) follower() {
  log.Println("enter in follower")
	select {
		case n := <-s.appendEntryChannel:
			log.Println("********** Message received ***********", n[0])
			s._state.StartElectionTimeout()
		case <-s._state.ElectionTimeout().C:
			if s._state.GetRole() != LEADER {
				s.startElection()
			}
		}	
}

/*
 * Handle candidates
*/
func (s *Server) candidate() {
  // TODO: implement candidate behaviour
}

func (s *Server) run() {
	for {
    switch s._state.GetRole() {
	  case LEADER: 
		  s.leader()
	  case FOLLOWER:
		  s.follower()
	  case CANDIDATE:
		  s.candidate()
	  }
  }
}
