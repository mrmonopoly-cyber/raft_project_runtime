package node

import (
	"bufio"
	"errors"
	"fmt"
	"io"
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

	var raw_mex string = ""
	var errMex error
    if this.recv.conn == nil {
        return "", errors.New("connection not instantiated")
    }
	this.recv.mu.Lock()
    log.Println("want to read")
    log.Println("reading")
    raw_mex, errMex = bufio.NewReader(this.recv.conn).ReadString('\n')
	this.recv.mu.Unlock()

    
    if errMex == io.EOF {
        log.Println("found EOF, received message: ", raw_mex)
        return raw_mex, nil
    }

	if errMex != nil {
        log.Println("found other error, received message: ", raw_mex)
		return "", errMex
	}
    
    log.Println("found no error, received message: ", raw_mex)
	return raw_mex, errMex

}

func NewNode(remoteAddr string, remotePort string) (Node, error) {
	return &node{
		addr: address.NewNodeAddress(remoteAddr, remotePort),
	}, nil
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
	this.send.mu.Lock()
    fmt.Fprintf(this.send.conn, string(mex))
	this.send.mu.Unlock()
    return nil
	
}

func (this *node) GetIp() string {
	return this.addr.GetIp()
}

func (this *node) GetPort() string {
	return this.addr.GetPort()
}
