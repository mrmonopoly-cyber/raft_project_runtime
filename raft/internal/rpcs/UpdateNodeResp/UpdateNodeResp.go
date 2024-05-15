package UpdateNodeResp

import (
	"log"
	"raft/internal/node"
	"raft/internal/raftstate"
	"raft/internal/rpcs"
	"raft/pkg/raft-rpcProtobuf-messages/rpcEncoding/out/protobuf"

	"google.golang.org/protobuf/proto"
)

type UpdateNodeResp struct {
    pMex protobuf.UpdateNodeResp
}

func NewUpdateNodeRespRPC() rpcs.Rpc {
    return &UpdateNodeResp{
        pMex: protobuf.UpdateNodeResp{
            Ack: true,
        },
    }
}

// Manage implements rpcs.Rpc.
func (this *UpdateNodeResp) Execute(state raftstate.State, sender node.Node) *rpcs.Rpc {
    sender.NextIndexStep()
    return nil
}

// ToString implements rpcs.Rpc.
func (this *UpdateNodeResp) ToString() string {
    if this.pMex.Ack {
        return "TRUE"
    }
    return "FALSE"
}

func (this *UpdateNodeResp) Encode() ([]byte, error) {
    var mess []byte
    var err error

    mess, err = proto.Marshal(&(*this).pMex)
    if err != nil {
        log.Panicln("error in Encoding Request Vote: ", err)
    }

	return mess, err
}
func (this *UpdateNodeResp) Decode(b []byte) error {
	err := proto.Unmarshal(b,&this.pMex)
    if err != nil {
        log.Panicln("error in Encoding Request Vote: ", err)
    }
	return err
}
