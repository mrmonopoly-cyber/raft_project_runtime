package NEW_RPC

import (
	"log"
	"raft/internal/raftstate"
	"raft/internal/rpcs"

	"google.golang.org/protobuf/proto"
)

type NEW_RPC struct {
}

func NewNEW_RPCRPC(term uint64) rpcs.Rpc {
    return &NEW_RPC{
    }
}

// Manage implements rpcs.Rpc.
func (this *NEW_RPC) Execute(state *raftstate.State) *rpcs.Rpc {
    panic("dummy implementation")
}

// ToString implements rpcs.Rpc.
func (this *NEW_RPC) ToString() string {
    panic("dummy implementation")
}

func (this *NEW_RPC) Encode() ([]byte, error) {
    var mess []byte
    var err error

    mess, err = proto.Marshal(&(*this).pMex)
    if err != nil {
        log.Panicln("error in Encoding Request Vote: ", err)
    }

	return mess, err
}
func (this *NEW_RPC) Decode(b []byte) error {
	err := proto.Unmarshal(b,&this.pMex)
    if err != nil {
        log.Panicln("error in Encoding Request Vote: ", err)
    }
	return err
}
