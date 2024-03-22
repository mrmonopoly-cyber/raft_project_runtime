package node

import (
	"errors"
	"log"
	"net"
	"raft/internal/node/address"
	"sync"
)

type Node interface {
	Send(mex []byte) error
	Recv() (string, error)
	GetIp() string
	GetPort() string
	AddConnIn(conn *net.Conn)
	AddConnOut(conn *net.Conn)
}

type safeConn struct {
	mu   sync.Mutex
	conn net.Conn
}

type node struct {
	addr address.NodeAddress
	recv safeConn
	send safeConn
}


// Read_rpc implements Node.
func (this *node) Recv() (string, error) {

    const bufferSize = 1024

    var outMex = ""
	var byteRead int = -1
	var errMex error
    var buffer []byte = make([]byte, bufferSize)
    if this.recv.conn == nil {
        // return "", errors.New("connection not instantiated")
        return "", nil
    }
	this.recv.mu.Lock()
    log.Println("want to read")
    log.Println("reading")
    for byteRead == -1 || byteRead == bufferSize{
        byteRead, errMex = this.recv.conn.Read(buffer)
        if errMex != nil {
            panic("error reading")
            log.Println("found other error, received message: ", byteRead)
            return "", errMex
        }   
        outMex += string(buffer)
    }
	this.recv.mu.Unlock()
    
    log.Println("found no error, received message: ", byteRead)
	return outMex, errMex

}

func NewNode(remoteAddr string, remotePort string) (Node) {
	return &node{
		addr: address.NewNodeAddress(remoteAddr, remotePort),
	}
}

func (this *node) AddConnIn(conn *net.Conn) {
	this.recv.conn = *conn
}

func (this *node) AddConnOut(conn *net.Conn) {
	this.send.conn = *conn
}

func (this *node) Send(mex []byte) error{
    if this.send.conn == nil {
        return errors.New("Connection with node " + this.GetIp() + " not enstablish, Dial Done?")
    }
	this.recv.mu.Lock()
    this.recv.conn.Write(mex)
	this.recv.mu.Unlock()
    return nil
	
}

func (this *node) GetIp() string {
	return this.addr.GetIp()
}

func (this *node) GetPort() string {
	return this.addr.GetPort()
}
