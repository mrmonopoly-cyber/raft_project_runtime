package CopyStateRPC

import (
	"google.golang.org/protobuf/proto"
	"raft/internal/messages"
	p "raft/pkg/protobuf"
  "strconv"
)

type CopyStateRPC struct {
	term    uint64
	index   uint64
	voting  bool
	entries []*p.Entry
}

// ToMessage implements messages.Rpc.
func (this *CopyStateRPC) ToMessage() messages.Message {
  enc, _ := this.Encode()
  return messages.Message{
    Mex_type: messages.COPY_STATE,
    Payload: enc,
  }  
}

// ToString implements messages.Rpc.
func (this *CopyStateRPC) ToString() string {	
  var entries string 
  for _, el := range this.entries {
    entries += el.String()
  }
  return "{term : " + strconv.Itoa(int(this.term)) + ", \nindex: " + strconv.Itoa(int(this.index)) + ",\nvoting: " + strconv.FormatBool(this.voting) + ", \nentries: " + entries + "}"
}

func newCopyStateRPC(term uint64, index uint64, voting bool,
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
