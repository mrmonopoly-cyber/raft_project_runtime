package RequestVoteResponse

import (
	"raft/internal/messages"
	"raft/internal/raftstate"
	p "raft/pkg/protobuf"
	"strconv"
	"sync"

	"google.golang.org/protobuf/proto"
)

type RequestVoteResponse struct {
	voteGranted bool
	term        uint64
}

// Manage implements messages.Rpc.
func (this *RequestVoteResponse) Execute(n *sync.Map, state raftstate.State) {
	panic("unimplemented")
}

// ToString implements messages.Rpc.
func (this *RequestVoteResponse) ToString() string {
	return "{term : " + strconv.Itoa(int(this.term)) + ", \nvoteGranted: " + strconv.FormatBool(this.voteGranted) + "}"
}

func newRequestVoteResponse(voteGranted bool, term uint64) messages.Rpc {
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

func (s RequestVoteResponse) otherNodeVoteCandidature() {

}
