package rpcs

import (
	"raft/internal/raftstate"
	"raft/internal/node/nodeState"
)

type Rpc interface {
	ToString() string
	Execute(state *raftstate.State, senderState *nodeState.VolatileNodeState) *Rpc
	Encode() ([]byte, error)
	Decode(rawMex []byte) error
}
