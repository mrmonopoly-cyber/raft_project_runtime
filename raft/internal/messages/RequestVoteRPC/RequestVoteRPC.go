package RequestVoteRPC

import (
	"log"
	"raft/internal/messages"
	"raft/internal/messages/RequestVoteResponse"
	"raft/internal/raftstate"
	p "raft/pkg/protobuf"
	"strconv"

	"google.golang.org/protobuf/proto"
)

type RequestVoteRPC struct {
	term         uint64
	candidateId  string
	lastLogIndex uint64
	lastLogTerm  uint64
}

func NewRequestVoteRPC(term uint64, candidateId string,
	lastLogIndex uint64, lastLogTerm uint64) messages.Rpc {
	return &RequestVoteRPC{
		term,
		candidateId,
		lastLogIndex,
		lastLogTerm,
	}
}

// GetId messages.Rpc.
func (this RequestVoteRPC) GetId() string {
  return this.candidateId
}

// ToString messages.Rpc.
func (this *RequestVoteRPC) ToString() string {
	var mex string = "{term : " + strconv.Itoa(int(this.term)) + ", leaderId: " + this.candidateId + ",lastLogIndex: " + strconv.Itoa(int(this.lastLogIndex)) + ", lastLogTerm: " + strconv.Itoa(int(this.lastLogTerm)) + "}"

    log.Println("rpc RequestVote :", mex)

    return mex
}

// GetTerm messages.Rpc.
func (this RequestVoteRPC) GetTerm() uint64 {
	return this.term
}

// Encode messages.Rpc.
func (this RequestVoteRPC) Encode() ([]byte, error) {
	reqVote := &p.RequestVote{
		Term:         proto.Uint64(this.term),
		CandidateId:  proto.String(this.candidateId),
		LastLogIndex: proto.Uint64(this.lastLogIndex),
		LastLogTerm:  proto.Uint64(this.lastLogTerm),
	}

	mess, err := proto.Marshal(reqVote)
	return mess, err
}

// Decode messages.Rpc.
func (this RequestVoteRPC) Decode(b []byte) error {
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

// GetCandidateId messages.Rpc.
func (this RequestVoteRPC) GetCandidateId() string {
	return this.candidateId
}

// GetLastLogIndex  messages.Rpc.
func (this RequestVoteRPC) GetLastLogIndex() uint64 {
	return this.lastLogIndex
}

// GetLastLogTerm messages.Rpc.
func (this RequestVoteRPC) GetLastLogTerm() uint64 {
	return this.lastLogTerm
}

// Manage implements messages.Rpc.
func (this *RequestVoteRPC) Execute(state *raftstate.State) *messages.Rpc{
    var myVote string = (*state).GetVoteFor()
    var sender = this.GetCandidateId()

    if ! (*state).CanVote(){
        return nil
    }
    if this.term < (*state).GetTerm() {
        return this.respondeVote(state, &sender, false)
    }

    if ! (*state).MoreRecentLog(this.GetLastLogIndex(), this.GetLastLogTerm()) {
        return this.respondeVote(state, &sender,false)
    }else if myVote == "" || myVote == this.GetCandidateId(){
        this.respondeVote(state,&sender,true)
        (*state).VoteFor(sender)
    }

    return this.respondeVote(state,&sender,false)
}

func (this *RequestVoteRPC) respondeVote(state *raftstate.State, sender *string, vote bool) *messages.Rpc{
        var resp = RequestVoteResponse.NewRequestVoteResponse(*sender,vote,(*state).GetTerm())
        return &resp
}
