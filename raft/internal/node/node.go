package node

import (
	"net"
	"raft/internal/node/address"
	"raft/internal/node/nodeState"
	nodematchidx "raft/internal/raftstate/nodeMatchIdx"
)

type Node interface {
    CloseConnection()
	Send(mex []byte) error
	Recv() ([]byte, error)
	GetIp() string
	GetPort() string

    nodeState.VolatileNodeState
}


func NewNode(remoteAddr string, remotePort string, nodeConn net.Conn, statePool nodematchidx.NodeCommonMatch) Node {
	var nNode *node = &node{
		addr: address.NewNodeAddress(remoteAddr, remotePort),
        conn: nodeConn,
        volPrivState: volPrivState{
            statepool: statePool,
        },
	}
    nNode.statepool.AddNode(remoteAddr)
    return nNode
}

