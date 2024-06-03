package ClientReq

import (
	"log"
	"raft/internal/node"
	"raft/internal/raft_log"
	"raft/internal/raftstate"
	"raft/internal/rpcs"
	"raft/pkg/raft-rpcProtobuf-messages/rpcEncoding/out/protobuf"

	"google.golang.org/protobuf/proto"
)

type ClientReq struct {
    pMex protobuf.ClientReq
}


// Manage implements rpcs.Rpc.
func (this *ClientReq) Execute(state raftstate.State, sender node.Node) rpcs.Rpc {
    var operation protobuf.Operation = (*this).pMex.Op
    var newEntries []*protobuf.LogEntry =make([]*protobuf.LogEntry, 1)
    var newLogEntry protobuf.LogEntry = protobuf.LogEntry{}
    var newLogEntryWrp []*raft_log.LogInstance = state.NewLogInstanceBatch(newEntries,[]func(){})

    newLogEntry.Term = state.GetTerm()
    newEntries[0] = &newLogEntry

    newLogEntry.OpType = operation

    newLogEntry.Description = "new " + string(operation) + " operation on file" + string((*this).pMex.Others)

    state.AppendEntries(newLogEntryWrp)

    return nil
}

// ToString implements rpcs.Rpc.
func (this *ClientReq) ToString() string {
    return "{" + (*this).pMex.GetOp().String() + "}"
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
