package AppendEntryResponse

import (
    p "raft/pkg/protobuf"
	"raft/internal/messages"
)

type AppendEntryResponse struct
{
    
}

func new_AppendEntryResponse() messages.Rpc{
    return &AppendEntryResponse{}
}

func (this AppendEntryResponse) GetTerm() uint64{
    return 0
}
func (this AppendEntryResponse) GetVoting() bool{
    return false
}
func (this AppendEntryResponse) GetEntries() []*p.Entry{
    return nil
}
func (this AppendEntryResponse) GetLeaderId() string{
    return ""
}
func (this AppendEntryResponse) GetCandidateId() uint64{
    return 0
}
func (this AppendEntryResponse) GetLastLogTerm() uint64{
    return 0
}
func (this AppendEntryResponse) GetLastLogIndex() uint64{
    return 0
}
func (this AppendEntryResponse) GetPrevLogTerm() uint64{
    return 0
}
func (this AppendEntryResponse) GetPrevLogIndex() uint64{
    return 0
}
func (this AppendEntryResponse) GetLeaderCommit() uint64{
    return 0
}
func (this AppendEntryResponse) HasSucceded() bool{
    return false
}
func (this AppendEntryResponse) VoteGranted() bool{
    return false
}
func (this AppendEntryResponse) GetIndex() uint64{
    return 0
}
func (this AppendEntryResponse) Encode() ([]byte, error){
    return nil,nil
}
func (this AppendEntryResponse) Decode() error{
    return nil
}

