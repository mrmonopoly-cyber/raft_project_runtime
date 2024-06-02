package ClientReturnValue

import (
	"log"
	"raft/internal/node"
	"raft/internal/raftstate"
	"raft/internal/rpcs"
	"raft/pkg/raft-rpcProtobuf-messages/rpcEncoding/out/protobuf"

	"google.golang.org/protobuf/proto"
)

type clientReturnValue struct {
    pMex protobuf.ClientReturnValue
}

func NewclientReturnValueRPC(exitStatus protobuf.STATUS, description string) rpcs.Rpc {
    return &clientReturnValue{
        pMex: protobuf.ClientReturnValue{
            ExitStatus: exitStatus,
            Description: description,
        },
    }
}

// Manage implements rpcs.Rpc.
func (this *clientReturnValue) Execute(state raftstate.State, sender node.Node) rpcs.Rpc {
    panic("dummy implementation")
}

// ToString implements rpcs.Rpc.
func (this *clientReturnValue) ToString() string {
    return this.pMex.String()
}

func (this *clientReturnValue) Encode() ([]byte, error) {
    var mess []byte
    var err error

    mess, err = proto.Marshal(&(*this).pMex)
    if err != nil {
        log.Panicln("error in Encoding Request Vote: ", err)
    }

	return mess, err
}
func (this *clientReturnValue) Decode(b []byte) error {
	err := proto.Unmarshal(b,&this.pMex)
    if err != nil {
        log.Panicln("error in Encoding Request Vote: ", err)
    }
	return err
}
