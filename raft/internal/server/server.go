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

func NewServer(ipAddPrivate string, ipAddrPublic string, port string, serversIp []string, fsRootDir string) Server {
	listener, err := net.Listen("tcp",":"+port)

	if err != nil {
		log.Fatalf("Failed to listen on port %s: %s", port, err)
	}

    log.Printf("my ip are: %v, %v\n",  ipAddPrivate, ipAddrPublic)

	var server = &server{
		_state:         state.NewState(ipAddPrivate, ipAddrPublic, fsRootDir),
		unstableNodes:    &sync.Map{},
        clientNodes:    &sync.Map{},
		messageChannel: make(chan pairMex),
		listener:       listener,
	}

    log.Println("number of others ip: ", len(serversIp))
    log.Printf("other ips: %v\n",serversIp)
    // server.connectToNodes(serversIp,port)
	return server
}
