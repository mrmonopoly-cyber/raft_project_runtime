package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
	"raft/internal/messages"
	append "raft/internal/messages/AppendEntryRPC"
	res_ap "raft/internal/messages/AppendEntryResponse"
	cop_st "raft/internal/messages/CopyStateRPC"
	req_vo "raft/internal/messages/RequestVoteRPC"
	res_vo "raft/internal/messages/RequestVoteResponse"
	"raft/internal/node"
	state "raft/internal/raftstate"
	p "raft/pkg/protobuf"
	"sync"
)




type Server struct {
	_state state.State
	connections     *sync.Map
    other_nodes     *node.Node 
    messageChannel chan messages.Rpc
}


func NewServer(term uint64, id string, role state.Role, serversId []string) *Server {
	var server = &Server{
    state.NewState(term, id, role, serversId),
		&sync.Map{},
        nil,
        nil,
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

          // for every []byte: decode it to RPC
          if message != "" {
              var mess *messages.Message = messages.New_message([]byte(message))
              var rec_rpc messages.Rpc
              switch {
              case mess.Mex_type == messages.APPEND_ENTRY:
                  var rpc append.AppendEntryRPC
                  rec_rpc = &rpc
              case mess.Mex_type == messages.APPEND_ENTRY_RESPONSE:
                  var rpc res_ap.AppendEntryResponse
                  rec_rpc = &rpc
              case mess.Mex_type == messages.REQUEST_VOTE:
                  var rpc req_vo.RequestVoteRPC
                  rec_rpc = &rpc
              case mess.Mex_type == messages.REQUEST_VOTE_RESPONSE:
                  var rpc res_vo.RequestVoteResponse
                  rec_rpc = &rpc
              case mess.Mex_type == messages.COPY_STATE:
                  var rpc cop_st.CopyStateRPC
                  rec_rpc = &rpc
              default:
                  panic("impossible case")
              }
              rec_rpc.Decode(mess.Payload)
              s.messageChannel <- rec_rpc
          }

          return true
      })
  }
}

/*
 *  Send to every connected server a bunch of bytes
*/
func (s *Server) sendAll(mex messages.Rpc) {

	s.connections.Range(func(k, conn interface{}) bool {
    log.Println(mex.ToString())
    fmt.Fprintf(conn.(net.Conn), string(mex.ToMessage().ToByte()) + "\n")
		return true
	})
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


    s.sendAll(appendEntry)
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
        case mess := <- s.messageChannel: 
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
        hearthbit := append.New_AppendEntryRPC (
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
    case mess := <- s.messageChannel:

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
  case mess := <- s.messageChannel:
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

func (s *Server) run(wg *sync.WaitGroup) {
  defer wg.Done()
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
