package RequestVoteResponse

import (
    p "raft/pkg/protobuf"
	"raft/internal/messages"
)

type RequestVoteResponse struct
{
    
}

func new_RequestVoteResponse() messages.Rpc{
    return &RequestVoteResponse{}
}

func (this RequestVoteResponse) GetTerm() uint64{
    return 0
}
func (this RequestVoteResponse) GetVoting() bool{
    return false
}
func (this RequestVoteResponse) GetEntries() []*p.Entry{
    return nil
}
func (this RequestVoteResponse) GetLeaderId() string{
    return ""
}
func (this RequestVoteResponse) GetCandidateId() uint64{
    return 0
}
func (this RequestVoteResponse) GetLastLogTerm() uint64{
    return 0
}
func (this RequestVoteResponse) GetLastLogIndex() uint64{
    return 0
}
func (this RequestVoteResponse) GetPrevLogTerm() uint64{
    return 0
}
func (this RequestVoteResponse) GetPrevLogIndex() uint64{
    return 0
}
func (this RequestVoteResponse) GetLeaderCommit() uint64{
    return 0
}
func (this RequestVoteResponse) HasSucceded() bool{
    return false
}
func (this RequestVoteResponse) VoteGranted() bool{
    return false
}
func (this RequestVoteResponse) GetIndex() uint64{
    return 0
}
func (this RequestVoteResponse) Encode() ([]byte, error){
    return nil,nil
}
func (this RequestVoteResponse) Decode() error{
    return nil
}

