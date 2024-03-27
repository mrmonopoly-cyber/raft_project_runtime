package rpcs

import (
	"raft/internal/raftstate"
)

type Rpc interface {
	ToString() string
	Execute(state *raftstate.State) *Rpc
    Encode() ([]byte, error)
    Decode(rawMex []byte) (error)
}

