package node

import (
	"net"
	"raft/internal/node/address"
	"raft/internal/node/nodeState"
)

type Node interface {
    CloseConnection()
	Send(mex []byte) error
	Recv() ([]byte, error)
	GetIp() string
	GetPort() string
    GetNodeState() *nodeState.VolatileNodeState
    ResetState(lastLogIndex int)
}


func NewNode(remoteAddr string, remotePort string, nodeConn net.Conn) (Node, error) {
  var addr address.NodeAddress
  var err error
  addr, err = address.NewNodeAddress(remoteAddr, remotePort)
	return &node{
		addr: addr,
        conn: nodeConn,
        nodeState: nodeState.NewVolatileState(),
	}, err
}

