package UpdateNode

import (
	"log"
	"raft/internal/raftstate"
	"raft/internal/rpcs"
	"raft/internal/node/nodeState"
    "raft/pkg/raft-rpcProtobuf-messages/rpcEncoding/out/protobuf"

	"google.golang.org/protobuf/proto"
)

type UpdateNode struct {
    pMex protobuf.UpdateNode
}

func NewUpdateNodeRPC(voteAble bool, log *protobuf.LogEntry) rpcs.Rpc {
    return &UpdateNode{
        pMex: protobuf.UpdateNode{
            Votante: false,
            Log: log,
        },
    }
}

// Manage implements rpcs.Rpc.
func (this *UpdateNode) Execute(state *raftstate.State, senderState *nodeState.VolatileNodeState) *rpcs.Rpc {
    log.Printf("updating log entry with new entry %v\n",this.pMex.Log)
    (*state).VoteRight(this.pMex.Votante)
    if this.pMex.Log != nil {
        (*state).AppendEntries([]*protobuf.LogEntry{this.pMex.Log},int((*state).GetCommitIndex()+1))
    }

    return nil
}

// ToString implements rpcs.Rpc.
func (this *UpdateNode) ToString() string {
    panic("dummy implementation")
}

func (this *UpdateNode) Encode() ([]byte, error) {
    var mess []byte
    var err error

    mess, err = proto.Marshal(&(*this).pMex)
    if err != nil {
        log.Panicln("error in Encoding Request Vote: ", err)
    }

	return mess, err
}
func (this *UpdateNode) Decode(b []byte) error {
	err := proto.Unmarshal(b,&this.pMex)
    if err != nil {
        log.Panicln("error in Encoding Request Vote: ", err)
    }
	return err
}
