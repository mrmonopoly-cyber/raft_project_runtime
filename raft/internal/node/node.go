package node

import (
	"net"
	"raft/internal/messages"
	"raft/internal/node/address"
	"strconv"
	"strings"
)

type Node interface {
	SendRpc(mex messages.Rpc) error
	ReadRpc() ([]byte,error)
	GetIp() string
  AddConnIn(conn *net.Conn) 
  AddConnOut(conn *net.Conn)
}

type node struct {
	addr address.NodeAddress
}

// Read_rpc implements Node.
// before: func (this *node) ReadRpc() (messages.Rpc, error) 
func (this *node) ReadRpc() ([]byte, error) {
  return this.addr.Receive()
  // TODO: return messages.Rpc not []byte
}

func NewNode(remoteAddr string) (Node, error) {
  addr := strings.Split(remoteAddr, ":")[0]
  p := strings.Split(remoteAddr, ":")[1]
  port, err := strconv.Atoi(p)
  
  if err != nil {
    return nil, err 
  }

	return &node{
		addr: address.NewNodeAddress(addr, uint16(port)),
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
