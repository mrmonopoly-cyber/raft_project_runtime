package AppendEntryResponse

import (
    p "raft/pkg/protobuf"
	"raft/internal/messages"
	"google.golang.org/protobuf/proto"
)

type AppendEntryResponse struct
{
  success bool
  term uint64
  logIndexError uint64
}

func new_AppendEntryResponse() messages.Rpc{
    return &AppendEntryResponse{}
}

func (this AppendEntryResponse) GetTerm() uint64{
    return this.term
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
    return this.success
}
func (this AppendEntryResponse) VoteGranted() bool{
    return false
}
func (this AppendEntryResponse) GetIndex() uint64{
    return this.logIndexError
}
func (this AppendEntryResponse) Encode() ([]byte, error){

    response := &p.AppendEntryResponse{
        Success: proto.Bool(this.success),
        Term: proto.Uint64(this.term),
        LogIndexError: proto.Uint64(this.logIndexError),
    }

    return proto.Marshal(response)
}
func (this AppendEntryResponse) Decode(b []byte) error{

    pb := new(p.AppendEntryResponse)
    err := proto.Unmarshal(b, pb)

    if err != nil {
        this.term = pb.GetTerm()
        this.success = pb.GetSuccess()
        this.logIndexError = pb.GetLogIndexError()
    }

    return err
}

