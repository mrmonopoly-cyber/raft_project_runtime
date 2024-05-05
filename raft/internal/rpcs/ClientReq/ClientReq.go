package ClientReq

import (
	"log"
	"raft/internal/node/nodeState"
	"raft/internal/raftstate"
	"raft/internal/rpcs"
	"raft/internal/rpcs/ClientResponse"
	"raft/pkg/raft-rpcProtobuf-messages/rpcEncoding/out/protobuf"

	"google.golang.org/protobuf/proto"
)

type ClientReq struct {
    pMex protobuf.ClientReq
}


// Manage implements rpcs.Rpc.
func (this *ClientReq) Execute(state *raftstate.State, senderState *nodeState.VolatileNodeState) *rpcs.Rpc {
    var operation protobuf.Operation = (*this).pMex.Op
    var newEntries []*protobuf.LogEntry =make([]*protobuf.LogEntry, 1)
    var newLogEntry protobuf.LogEntry
    var op string = "NULL"
    var fileName string = string(this.pMex.GetFileName())
    var response []byte = nil

    newLogEntry.Term = (*state).GetTerm()
    newLogEntry.FilenName = fileName
    newEntries[0] = &newLogEntry

    switch operation{
    case protobuf.Operation_READ:
        newLogEntry.OpType = protobuf.Operation_READ
        op = "READ"
        // retreive file
        response = nil 
    case protobuf.Operation_WRITE:
        newLogEntry.OpType = protobuf.Operation_WRITE
        newLogEntry.Payload = this.pMex.Others
        op = "WRITE"
    case protobuf.Operation_DELETE:
        newLogEntry.OpType = protobuf.Operation_DELETE
        op = "DELETE"
    case protobuf.Operation_RENAME:
        newLogEntry.OpType = protobuf.Operation_RENAME
        newLogEntry.Payload = this.pMex.Others
        op = "RENAME"
    case protobuf.Operation_CREATE:
        newLogEntry.OpType = protobuf.Operation_CREATE
        newLogEntry.Payload = this.pMex.Others
        op = "CREATE"
    default:
        log.Printf("NOT IMPLMENTED OPERATION %v\n", operation)
        return nil
    }

    newLogEntry.Description = "new " + op + " operation on file" + fileName

    (*state).AppendEntries(newEntries,(*state).GetLastLogIndex()+1)

    // Create client response and return it
    var clientReponse rpcs.Rpc = ClientResponse.NewClientResponseRPC(response)
    return &clientReponse
}

// ToString implements rpcs.Rpc.
func (this *ClientReq) ToString() string {
    return "{" + (*this).pMex.GetOp().String() + "}"
}

func (this *ClientReq) Encode() ([]byte, error) {
    return proto.Marshal(&(*this).pMex)
}
func (this *ClientReq) Decode(b []byte) error {
	return proto.Unmarshal(b,&this.pMex)
}
