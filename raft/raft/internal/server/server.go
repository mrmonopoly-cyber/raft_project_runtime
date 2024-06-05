package server

import (
	"log"
	"net"
	clustermetadata "raft/internal/raftstate/clusterMetadata"
	confpool "raft/internal/raftstate/confPool"
	"sync"
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
        wg: sync.WaitGroup{},
        listener:       listener,
        clientList: sync.Map{},
        messageChannel: make(chan pairMex),
        ClusterMetadata: clustermetadata.NewClusterMetadata(ipAddPrivate,ipAddrPublic),
        ConfPool: nil,
	}
    server.ConfPool = confpool.NewConfPoll(fsRootDir,server.ClusterMetadata)

    log.Println("number of others ip: ", len(serversIp))
    log.Printf("other ips: %v\n",serversIp)
    server.connectToNodes(serversIp,port)
	return server
}
