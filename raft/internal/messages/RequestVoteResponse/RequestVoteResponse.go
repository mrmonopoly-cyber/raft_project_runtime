package RequestVoteResponse

import (
	"raft/internal/messages"
	"raft/internal/raftstate"
	p "raft/pkg/protobuf"
	"strconv"

	"google.golang.org/protobuf/proto"
)

type RequestVoteResponse struct {
	id          string
	voteGranted bool
	term        uint64
}

func NewRequestVoteResponse(id string, voteGranted bool, term uint64) messages.Rpc {
	return &RequestVoteResponse{
		id:          id,
		voteGranted: voteGranted,
		term:        term,
	}
}

// Manage implements messages.Rpc.
func (this *RequestVoteResponse) Execute(state *raftstate.State) *messages.Rpc {
	panic("unimplemented")
}

// ToString implements messages.Rpc.
func (this *RequestVoteResponse) ToString() string {
	return "{term : " + strconv.Itoa(int(this.term)) + ", \nvoteGranted: " + strconv.FormatBool(this.voteGranted) + "}"
}


func (this RequestVoteResponse) GetId() string {
	return this.id
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
