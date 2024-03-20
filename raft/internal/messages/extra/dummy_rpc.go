package NEW_RPC

import (
	"raft/internal/messages"
	"raft/internal/raftstate"
	p "raft/pkg/protobuf"
	"sync"

	"google.golang.org/protobuf/proto"
)

type NEW_RPC struct {
	term uint64
}

func NewNEW_RPCRPC(term uint64) messages.Rpc {
    return &NEW_RPC{
        term: 0,
    }
}

// GetId implements messages.Rpc.
func (this *NEW_RPC) GetId() string {
	panic("unimplemented")
}


// Manage implements messages.Rpc.
func (this *NEW_RPC) Execute(state *raftstate.State) *messages.Rpc {
	panic("unimplemented")
}

// ToString implements messages.Rpc.
func (this *NEW_RPC) ToString() string {
	panic("unimplemented")
}

func (this NEW_RPC) GetTerm() uint64 {
	return this.term
}

func (this NEW_RPC) Encode() ([]byte, error) {
	appendEntry := &p.NEW_RPC{
		Term: proto.Uint64(this.term),
	}

	mess, err := proto.Marshal(appendEntry)
	return mess, err
}
func (this NEW_RPC) Decode(b []byte) error {
	pb := new(p.NEW_RPC)
	err := proto.Unmarshal(b, pb)

	if err != nil {
		this.term = pb.GetTerm()
	}

	return err
}
