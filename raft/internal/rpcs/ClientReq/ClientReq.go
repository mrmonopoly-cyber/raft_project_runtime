package ClientReq

import (
	"log"
	"raft/internal/raftstate"
	"raft/internal/rpcs"
	"raft/internal/node/nodeState"
    "raft/pkg/rpcEncoding/out/protobuf"

	"google.golang.org/protobuf/proto"
)

type ClientReq struct {
    pMex protobuf.ClientReq
}


// Manage implements rpcs.Rpc.
func (this *ClientReq) Execute(state *raftstate.State, senderState *nodeState.VolatileNodeState) *rpcs.Rpc {
    var operation protobuf.ClientOperation = (*this).pMex.Op
    var newLogEntryArr []*protobuf.LogEntry = make([]*protobuf.LogEntry, 1)
    var newLogEntry protobuf.LogEntry
    var op string = "NULL"

    newLogEntryArr[0] = &newLogEntry
    *newLogEntry.Term =  (*state).GetTerm()

    switch operation{
    case protobuf.ClientOperation_READ:
        log.Printf("testing operation READ, TO IMPLEMENT")
        *newLogEntry.OpType = protobuf.Operation_READ
        op = "READ"
    case protobuf.ClientOperation_WRTE:
        log.Printf("testing operation WRITE, TO IMPLEMENT")
        *newLogEntry.OpType = protobuf.Operation_WRITE
        op = "WRITE"
    case protobuf.ClientOperation_DELETE:
        log.Printf("testing operation DELETE, TO IMPLEMENT")
        *newLogEntry.OpType = protobuf.Operation_DELETE
        op = "DELETE"
    default:
        log.Printf("NOT IMPLMENTED OPERATION %v\n", operation)
        return nil
    }

    *newLogEntry.Description = "new " + op + " operation on file" + string((*this).pMex.Payload)

    (*state).AppendEntries(newLogEntryArr,0)

    return nil
}

// ToString implements rpcs.Rpc.
func (this *ClientReq) ToString() string {
    panic("dummy implementation")
}

func (this *ClientReq) Encode() ([]byte, error) {
    var mess []byte
    var err error

    mess, err = proto.Marshal(&(*this).pMex)
    if err != nil {
        log.Panicln("error in Encoding Request Vote: ", err)
    }

	return mess, err
}
func (this *ClientReq) Decode(b []byte) error {
	err := proto.Unmarshal(b,&this.pMex)
    if err != nil {
        log.Panicln("error in Encoding Request Vote: ", err)
    }
	return err
}
