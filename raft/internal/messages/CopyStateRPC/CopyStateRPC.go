package CopyStateRPC

import (
	"google.golang.org/protobuf/proto"
	"raft/internal/messages"
	p "raft/pkg/protobuf"
)

type CopyStateRPC struct {
	term    uint64
	index   uint64
	voting  bool
	entries []*p.Entry
}

// ToMessage implements messages.Rpc.
func (this *CopyStateRPC) ToMessage() messages.Message {
	panic("unimplemented")
}

// ToString implements messages.Rpc.
func (this *CopyStateRPC) ToString() string {
	panic("unimplemented")
}

func new_CopyStateRPC(term uint64, index uint64, voting bool,
	entries []*p.Entry) messages.Rpc {
	return &CopyStateRPC{
		term:    term,
		index:   index,
		voting:  voting,
		entries: entries,
	}
}

func (this CopyStateRPC) GetTerm() uint64 {
	return this.term
}
func (this CopyStateRPC) Encode() ([]byte, error) {

	copyState := &p.CopyState{
		Term:    proto.Uint64(this.term),
		Voting:  proto.Bool(this.voting),
		Index:   proto.Uint64(this.index),
		Entries: this.entries,
	}

	return proto.Marshal(copyState)
}
func (this CopyStateRPC) Decode(b []byte) error {
	pb := new(p.CopyState)
	err := proto.Unmarshal(b, pb)

	if err != nil {
		this.term = pb.GetTerm()
		this.index = pb.GetIndex()
		this.voting = pb.GetVoting()
		this.entries = pb.GetEntries()
	}

	return err
}
func (this CopyStateRPC) GetVoting() bool {
	return this.voting
}
func (this CopyStateRPC) GetEntries() []*p.Entry {
	return this.entries
}
