package ClientReq

import (
	"log"
	"raft/internal/raft_log"
	clustermetadata "raft/internal/raftstate/clusterMetadata"
	nodestate "raft/internal/raftstate/confPool/NodeIndexPool/nodeState"
	"raft/internal/rpcs"
	"raft/pkg/raft-rpcProtobuf-messages/rpcEncoding/out/protobuf"

	"google.golang.org/protobuf/proto"
)

type ClientReq struct {
    pMex protobuf.ClientReq
}


// Manage implements rpcs.Rpc.
func (this *ClientReq) Execute(
            intLog raft_log.LogEntry,
            metadata clustermetadata.ClusterMetadata,
            senderState nodestate.NodeState) rpcs.Rpc {
    var operation protobuf.Operation = (*this).pMex.Op
    var newEntries []*protobuf.LogEntry =make([]*protobuf.LogEntry, 1)
    var newLogEntry protobuf.LogEntry = protobuf.LogEntry{}
    var newLogEntryWrp []*raft_log.LogInstance = intLog.NewLogInstanceBatch(newEntries,[]func(){})

    newLogEntry.Term = metadata.GetTerm()
    newEntries[0] = &newLogEntry

    newLogEntry.OpType = operation

    newLogEntry.Description = "new " + string(operation) + " operation on file" + string((*this).pMex.Others)

    for _,v := range newLogEntryWrp {
        intLog.AppendEntryLast(v)
    }

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
