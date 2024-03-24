package node

import (
	"bufio"
	"errors"
	"log"
	"net"
	"raft/internal/node/address"
)

type Node interface {
	Send(mex []byte) error
	Recv() ([]byte, error)
	GetIp() string
	GetPort() string
	AddConn(conn net.Conn)
}

type node struct {
	addr address.NodeAddress
    conn net.Conn
}


// Read_rpc implements Node.
func (this *node) Recv() ([]byte, error) {

    const bufferSize = 1024

    var mexS string
    var outMex []byte 
    // var tempBuffer[] byte = make([]byte, bufferSize)
	var errMex error
    var byteRead int  = bufferSize

    if this.conn == nil {
        return nil, errors.New("connection not instantiated")
    }
    log.Println("want to read")
    log.Printf("start reading from %v\n", this.GetIp())
    mexS,errMex = bufio.NewReader(this.conn).ReadString('\n')
    // for byteRead == bufferSize{
    //     byteRead,errMex = this.conn.Read(tempBuffer)
    //     outMex = append(outMex, tempBuffer...)
    // }
    if errMex != nil {
        log.Println("found other error, received message: ", byteRead)
        return nil, errMex
    }   
    log.Printf("end reading from %v : %v\n", this.GetIp(), outMex)
    
    log.Println("found no error, received message: ", byteRead)
	// return outMex, errMex
	return []byte(mexS), errMex

}

func NewNode(remoteAddr string, remotePort string) (Node) {
	return &node{
		addr: address.NewNodeAddress(remoteAddr, remotePort),
	}
}

func (this *node) AddConn(conn net.Conn) {
	this.conn = conn
}

func (this *node) Send(mex []byte) error{
    if this.conn == nil {
        return errors.New("Connection with node " + this.GetIp() + " not enstablish, Dial Done?")
    }
    log.Printf("start sending message to %v", this.GetIp())
    // var mexTerm = string(mex) + "\n"
    this.conn.Write(append(mex, []byte("\n")...))
    log.Printf("message sended to %v", this.GetIp())
    return nil
	
}

func (this *node) GetIp() string {
	return this.addr.GetIp()
}

func (this *node) GetPort() string {
	return this.addr.GetPort()
}
