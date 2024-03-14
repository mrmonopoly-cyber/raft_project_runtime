package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	p "raft/pkg/protobuf"
	"sync"
    append "raft/internal/messages/AppendEntryRPC"
    res_ap "raft/internal/messages/AppendEntryResponse"
    req_vo "raft/internal/messages/RequestVoteRPC"
    res_vo "raft/internal/messages/RequestVoteResponse"
    cop_st "raft/internal/messages/CopyStateRPC"
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

  wg.Add(4)
  listener := s.startListen(&wg)

  go acceptIncomingConn(listener, &wg)

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
		if serverId != s._state.GetID() {
      var conn, err = net.Dial("tcp", "localhost:" + serverId)

			for err != nil {
        log.Println("Dial error: ", err)
        log.Println("retrying...")
        conn, err = net.Dial("tcp", "localhost:" + serverId)
			}

      if conn != nil {
			  s.connections.Store(serverId, conn)
      }
		}
	}
}

func (s *Server) startListen(wg *sync.WaitGroup) net.Listener {
  listener, err := net.Listen("tcp", "localhost:" + s._state.GetID())
  
  if err != nil {
    log.Println(err) 
  }
  
  go func () {
    defer wg.Done()
    for {
      http.Serve(listener, nil)
    }
  } ()
  //go http.Serve(listener, nil)
  log.Println("Listener up on port: " + listener.Addr().String())
  return listener
}

func acceptIncomingConn(listener net.Listener, wg *sync.WaitGroup) {
  defer wg.Done()
  for {
    listener.Accept()
  }
}

/*
 * Handle incoming messages and send it to the corresponding channel
*/
func (s *Server) handleResponse(wg *sync.WaitGroup) {
  defer wg.Done()
  // iterating over the connections map and receive byte messages 
	for {
    s.connections.Range(func(k, conn interface{}) bool {
		  message, _ := bufio.NewReader(conn.(net.Conn)).ReadString('\n')
      
      // for every []byte: decode it to Message type
      if message != "" {
        var mess Message
        mess.toMessage([]byte(message))
        s.messageChannel <- &mess
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
    log.Println(string(mess))
    fmt.Fprintf(conn.(net.Conn), string(mess) + "\n")
		return true
	})
}

/*
 *  Send heartbeats (Empty AppendEntryRPC)
*/
func (s *Server) sendHeartbeat() {
  appendEntry := append.New_AppendEntryRPC (
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

/*
 *  Send AppendEntry message
*/
func (s *Server) sendAppendEntryRPC() {

  prevLogIndex := len(s._state.GetEntries())-2
  prevLogTerm := s._state.GetEntries()[prevLogIndex].GetTerm()

  appendEntry := append.New_AppendEntryRPC(
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
  requestVote := req_vo.New_RequestVoteRPC(
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
  //s._state.ElectionTimeout()
  //s.sendRequestVoteRPC() 
}

func (s *Server) More_recent_log(last_log_index uint64, last_log_term uint64) bool {
    //TODO implement More recent log
    return false
}

func (s *Server) other_node_vote_candidature(mex req_vo.RequestVoteRPC) {
    if !s._state.CanVote(){
        return
    }
    if mex.GetTerm() < s._state.GetTerm() {
        
        //TODO: send response (my_term,false)
        return
    }
    
    if !s.More_recent_log(mex.GetLastLogIndex(),mex.GetLastLogTerm()) {
        
    }else if s._state.voteFor == "" || s._state.voteFor == mex.GetCandidateId(){
        s._state.voteFor = mex.GetCandidateId()
        //TODO: send response to sender (my_term, true)
        return
    }

    //TODO: send response to server (my_term,false)
    return
    
}

/*
 * Leader behaviour
*/
func (s *Server) leader() {
	select {
    case mess := <- s.messageChannel: 
      switch mess.Ty {
        case APPEND_ENTRY:
          var appendEntry = &append.AppendEntryRPC{}
          appendEntry.Decode(mess.Payload)
          //TODO Handle append entry

        case APPEND_RESPONSE: 
          var appendResponse = &res_ap.AppendEntryResponse{}
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
	select {
  case mess := <- s.messageChannel:

    switch mess.Ty {
      case APPEND_ENTRY:
        var appendEntry = &append.AppendEntryRPC{}
        appendEntry.Decode(mess.Payload)
			  log.Println("********** Message received ***********", appendEntry)
        s._state.StartElectionTimeout()
        // TODO: Handle append entry

      case REQUEST_VOTE:
        var reqVote = &req_vo.RequestVoteRPC{}
        reqVote.Decode(mess.Payload)
        s.other_node_vote_candidature(*reqVote)
        // TODO: Handle request vote

      case COPY_STATE:
        var copyState = &cop_st.CopyStateRPC{}
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
        var appendEntry = &append.AppendEntryRPC{}
        appendEntry.Decode(mess.Payload)
        // TODO Handle response

      case REQUEST_VOTE:
        var reqVote = &req_vo.RequestVoteRPC{}
        reqVote.Decode(mess.Payload)
        // TODO Handle request vote

      case VOTE_RESPONSE:
        var voteResponse = &res_vo.RequestVoteResponse{}
        voteResponse.Decode(mess.Payload)
        // TODO Handle Response vote

      case COPY_STATE:
        var copyState = &cop_st.CopyStateRPC{}
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
