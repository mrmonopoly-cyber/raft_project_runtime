package messages

import (
	"raft/internal/raftstate"
	"sync"
)

type Rpc interface {
  GetTerm() uint64
  Encode() ([]byte, error)
  Decode(b []byte) error
  ToString() string
  Execute(n *sync.Map, state raftstate.State) 
}


