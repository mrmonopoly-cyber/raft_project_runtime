package node

import (
	"bufio"
	"fmt"
	"net"
	"raft/internal/messages"
	"raft/internal/messages/AppendEntryRPC"
	"raft/internal/messages/AppendEntryResponse"
	"raft/internal/messages/CopyStateRPC"
	"raft/internal/messages/RequestVoteRPC"
	"raft/internal/messages/RequestVoteResponse"
	"raft/internal/node/address"
	"sync"
)

type Node interface {
	SendRpc(mex messages.Rpc) error
	ReadRpc() (*messages.Rpc,error)
	GetIp() string
	GetPort() string
    AddConnIn(conn *net.Conn) 
    AddConnOut(conn *net.Conn)
}

type safeConn struct
{
    mu sync.Mutex
    conn net.Conn
}

type node struct {
	addr address.NodeAddress
    recv safeConn
    send safeConn
}

// Read_rpc implements Node.
// before: func (this *node) ReadRpc() (*messages.Rpc, error) 
func (this *node) ReadRpc() (*messages.Rpc, error) {
    
    var raw_mex string
    var errMex error 
    this.recv.mu.Lock()
    raw_mex, errMex = bufio.NewReader(this.recv.conn).ReadString('\n')
    this.recv.mu.Unlock()


    if errMex != nil {
        return nil,errMex
    }

    var mess *messages.Message = messages.NewMessage([]byte(raw_mex))

    var recRpc messages.Rpc
    switch mess.Mex_type {
    case messages.APPEND_ENTRY:
        var rpc AppendEntryRPC.AppendEntryRPC
        recRpc = &rpc
    case messages.APPEND_ENTRY_RESPONSE:
        var rpc AppendEntryResponse.AppendEntryResponse 
        recRpc = &rpc
    case messages.REQUEST_VOTE:
        var rpc RequestVoteRPC.RequestVoteRPC
        recRpc = &rpc
    case messages.REQUEST_VOTE_RESPONSE:
        var rpc RequestVoteResponse.RequestVoteResponse
        recRpc = &rpc
    case messages.COPY_STATE:
        var rpc CopyStateRPC.CopyStateRPC
        recRpc = &rpc
    default:
        panic("impossible case invalid type RPC " + string(mess.Mex_type))
    }
    recRpc.Decode(mess.Payload)

    return &recRpc,errMex
}

func NewNode(remoteAddr string, remotePort string) (Node, error) {
	return &node{
		addr: address.NewNodeAddress(remoteAddr, remotePort),
	}, nil
}

func (this *node) AddConnIn(conn *net.Conn) {
    this.send.conn = *conn
}

func (this *node) AddConnOut(conn *net.Conn) {
    this.recv.conn = *conn
}

func (this *node) SendRpc(mex messages.Rpc) error {
    var mexByte []byte
    var err error
	mexByte, err = mex.Encode()
    if err != nil {
        println("fail to send message to %v", this.GetIp())
        return err
    }
    this.send.mu.Lock()
    fmt.Fprintf(this.send.conn,string(mexByte))
    this.send.mu.Unlock()
	return nil
}

func (this *node) GetIp() string {
	return this.addr.GetIp()
}

func (this *node) GetPort() string{
    return this.addr.GetPort()
}
