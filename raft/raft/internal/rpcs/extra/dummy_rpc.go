package NEW_RPC

import (
	"raft/internal/messages"
	"raft/internal/raftstate"
	// "sync"
 //    "raft/pkg/protobuf/t/protobuf"
	//
	// "google.golang.org/protobuf/proto"
)

type NEW_RPC struct {
}

func NewNEW_RPCRPC(term uint64) messages.Rpc {
    return &NEW_RPC{
    }
}

// GetId implements messages.Rpc.
func (this *NEW_RPC) GetId() string {
    panic("dummy implementation")
}


// Manage implements messages.Rpc.
func (this *NEW_RPC) Execute(state *raftstate.State) *messages.Rpc {
    panic("dummy implementation")
}

// ToString implements messages.Rpc.
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
