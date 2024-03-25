package NEW_RPC

import (
	"raft/internal/rpcs"
	"raft/internal/raftstate"
	// "sync"
 //    "raft/pkg/protobuf/t/protobuf"
	//
	// "google.golang.org/protobuf/proto"
)

type NEW_RPC struct {
}

func NewNEW_RPCRPC(term uint64) rpcs.Rpc {
    return &NEW_RPC{
    }
}

// GetId implements rpcs.Rpc.
func (this *NEW_RPC) GetId() string {
    panic("dummy implementation")
}


// Manage implements rpcs.Rpc.
func (this *NEW_RPC) Execute(state *raftstate.State) *rpcs.Rpc {
    panic("dummy implementation")
}

// ToString implements rpcs.Rpc.
func (this *NEW_RPC) ToString() string {
    panic("dummy implementation")
}

func (this *NEW_RPC) GetTerm() uint64 {
    panic("dummy implementation")
}

func (this *NEW_RPC) Encode() ([]byte, error) {
    panic("dummy implementation")
}
func (this *NEW_RPC) Decode(b []byte) error {
    panic("dummy implementation")
}
