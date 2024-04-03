package node

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net"
	"raft/internal/node/address"
	"raft/internal/node/nodeState"
)

type Node interface {
	Send(mex []byte) error
	Recv() ([]byte, error)
	GetIp() string
	GetPort() string
    GetNodeState() *nodeState.VolatileNodeState
    ResetState()
}

type node struct {
	addr address.NodeAddress
    conn net.Conn
    nodeState nodeState.VolatileNodeState
}

func NewNode(remoteAddr string, remotePort string, nodeConn net.Conn) (Node) {
	return &node{
		addr: address.NewNodeAddress(remoteAddr, remotePort),
        conn: nodeConn,
        nodeState: nodeState.NewVolatileState(),
	}
}

// Read_rpc implements Node.
func (this *node) Recv() ([]byte, error) {

    buffer := &bytes.Buffer{}

	// Create a temporary buffer to store incoming data
	tmp := make([]byte, 1024) // Initial buffer size


    if this.conn == nil {
        return nil, errors.New("connection not instantiated")
    }
    //log.Println("want to read")
    //log.Printf("start reading from %v\n", this.GetIp())

    var bytesRead int = len(tmp)
    var errConn error
    var errSavi error
    for bytesRead == len(tmp){
		// Read data from the connection
		bytesRead, errConn = this.conn.Read(tmp)

        // Write the read data into the buffer
        _, errSavi = buffer.Write(tmp[:bytesRead])
        
        // check error saving
        if errSavi != nil {
            return nil, errSavi
        }

		if errConn != nil {
			if errConn != io.EOF {
				// Handle other errConnors
				return nil, errConn
			}
			break
		}
	}

    //log.Printf("end reading from %v : %v\n", this.GetIp(), buffer)
    
    //log.Println("found no error, received message: ", buffer)
	return buffer.Bytes(), nil 

}

func (this *node)GetNodeState() *nodeState.VolatileNodeState{
    return &this.nodeState   
}

func (this *node) AddConn(conn net.Conn) {
	this.conn = conn
}

func (this *node) Send(mex []byte) error{
    if this.conn == nil {
        return errors.New("Connection with node " + this.GetIp() + " not enstablish, Dial Done?")
    }
    log.Printf("start sending message to %v", this.GetIp())
    this.conn.Write(mex)
    log.Printf("message sended to %v", this.GetIp())
    return nil
	
}

func (this *node) GetIp() string {
	return this.addr.GetIp()
}

func (this *node) GetPort() string {
	return this.addr.GetPort()
}
func (this *node) ResetState(){
    (*this).nodeState.InitVolatileState()
}
