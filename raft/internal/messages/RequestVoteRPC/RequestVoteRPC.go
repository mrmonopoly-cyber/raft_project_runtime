package RequestVoteRPC

import (
	"google.golang.org/protobuf/proto"
	"raft/internal/messages"
	p "raft/pkg/protobuf"
)

type RequestVoteRPC struct {
	term         uint64
	candidateId  string
	lastLogIndex uint64
	lastLogTerm  uint64
}

// ToMessage implements messages.Rpc.
func (this *RequestVoteRPC) ToMessage() messages.Message {
	panic("unimplemented")
}

func New_RequestVoteRPC(term uint64, candidateId string,
	lastLogIndex uint64, lastLogTerm uint64) messages.Rpc {
	return &RequestVoteRPC{
		term,
		candidateId,
		lastLogIndex,
		lastLogTerm,
	}
}

// ToString implements messages.Rpc.
func (this *RequestVoteRPC) ToString() string {
	return ""
}

func (this RequestVoteRPC) GetTerm() uint64 {
	return this.term
}

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
func (this RequestVoteRPC) GetCandidateId() string {
	return this.candidateId
}
func (this RequestVoteRPC) GetLastLogIndex() uint64 {
	return this.lastLogIndex
}
func (this RequestVoteRPC) GetLastLogTerm() uint64 {
	return this.lastLogTerm
}
