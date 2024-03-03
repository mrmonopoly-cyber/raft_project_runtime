package server

import (
	"log"
	"net"
	"net/http"
	str "strconv"
	"sync"
)

type Server struct {
	_state raftStateImpl
	//log                l.Log
	connections        *sync.Map
	appendEntryChannel chan []byte
}

func NewServer(term int, id int, role Role, serversId []int) *Server {
	var server = &Server{
		*NewState(term, id, role, serversId),
		&sync.Map{},
		make(chan []byte),
	}

	return server
}

func (s *Server) Start() {

	s.ConncectToServers()

	s._state.StartElectionTimeout()

	if s._state.Leader() {
		s._state.StartHearthbeatTimeout()
	}

	for {
		go s.listen()
	}
}

func (s *Server) ConncectToServers() {
	list, err := net.Listen("tcp", "localhost:"+str.Itoa(s._state.GetID()))

	if err != nil {
		log.Println("Listener error: ", err)
	}
	go http.Serve(list, nil)

	go acceptIncomingConn(list, &s.appendEntryChannel)

	for i := 0; i < len(s._state.GetServersID()); i++ {
		var serverId = s._state.GetServersID()[i]
		if serverId != s._state.id {
			conn, err := net.Dial("tcp", "localhost:"+str.Itoa(serverId))

			if err != nil {
				log.Println("Dial error", err)
			}
			log.Println("Server ", s._state.id)
			s.connections.Store(serverId, conn)
			s.connections.Range(func(k, conn interface{}) bool {
				log.Println("...", k, conn)
				return true
			})
		}

	}

}

func acceptIncomingConn(listener net.Listener, channel *chan []byte) {
	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}
		go handleResponse(conn, channel)
	}
}

func handleResponse(conn net.Conn, channel *chan []byte) {
	buffer := make([]byte, 2048)

	for {
		_, err := conn.Read(buffer)
		if err != nil {
			log.Println("Read error: ", err)
		}

		*channel <- buffer
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

func (s *Server) listen() {
	select {
	case n := <-s.appendEntryChannel:
		log.Println("********** Message received ***********", n[0])
		s._state.StartElectionTimeout()
	case <-s._state.ElectionTimeout().C:
		if s._state.role != LEADER {
			s.startElection()
		}
	case <-s._state.HeartbeatTimeout().C:
		s.sendHeartbeat()
	}

}
