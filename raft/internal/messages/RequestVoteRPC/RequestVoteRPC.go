package RequestVoteRPC

import (
	"raft/internal/messages"
	"raft/internal/raftstate"
	p "raft/pkg/protobuf"
	"strconv"

	"google.golang.org/protobuf/proto"
)

type RequestVoteRPC struct {
	term         uint64
	candidateId  string
	lastLogIndex uint64
	lastLogTerm  uint64
}

func NewRequestVoteRPC(term uint64, candidateId string,
	lastLogIndex uint64, lastLogTerm uint64) messages.Rpc {
	return &RequestVoteRPC{
		term,
		candidateId,
		lastLogIndex,
		lastLogTerm,
	}
}

// GetId messages.Rpc.
func (this RequestVoteRPC) GetId() string {
  return this.candidateId
}

// ToString messages.Rpc.
func (this *RequestVoteRPC) ToString() string {
	return "{term : " + strconv.Itoa(int(this.term)) + ", \nleaderId: " + this.candidateId + ",\nlastLogIndex: " + strconv.Itoa(int(this.lastLogIndex)) + ", \nlastLogTerm: " + strconv.Itoa(int(this.lastLogTerm)) + "}"

}

// GetTerm messages.Rpc.
func (this RequestVoteRPC) GetTerm() uint64 {
	return this.term
}

// Encode messages.Rpc.
func (this RequestVoteRPC) Encode() ([]byte, error) {
	reqVote := &p.RequestVote{
		Term:         proto.Uint64(this.term),
		CandidateId:  proto.String(this.candidateId),
		LastLogIndex: proto.Uint64(this.lastLogIndex),
		LastLogTerm:  proto.Uint64(this.lastLogTerm),
	}

	mess, err := proto.Marshal(reqVote)
	return mess, err
}

// Decode messages.Rpc.
func (this RequestVoteRPC) Decode(b []byte) error {
	pb := new(p.RequestVote)
	err := proto.Unmarshal(b, pb)

	if err != nil {
		this.term = pb.GetTerm()
		this.candidateId = pb.GetCandidateId()
		this.lastLogTerm = pb.GetLastLogTerm()
		this.lastLogIndex = pb.GetLastLogIndex()
	}

	return err
}

// GetCandidateId messages.Rpc.
func (this RequestVoteRPC) GetCandidateId() string {
	return this.candidateId
}

// GetLastLogIndex  messages.Rpc.
func (this RequestVoteRPC) GetLastLogIndex() uint64 {
	return this.lastLogIndex
}

// GetLastLogTerm messages.Rpc.
func (this RequestVoteRPC) GetLastLogTerm() uint64 {
	return this.lastLogTerm
}

// Manage implements messages.Rpc.
func (this *RequestVoteRPC) Execute(state *raftstate.State, sender *string) *messages.Rpc{
    panic("unimplemented")
}
