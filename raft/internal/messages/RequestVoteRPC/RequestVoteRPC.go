package RequestVoteRPC

import (
    p "raft/pkg/protobuf"
	"raft/internal/messages"
)

type RequestVoteRPC struct
{
    
}

func new_RequestVoteRPC() messages.Rpc{
    return &RequestVoteRPC{}
}

func (this RequestVoteRPC) GetTerm() uint64{
    return 0
}
func (this RequestVoteRPC) GetVoting() bool{
    return false
}
func (this RequestVoteRPC) GetEntries() []*p.Entry{
    return nil
}
func (this RequestVoteRPC) GetLeaderId() string{
    return ""
}
func (this RequestVoteRPC) GetCandidateId() uint64{
    return 0
}
func (this RequestVoteRPC) GetLastLogTerm() uint64{
    return 0
}
func (this RequestVoteRPC) GetLastLogIndex() uint64{
    return 0
}
func (this RequestVoteRPC) GetPrevLogTerm() uint64{
    return 0
}
func (this RequestVoteRPC) GetPrevLogIndex() uint64{
    return 0
}
func (this RequestVoteRPC) GetLeaderCommit() uint64{
    return 0
}
func (this RequestVoteRPC) HasSucceded() bool{
    return false
}
func (this RequestVoteRPC) VoteGranted() bool{
    return false
}
func (this RequestVoteRPC) GetIndex() uint64{
    return 0
}
func (this RequestVoteRPC) Encode() ([]byte, error){
    return nil,nil
}
func (this RequestVoteRPC) Decode() error{
    return nil
}

