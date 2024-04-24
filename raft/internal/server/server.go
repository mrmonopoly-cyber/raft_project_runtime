package server

import (
	"log"
	"net"
    "sync"
	"raft/internal/node"
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
		otherNodes:     &sync.Map{},
        clientNodes:    &sync.Map{},
		messageChannel: make(chan pairMex),
		listener:       listener,
	}

    log.Println("number of others ip: ", len(serversIp))
	for i := 0; i < len(serversIp)-1; i++ {
		var new_node node.Node
        var nodeConn net.Conn
        var erroConn error
        var nodeId string

        nodeConn,erroConn = net.Dial("tcp",serversIp[i]+":"+port)
        if erroConn != nil {
            log.Println("Failed to connect to node: ", serversIp[i])
            continue
        }
        new_node = node.NewNode(serversIp[i], port, nodeConn)
        nodeId = generateID(serversIp[i])
        server.otherNodes.Store(nodeId, new_node)
        server._state.IncreaseNodeInCluster()
        go server.handleResponseSingleNode(nodeId, &new_node)

	}
	return server
}
