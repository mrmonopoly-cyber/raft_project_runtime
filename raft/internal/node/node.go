package node

import (
	"raft/internal/messages"
	"raft/internal/node/address"
)

type Node interface {
	Send_rpc(mex messages.Rpc) error
	Read_rpc() (messages.Rpc,error)
	Get_ip() string
}

type node struct {
	addr address.Node_address
}

// Read_rpc implements Node.
func (this *node) Read_rpc() (messages.Rpc,error) {
	panic("unimplemented")
}

func New_node(addr string, port uint16) Node {
	return &node{
		addr: address.New_node_address(addr, port),
	}
}

func (this *node) Send_rpc(mex messages.Rpc) error {
	mex_byte, _ := mex.Encode()
	this.addr.Send(mex_byte)
	return nil
}

func (this *node) Get_ip() string {
	return this.addr.Get_ip()
}
