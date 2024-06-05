package clientReturnValue

import (
	"log"
	"raft/internal/node"
	"raft/internal/raft_log"
	clustermetadata "raft/internal/raftstate/clusterMetadata"
	"raft/internal/rpcs"
	"raft/pkg/raft-rpcProtobuf-messages/rpcEncoding/out/protobuf"

	"google.golang.org/protobuf/proto"
)

type ClientReturnValue struct {
    pMex protobuf.ClientReturnValue
}

func NewclientReturnValueRPC(exitStatus protobuf.STATUS, description string) rpcs.Rpc {
    return &ClientReturnValue{
        pMex: protobuf.ClientReturnValue{
            ExitStatus: exitStatus,
            Description: description,
        },
    }
}

// Manage implements rpcs.Rpc.
func (this *ClientReturnValue) Execute(
            intLog raft_log.LogEntry,
            metadata clustermetadata.ClusterMetadata,
            sender node.Node) rpcs.Rpc {
    panic("dummy implementation")
}

// ToString implements rpcs.Rpc.
func (this *ClientReturnValue) ToString() string {
    return this.pMex.String()
}

func (this *ClientReturnValue) Encode() ([]byte, error) {
    var mess []byte
    var err error

    mess, err = proto.Marshal(&(*this).pMex)
    if err != nil {
        log.Panicln("error in Encoding Request Vote: ", err)
    }

	return mess, err
}
func (this *ClientReturnValue) Decode(b []byte) error {
	err := proto.Unmarshal(b,&this.pMex)
    if err != nil {
        log.Panicln("error in Encoding Request Vote: ", err)
    }
	return err
}
