package node

import (
	"net"
	"raft/internal/node/address"
)

type Node interface {
    CloseConnection()
	Send(mex []byte) error
	Recv() ([]byte, error)
	GetIp() string
	GetPort() string
}


func NewNode(remoteAddr string, remotePort string, nodeConn net.Conn) Node {
	var nNode *node = &node{
		addr: address.NewNodeAddress(remoteAddr, remotePort),
        conn: nodeConn,
	}
    return nNode
}

