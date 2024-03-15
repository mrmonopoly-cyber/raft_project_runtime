package messages

import(
    state "raft/internal/raftstate"
)

type Rpc interface {
  GetTerm() uint64
  Encode() ([]byte, error)
  Decode(b []byte) error
  ToString() string
  Manage(state.State)
}
