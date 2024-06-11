package ClientReq

import (
	"log"
	"raft/internal/raft_log"
	clustermetadata "raft/internal/raftstate/clusterMetadata"
	nodestate "raft/internal/raftstate/confPool/NodeIndexPool/nodeState"
	confmetadata "raft/internal/raftstate/confPool/singleConf/confMetadata"
	"raft/internal/rpcs"
	"raft/internal/rpcs/clientReturnValue"
	"raft/internal/utiliy"
	"raft/pkg/raft-rpcProtobuf-messages/rpcEncoding/out/protobuf"

	"google.golang.org/protobuf/proto"
)

type ClientReq struct {
    pMex protobuf.ClientReq
    returnValue chan []byte
}


// Manage implements rpcs.Rpc.
func (this *ClientReq) Execute(
            intLog raft_log.LogEntry,
            metadata clustermetadata.ClusterMetadata,
            confMetadata confmetadata.ConfMetadata,
            senderState nodestate.NodeState) rpcs.Rpc {
    var operation protobuf.Operation = (*this).pMex.Op
    var newLogEntry protobuf.LogEntry = protobuf.LogEntry{
        Term: metadata.GetTerm(),
        OpType: operation,
        Description:    "new " + 
                        string(operation) + 
                        " operation on file" + 
                        string((*this).pMex.Others)}
    var newLogEntryWrp = intLog.NewLogInstanceBatch([]*protobuf.LogEntry{&newLogEntry})
    newLogEntryWrp[0].ReturnValue = make(chan utiliy.Pair[[]byte, error])

    intLog.AppendEntry(newLogEntryWrp,-2)

    var returnValue = <- newLogEntryWrp[0].ReturnValue

    var exitStatus = returnValue.Snd
    var retValue = returnValue.Fst
    var exitReturnValue protobuf.STATUS = protobuf.STATUS_SUCCESS

    if exitStatus != nil{
        exitReturnValue = protobuf.STATUS_FAILURE
    }

    return clientReturnValue.NewclientReturnValueRPC(exitReturnValue, retValue, exitStatus.Error())
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
    this.returnValue = make(chan []byte)
	return err
}
