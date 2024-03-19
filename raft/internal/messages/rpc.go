package messages

import (
	"raft/internal/raftstate"
)

type Rpc interface {
  GetId() string
  GetTerm() uint64
  Encode() ([]byte, error)
  Decode(b []byte) error
  ToString() string
  Execute(state *raftstate.State, resp *Rpc) 
}


