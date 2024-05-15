package rpcs

import (
	"raft/internal/node"
	"raft/internal/raftstate"
)

type Rpc interface {
	ToString() string
	Execute(state raftstate.State, sender node.Node) *Rpc
	Encode() ([]byte, error)
	Decode(rawMex []byte) error
}
