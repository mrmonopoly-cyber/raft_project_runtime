package messages

import (
	"raft/internal/raftstate"
)

type Rpc interface {
	GetId() string
	GetTerm() uint64
	ToString() string
	Execute(state *raftstate.State) *Rpc
    Encode() ([]byte, error)
}

func Decode(b []byte) *Rpc{
    panic("unimplemented")
}
