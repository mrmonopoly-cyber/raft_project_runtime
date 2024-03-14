package RequestVoteRPC

import (
    p "raft/pkg/protobuf"
	"raft/internal/messages"
	"google.golang.org/protobuf/proto"
)

type RequestVoteRPC struct
{
    term uint64
    candidateId string
    lastLogIndex uint64
    lastLogTerm uint64
}

func new_RequestVoteRPC(term uint64, candidateId string, 
                        lastLogIndex uint64, lastLogTerm uint64) messages.Rpc{
    return &RequestVoteRPC{
        term,
        candidateId,
        lastLogIndex,
        lastLogTerm,
    }
}

func (this RequestVoteRPC) GetTerm() uint64{
    return this.term
}
func (this RequestVoteRPC) Encode() ([]byte, error){
    reqVote := &p.RequestVote{
        Term: proto.Uint64(this.term),
        CandidateId: proto.String(this.candidateId),
        LastLogIndex: proto.Uint64(this.lastLogIndex),
        LastLogTerm: proto.Uint64(this.lastLogTerm),
    }

    mess, err := proto.Marshal(reqVote)
    return mess, err
}
func (this RequestVoteRPC) Decode(b []byte) error{
    pb := new(p.RequestVote)
    err := proto.Unmarshal(b, pb)

    if err != nil {
        this.term = pb.GetTerm()
        this.candidateId = pb.GetCandidateId()
        this.lastLogTerm = pb.GetLastLogTerm()
        this.lastLogIndex = pb.GetLastLogIndex()
    }

    return err
}
func (this RequestVoteRPC) GetCandidateId() string{
    return this.candidateId
}
func (this RequestVoteRPC) GetLastLogIndex() uint64{
    return this.lastLogIndex
}
func (this RequestVoteRPC) GetPrevLogTerm() uint64{
    return this.lastLogTerm
}
