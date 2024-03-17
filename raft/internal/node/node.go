package node

import (
	"net"
	"raft/internal/messages"
	"raft/internal/node/address"
)

type Node interface {
	SendRpc(mex messages.Rpc) error
	ReadRpc() ([]byte,error)
	GetIp() string
	GetPort() string
    AddConnIn(conn *net.Conn) 
    AddConnOut(conn *net.Conn)
}

type node struct {
	addr address.NodeAddress
    recv net.Conn
    send net.Conn
}

// Read_rpc implements Node.
// before: func (this *node) ReadRpc() (messages.Rpc, error) 
func (this *node) ReadRpc() ([]byte, error) {
  return this.addr.Receive()
  // TODO: return messages.Rpc not []byte
}

func NewNode(remoteAddr string, remotePort string) (Node, error) {
	return &node{
		addr: address.NewNodeAddress(remoteAddr, remotePort),
	}, nil
}

func (this *node) AddConnIn(conn *net.Conn) {
  this.addr.HandleConnIn(conn)
}

func (this *node) AddConnOut(conn *net.Conn) {
  this.addr.HandleConnOut(conn)
}

func (this *node) SendRpc(mex messages.Rpc) error {
	mexByte, _ := mex.Encode()
	this.addr.Send(mexByte)
	return nil
}

func (this *node) GetIp() string {
	return this.addr.GetIp()
}

func (this *node) GetPort() string{
    return this.addr.GetPort()
}
