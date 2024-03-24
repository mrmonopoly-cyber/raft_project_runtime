package messages

import (
	"raft/internal/raftstate"
)

type Rpc interface {
	GetId() string
	GetTerm() uint64
	Decode(b []byte) error
	ToString() string
	Execute(state *raftstate.State) *Rpc
    Encode() ([]byte, error)
}

func Decode(m string) *Rpc{
    panic("unimplemented")
}
