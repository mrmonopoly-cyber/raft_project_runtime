package node

import (
	"bufio"
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
	AddConn(conn *net.Conn)
}

type safeConn struct {
	mu   sync.Mutex
	conn *net.Conn
}

type node struct {
	addr address.NodeAddress
	safeConn safeConn
}


// Read_rpc implements Node.
func (this *node) Recv() (string, error) {

    const bufferSize = 1024

    var outMex = ""
	var byteRead string
	var errMex error
    if this.safeConn.conn == nil {
        return "", errors.New("connection not instantiated")
    }
	this.safeConn.mu.Lock()
    log.Println("want to read")
    log.Println("reading")
    outMex,errMex = bufio.NewReader(*this.safeConn.conn).ReadString('\n')
    if errMex != nil {
        log.Println("found other error, received message: ", byteRead)
        return "", errMex
    }   
	this.safeConn.mu.Unlock()
    
    log.Println("found no error, received message: ", byteRead)
	return outMex, errMex

}

func NewNode(remoteAddr string, remotePort string) (Node) {
	return &node{
		addr: address.NewNodeAddress(remoteAddr, remotePort),
	}
}

func (this *node) AddConn(conn *net.Conn) {
	this.safeConn.conn = conn
}

func (this *node) Send(mex []byte) error{
    if this.safeConn.conn == nil {
        return errors.New("Connection with node " + this.GetIp() + " not enstablish, Dial Done?")
    }
	this.safeConn.mu.Lock()
    (*this.safeConn.conn).Write(mex)
	this.safeConn.mu.Unlock()
    return nil
	
}

func (this *node) GetIp() string {
	return this.addr.GetIp()
}

func (this *node) GetPort() string {
	return this.addr.GetPort()
}
