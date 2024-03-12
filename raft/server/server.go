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
	messageChannel chan *Message
}


func NewServer(term uint64, id string, role Role, serversId []string) *Server {
	var server = &Server{
    *NewState(term, id, role, serversId),
		&sync.Map{},
		make(chan *Message),
	}

  return server
}

func (s *Server) Start() {
  var wg sync.WaitGroup

  wg.Add(2)
  s.startListen()

	s.ConnectToServers()
	
  s._state.StartElectionTimeout()

	if s._state.Leader() {
		s._state.StartHearthbeatTimeout()
	}

	go s.run(&wg)

  go s.handleResponse(&wg)
  
  wg.Wait()
}

/*
 * Create a connection between this server to all the others and populate the map containing these connections 
*/
func (s *Server) ConnectToServers() {
	//go acceptIncomingConn(list, &s.appendEntryChannel)

	for i := 0; i < len(s._state.GetServersID()); i++ {
		var serverId = s._state.GetServersID()[i]
		if serverId != s._state.id {
      var conn, err = net.Dial("tcp", "localhost:" + serverId)

			for err != nil {
        log.Println("Dial error: ", err)
        log.Println("retrying...")
        conn, err = net.Dial("tcp", "localhost:" + serverId)
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
    log.Println(err) 
  }
  
  go http.Serve(listener, nil)
  log.Println("Listener up on port: " + listener.Addr().String())
}

/*
 * Handle incoming messages and send it to the corresponding channel
*/
func (s *Server) handleResponse(wg *sync.WaitGroup) {
	buffer := make([]byte, 2048)
  defer wg.Done()
  // iterating over the connections map and receive byte messages 
	for {
    s.connections.Range(func(k, conn interface{}) bool {
		  _, err := conn.(net.Conn).Read(buffer)
		  if err != nil {
			  log.Println("Read error: ", err)
		  }

      log.Println("Buffer:", buffer)
      
      // for every []byte: decode it to Message type
      var mess Message
      mess.toMessage(buffer)

      s.messageChannel <- &mess
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
  log.Println("enter leader...")
	select {
    case mess := <- s.messageChannel: 
      switch mess.Ty {
        case APPEND_ENTRY:
          var appendEntry = &m.AppendEntryRPC{}
          appendEntry.Decode(mess.Payload)
          //TODO Handle append entry

        case APPEND_RESPONSE: 
          var appendResponse = &m.AppendEntryResponse{}
          appendResponse.Decode(mess.Payload)
          // TODO Handle Response

        default:
          // Do nothing
       }
	  case <-s._state.HeartbeatTimeout().C:
		  s.sendHeartbeat()
	}	
}

/*
 * Followers behaviour
*/
func (s *Server) follower() {
  log.Println("Enter follower....")
	select {
  case mess := <- s.messageChannel:

    switch mess.Ty {
      case APPEND_ENTRY:
        var appendEntry = &m.AppendEntryRPC{}
        appendEntry.Decode(mess.Payload)
			  log.Println("********** Message received ***********", appendEntry)
        s._state.StartElectionTimeout()
        // TODO: Handle append entry

      case REQUEST_VOTE:
        var reqVote = &m.RequestVoteRPC{}
        reqVote.Decode(mess.Payload)
        // TODO: Handle request vote

      case COPY_STATE:
        var copyState = &m.CopyStateRPC{}
        copyState.Decode(mess.Payload)
        // TODO: Handle copy state message 
      
      default: 
        // default behaviour: do nothing
    }
	case <-s._state.ElectionTimeout().C:
		if s._state.GetRole() != LEADER {
      s._state.SetRole(CANDIDATE)
			s.startElection()
		}
  case <- s.messageChannel:
  }
}

/*
 * Handle candidates
*/
func (s *Server) candidate() {
  select {
  case mess := <- s.messageChannel:
    switch mess.Ty {
      case APPEND_ENTRY:
        var appendEntry = &m.AppendEntryRPC{}
        appendEntry.Decode(mess.Payload)
        // TODO Handle response

      case REQUEST_VOTE:
        var reqVote = &m.RequestVoteRPC{}
        reqVote.Decode(mess.Payload)
        // TODO Handle request vote

      case VOTE_RESPONSE:
        var voteResponse = &m.RequestVoteResponse{}
        voteResponse.Decode(mess.Payload)
        // TODO Handle Response vote

      case COPY_STATE:
        var copyState = &m.CopyStateRPC{}
        copyState.Decode(mess.Payload)
        // TODO Handle copy state 

      default:
        // default behaviour: do nothing
    }
  case <-s._state.ElectionTimeout().C:
    // TODO Handle election timeout during an election phase
  }
}

func (s *Server) run(wg *sync.WaitGroup) {
  defer wg.Done()
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
