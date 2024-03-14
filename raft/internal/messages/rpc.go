package messages

import (
	p "raft/pkg/protobuf"

	//"golang.org/x/term"
	"google.golang.org/protobuf/proto"
)

type Rpc interface {
  GetTerm() uint64
  GetVoting() bool 
  GetEntries() []*p.Entry
  GetLeaderId() string
  GetLastLogTerm() uint64
  GetLastLogIndex() uint64
  GetPrevLogTerm() uint64
  GetPrevLogIndex() uint64
  GetLeaderCommit() uint64
  HasSucceded() bool
  VoteGranted() bool
  GetIndex() uint64
  Encode() ([]byte, error)
  Decode([]byte) error
  // NewRequestVote(term uint64, candidateId string, lastLogIndex uint64, lastLogTerm uint64) RequestVoteRPC
}

type RequestVoteResponse struct {
  voteGranted bool
  term uint64
}

func (m *RequestVoteResponse) GetTerm() uint64 {
  return m.term
}

func (m *RequestVoteResponse) VoteGranted() bool {
  return m.voteGranted
}

func (m *RequestVoteResponse) Encode() ([]byte, error) {
  response := &p.RequestVoteResponse{
    VoteGranted: proto.Bool(m.voteGranted),
    Term: proto.Uint64(m.term),
  }

  return proto.Marshal(response)
}

func (m *RequestVoteResponse) Decode(b []byte) error {
  pb := new(p.RequestVoteResponse)
  err := proto.Unmarshal(b, pb)

  if err != nil {
    m.term = pb.GetTerm()
    m.voteGranted = pb.GetVoteGranted()
  }

  return err
}

func NewRequestVoteResponse(voteGranted bool, term uint64) RequestVoteResponse {
    return RequestVoteResponse{
        voteGranted,
        term,
    }
}

