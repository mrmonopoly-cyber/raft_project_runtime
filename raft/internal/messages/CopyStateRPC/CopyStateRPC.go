package CopyStateRPC

import (
	"raft/internal/messages"
	"raft/internal/raftstate"
	p "raft/pkg/protobuf"
	"strconv"

	"google.golang.org/protobuf/proto"
)

type CopyStateRPC struct {
	id      string
	term    uint64
	index   uint64
	voting  bool
	entries []*p.Entry
}

// Manage implements messages.Rpc.
func (this *CopyStateRPC) Execute(state *raftstate.State) *messages.Rpc{
	panic("unimplemented")
}

// ToString implements messages.Rpc.
func (this *CopyStateRPC) ToString() string {
	var entries string
	for _, el := range this.entries {
		entries += el.String()
	}
	return "{term : " + strconv.Itoa(int(this.term)) + ", \nindex: " + strconv.Itoa(int(this.index)) + ",\nvoting: " + strconv.FormatBool(this.voting) + ", \nentries: " + entries + "}"
}

func newCopyStateRPC(id string, term uint64, index uint64, voting bool,
	entries []*p.Entry) messages.Rpc {
	return &CopyStateRPC{
		id:      id,
		term:    term,
		index:   index,
		voting:  voting,
		entries: entries,
	}
}

func (this CopyStateRPC) GetId() string {
	return this.id
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
