package extra

import (
    p "raft/pkg/protobuf"
	"raft/internal/messages"
	"google.golang.org/protobuf/proto"
)

type NEW_RPC struct
{
    
}

func new_NEW_RPC() messages.Rpc{
    return &NEW_RPC{}
}

func (this NEW_RPC) GetTerm() uint64{
    return 0
}
func (this NEW_RPC) GetVoting() bool{
    return false
}
func (this NEW_RPC) GetEntries() []*p.Entry{
    return nil
}
func (this NEW_RPC) GetLeaderId() string{
    return ""
}
func (this NEW_RPC) GetCandidateId() uint64{
    return 0
}
func (this NEW_RPC) GetLastLogTerm() uint64{
    return 0
}
func (this NEW_RPC) GetLastLogIndex() uint64{
    return 0
}
func (this NEW_RPC) GetPrevLogTerm() uint64{
    return 0
}
func (this NEW_RPC) GetPrevLogIndex() uint64{
    return 0
}
func (this NEW_RPC) GetLeaderCommit() uint64{
    return 0
}
func (this NEW_RPC) HasSucceded() bool{
    return false
}
func (this NEW_RPC) VoteGranted() bool{
    return false
}
func (this NEW_RPC) GetIndex() uint64{
    return 0
}
func (this NEW_RPC) Encode() ([]byte, error){
    return nil,nil
}
func (this NEW_RPC) Decode(b []byte) error{
    return nil
}

