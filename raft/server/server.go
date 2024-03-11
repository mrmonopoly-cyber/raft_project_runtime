package server

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	m "raft/messages"
	p "raft/protobuf"
	"sync"
)

type EnumType int

const (
  APPEND_ENTRY EnumType = iota
  REQUEST_VOTE 
  APPEND_RESPONSE
  VOTE_RESPONSE
  COPY_STATE
)

type Message struct {
  Ty EnumType
  Payload []byte
}

func (m *Message) toByte() []byte {
  result, _ := json.Marshal(m)
  return result
}

func (m *Message) toMessage(b []byte) {
  json.Unmarshal(b, m)
}

type Server struct {
	_state raftStateImpl
	connections        *sync.Map
	appendEntryChannel chan m.AppendEntryRPC
  requestVoteChannel chan m.RequestVoteRPC
  appendEntryResponseChannel chan m.AppendEntryResponse
  requesVoteRespChannel chan m.RequestVoteResponse
  copyStateChannel chan m.CopyStateRPC
}


func NewServer(term uint64, id string, role Role, serversId []string) *Server {
	var server = &Server{
    *NewState(term, id, role, serversId),
		&sync.Map{},
		make(chan m.AppendEntryRPC),
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

  go s.handleResponse()

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
			var conn, err = net.Dial("tcp", serverId + "8080")

			for err == nil {
        log.Println("Dial error: ", err)
        log.Println("retrying...")
        conn, err = net.Dial("tcp", serverId + "8080")
			}
			log.Println("Server ", s._state.id)

      if conn != nil {
			  s.connections.Store(serverId, conn)
      }
		}
	}
}

func (s *Server) startListen() {
  listener, err := net.Listen("tcp", "localhost:8080")
  
  if err != nil {
    log.Panic(err) 
  }

  go http.Serve(listener, nil)
  log.Println("Listener up on port: " + listener.Addr().String())
}

/*
 * Handle incoming messages and send it to the corresponding channel
*/
func (s *Server) handleResponse() {
	buffer := make([]byte, 2048)

  // iterating over the connections map and receive byte messages 
	for {
    s.connections.Range(func(k, conn interface{}) bool {
		  _, err := conn.(net.Conn).Read(buffer)
		  if err != nil {
			  log.Println("Read error: ", err)
		  }
      
    // for every []byte: decode it to Message type
    var mess Message
    mess.toMessage(buffer)

    // filter by type of message and send it to the appropriate channel
    switch mess.Ty {
      case APPEND_ENTRY:
        var appendEntry = &m.AppendEntryRPC{}
        appendEntry.Decode(mess.Payload)
        s.appendEntryChannel <- *appendEntry

      case REQUEST_VOTE:
        var reqVote = &m.RequestVoteRPC{}
        reqVote.Decode(mess.Payload)
        s.requestVoteChannel <- *reqVote

      case APPEND_RESPONSE: 
        var appendResponse = &m.AppendEntryResponse{}
        appendResponse.Decode(mess.Payload)
        s.appendEntryResponseChannel <- *appendResponse

      case VOTE_RESPONSE:
        var voteResponse = &m.RequestVoteResponse{}
        voteResponse.Decode(mess.Payload)
        s.requesVoteRespChannel <- *voteResponse

      case COPY_STATE:
        var copyState = &m.CopyStateRPC{}
        copyState.Decode(mess.Payload)
        s.copyStateChannel <- *copyState
    }
      
    return true
    })
	}

}

/*
 *  Send to every connected server a bunch of bytes
*/
func (s *Server) sendAll(mess []byte) {

	s.connections.Range(func(k, conn interface{}) bool {
		conn.(net.Conn).Write(mess)
		return true
	})
}

/*
 *  Send heartbeats (Empty AppendEntryRPC)
*/
func (s *Server) sendHeartbeat() {
	if s._state.Leader() {

    appendEntry := m.NewAppendEntry(
      s._state.GetTerm(), 
      s._state.GetID(), 
      0, 
      0, 
      make([]*p.Entry, 0), 
      0)

    enc, err := appendEntry.Encode()
    
    if err != nil {
      log.Println("Error encoding", err)
    } else {
      mess := Message{
        Ty: APPEND_ENTRY,
        Payload: enc,
      }
 		  log.Println("Send heartbeat...", s._state.id)
		  s._state.StartHearthbeatTimeout()
		  s.sendAll(mess.toByte())
    }
	}
}

/*
 *  Send AppendEntry message
*/
func (s *Server) sendAppendEntryRPC() {

  prevLogIndex := len(s._state.GetEntries())-2
  prevLogTerm := s._state.GetEntries()[prevLogIndex].GetTerm()

  appendEntry := m.NewAppendEntry(
    s._state.GetTerm(), 
    s._state.GetID(), 
    uint64(prevLogIndex), 
    prevLogTerm, 
    make([]*p.Entry, 0), 
    s._state.GetCommitIndex())

  enc, err := appendEntry.Encode()
    
  if err != nil {
    log.Println("Error encoding", err)
  } else {
    mess := Message{
      Ty: APPEND_ENTRY,
      Payload: enc,
    }
	  s.sendAll(mess.toByte())
  }
}

/*
 *  Send RequestVote messages
*/
func (s *Server) sendRequestVoteRPC() {
  lastLogIndex := len(s._state.GetEntries())-1
  lastLogTerm := s._state.GetEntries()[lastLogIndex].GetTerm()
  requestVote := m.NewRequestVote(
    s._state.GetTerm(),
    s._state.GetID(),
    uint64(lastLogIndex), 
    lastLogTerm)

  enc, err := requestVote.Encode()

  if err!= nil { 
    log.Println("During encoding a request vote: ", err)
  } else {
    mess := Message {
      Ty: REQUEST_VOTE,
      Payload: enc,
    }
    s.sendAll(mess.toByte())
  }
}

func (s *Server) startElection() {
	log.Println("/////////////////////////////////////Starting new election...", s._state.id)
  s._state.IncrementTerm()
  s._state.VoteFor(s._state.GetID())
  s._state.ElectionTimeout()
  s.sendRequestVoteRPC() 
}

/*
 * Leader behaviour
*/
func (s *Server) leader() {
  log.Println("enter in leader")
	select {
	case <-s._state.HeartbeatTimeout().C:
		s.sendHeartbeat()
	}
	
}

/*
 * Followers behaviour
*/
func (s *Server) follower() {
  log.Println("enter in follower")
	select {
		case n := <-s.appendEntryChannel:
			log.Println("********** Message received ***********", n)
    s._state.StartElectionTimeout()
		case <-s._state.ElectionTimeout().C:
			if s._state.GetRole() != LEADER {
        s._state.SetRole(CANDIDATE)
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
