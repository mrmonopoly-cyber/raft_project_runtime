package CopyStateRPC

import (
    p "raft/pkg/protobuf"
	"raft/internal/messages"
)

type CopyStateRPC struct
{
    
}

func new_CopyStateRPC() messages.Rpc{
    return &CopyStateRPC{}
}

func (this CopyStateRPC) GetTerm() uint64{
    return 0
}
func (this CopyStateRPC) GetVoting() bool{
    return false
}
func (this CopyStateRPC) GetEntries() []*p.Entry{
    return nil
}
func (this CopyStateRPC) GetLeaderId() string{
    return ""
}
func (this CopyStateRPC) GetCandidateId() uint64{
    return 0
}
func (this CopyStateRPC) GetLastLogTerm() uint64{
    return 0
}
func (this CopyStateRPC) GetLastLogIndex() uint64{
    return 0
}
func (this CopyStateRPC) GetPrevLogTerm() uint64{
    return 0
}
func (this CopyStateRPC) GetPrevLogIndex() uint64{
    return 0
}
func (this CopyStateRPC) GetLeaderCommit() uint64{
    return 0
}
func (this CopyStateRPC) HasSucceded() bool{
    return false
}
func (this CopyStateRPC) VoteGranted() bool{
    return false
}
func (this CopyStateRPC) GetIndex() uint64{
    return 0
}
func (this CopyStateRPC) Encode() ([]byte, error){
    return nil,nil
}
func (this CopyStateRPC) Decode() error{
    return nil
}

