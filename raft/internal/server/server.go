package server

import (
	"log"
	"net"
    "sync"
	state "raft/internal/raftstate"
)


type Server interface{
    Start()
}

func NewServer(term uint64, ipAddPrivate string, ipAddrPublic string, port string, serversIp []string, fsRootDir string) Server {
	listener, err := net.Listen("tcp",":"+port)

	if err != nil {
		log.Fatalf("Failed to listen on port %s: %s", port, err)
	}

    log.Printf("my ip are: %v, %v\n",  ipAddPrivate, ipAddrPublic)

	var server = &server{
		_state:         state.NewState(term, ipAddPrivate, ipAddrPublic, state.FOLLOWER, fsRootDir),
		unstableNodes:    &sync.Map{},
		stableNodes:    &sync.Map{},
        clientNodes:    &sync.Map{},
		messageChannel: make(chan pairMex),
		listener:       listener,
	}

    log.Println("number of others ip: ", len(serversIp))
    server.connectToNodes(serversIp,port)
	return server
}
