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
  GetCandidateId() uint64
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


type RequestVoteRPC struct {
  term uint64
  candidateId string
  lastLogIndex uint64
  lastLogTerm uint64
}


type CopyStateRPC struct {
  term uint64
  index uint64
  voting bool
  entries []*p.Entry
}

type RequestVoteResponse struct {
  voteGranted bool
  term uint64
}

func (m *RequestVoteRPC) GetTerm() uint64 {
  return m.term
}

func (m *RequestVoteResponse) GetTerm() uint64 {
  return m.term
}

func (m *CopyStateRPC) GetTerm() uint64 {
  return m.term
}

func (m *CopyStateRPC) GetVoting() bool {
  return m.voting
}

func (m *CopyStateRPC) GetEntries() []*p.Entry {
  return m.entries
}

func (m *RequestVoteRPC) GetCandidateId() string {
  return m.candidateId
}

func (m *RequestVoteRPC) GetLastLogIndex() uint64 {
  return m.lastLogIndex
}

func (m *RequestVoteRPC) GetLastLogTerm() uint64 {
  return m.lastLogTerm
}

func (m *RequestVoteResponse) VoteGranted() bool {
  return m.voteGranted
}

func (m *CopyStateRPC) GetIndex() uint64 {
  return m.index
}

func (m *RequestVoteRPC) Encode() ([]byte, error) {
  reqVote := &p.RequestVote{
    Term: proto.Uint64(m.term),
    CandidateId: proto.String(m.candidateId),
    LastLogIndex: proto.Uint64(m.lastLogIndex),
    LastLogTerm: proto.Uint64(m.lastLogTerm),
  }

  mess, err := proto.Marshal(reqVote)
  return mess, err
}

func (m *CopyStateRPC) Encode() ([]byte, error) {
  copyState := &p.CopyState{
    Term: proto.Uint64(m.term),
    Voting: proto.Bool(m.voting),
    Index: proto.Uint64(m.index),
    Entries: m.entries,
  }

  return proto.Marshal(copyState)
}

func (m *RequestVoteResponse) Encode() ([]byte, error) {
  response := &p.RequestVoteResponse{
    VoteGranted: proto.Bool(m.voteGranted),
    Term: proto.Uint64(m.term),
  }

  return proto.Marshal(response)
}

func (m *RequestVoteRPC) Decode(b []byte) error {
  pb := new(p.RequestVote)
  err := proto.Unmarshal(b, pb)
  
  if err != nil {
    m.term = pb.GetTerm()
    m.candidateId = pb.GetCandidateId()
    m.lastLogTerm = pb.GetLastLogTerm()
    m.lastLogIndex = pb.GetLastLogIndex()
  }

  return err
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

func (m *CopyStateRPC) Decode(b []byte) error {
  pb := new(p.CopyState)
  err := proto.Unmarshal(b, pb)

  if err != nil {
    m.term = pb.GetTerm()
    m.index = pb.GetIndex()
    m.voting = pb.GetVoting()
    m.entries = pb.GetEntries()
  }

  return err
}

func NewRequestVote(term uint64, candidateId string, lastLogIndex uint64, lastLogTerm uint64) RequestVoteRPC {
  return RequestVoteRPC{
    term,
    candidateId,
    lastLogIndex,
    lastLogTerm,
  }
}

func NewRequestVoteResponse(voteGranted bool, term uint64) RequestVoteResponse {
    return RequestVoteResponse{
        voteGranted,
        term,
    }
}

