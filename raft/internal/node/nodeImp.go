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

type node struct {
	addr address.NodeAddress
    conn net.Conn
    nodeState nodeState.VolatileNodeState
}

// Read_rpc implements Node.
func (this *node) Recv() ([]byte, error) {

    buffer := &bytes.Buffer{}

	// Create a temporary buffer to store incoming data
	tmp := make([]byte, 1024) // Initial buffer size


    if this.conn == nil {
        return nil, errors.New("connection not instantiated")
    }

    var bytesRead int = len(tmp)
    var errConn error
    var errSavi error
    for bytesRead == len(tmp){
		bytesRead, errConn = this.conn.Read(tmp)
        _, errSavi = buffer.Write(tmp[:bytesRead])
        if errSavi != nil {
            return nil, errSavi
        }

		if errConn != nil {
			if errConn != io.EOF {
				return nil, errConn
			}
            if errConn == io.EOF {
                return nil, errConn
            }
			break
		}
	}
	return buffer.Bytes(), nil 
}

func (this *node)GetNodeState() nodeState.VolatileNodeState{
    return this.nodeState   
}

func (this *node) CloseConnection(){
    (*this).conn.Close()
}

func (this *node) AddConn(conn net.Conn) {
	this.conn = conn
}

func (this *node) Send(mex []byte) error{
    if this.conn == nil {
        return errors.New("Connection with node " + this.GetIp() + " not enstablish, Dial Done?")
    }
    var _,err = this.conn.Write(mex)
    if err != nil {
        log.Panicf("error sending to %v, error %v\n", (*this).GetIp(), err)
        return err
    }
    return nil
	
}

func (this *node) GetIp() string {
	return this.addr.GetIp()
}

func (this *node) GetPort() string {
	return this.addr.GetPort()
}
func (this *node) ResetState(lastLogIndex int){
    (*this).nodeState.InitVolatileState(lastLogIndex)
}
