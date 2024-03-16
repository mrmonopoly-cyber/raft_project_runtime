package RequestVoteResponse

import (
	"google.golang.org/protobuf/proto"
	"raft/internal/messages"
	p "raft/pkg/protobuf"
)

type RequestVoteResponse struct {
	voteGranted bool
	term        uint64
}

// ToMessage implements messages.Rpc.
func (this *RequestVoteResponse) ToMessage() messages.Message {
	panic("unimplemented")
}

// ToString implements messages.Rpc.
func (this *RequestVoteResponse) ToString() string {
	panic("unimplemented")
}

func new_RequestVoteResponse(voteGranted bool, term uint64) messages.Rpc {
	return &RequestVoteResponse{
		voteGranted: voteGranted,
		term:        term,
	}
}

func (this RequestVoteResponse) GetTerm() uint64 {
	return this.term
}

func (this RequestVoteResponse) Encode() ([]byte, error) {

	response := &p.RequestVoteResponse{
		VoteGranted: proto.Bool(this.voteGranted),
		Term:        proto.Uint64(this.term),
	}

	return proto.Marshal(response)
}
func (this RequestVoteResponse) Decode(b []byte) error {
	pb := new(p.RequestVoteResponse)
	err := proto.Unmarshal(b, pb)

	if err != nil {
		this.term = pb.GetTerm()
		this.voteGranted = pb.GetVoteGranted()
	}

	return err
}

func (this RequestVoteResponse) VoteGranted() bool {
	return this.voteGranted
}

func (s RequestVoteResponse) other_node_vote_candidature() {

}
