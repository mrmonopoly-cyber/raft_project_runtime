package RequestVoteRPC

import (
	"log"
	"raft/internal/messages"
	"raft/internal/raftstate"
	"strconv"
    "raft/pkg/rpcEncoding/out/protobuf"

	"google.golang.org/protobuf/proto"
)

type RequestVoteRPC struct {
    pMex protobuf.RequestVote
}

func NewRequestVoteRPC(term uint64, candidateId string,
	lastLogIndex uint64, lastLogTerm uint64) messages.Rpc {
	return &RequestVoteRPC{
        protobuf.RequestVote{
            Ty: protobuf.MexType_REQUEST_VOTE,
            Term: term,
            CandidateId: candidateId,
            LastLogIndex: lastLogIndex,
            LastLogTerm: lastLogTerm,
        },
	}
}

// GetId messages.Rpc.
func (this *RequestVoteRPC) GetId() string {
  return this.pMex.CandidateId
}

// ToString messages.Rpc.
func (this *RequestVoteRPC) ToString() string {
	var mex string = "{term : " + strconv.Itoa(int(this.pMex.Term)) + ", leaderId: " + this.pMex.CandidateId + ",lastLogIndex: " + strconv.Itoa(int(this.pMex.LastLogIndex)) + ", lastLogTerm: " + strconv.Itoa(int(this.pMex.LastLogTerm)) + "}"

    log.Println("rpc RequestVote :", mex)

    return mex
}

// GetTerm messages.Rpc.
func (this *RequestVoteRPC) GetTerm() uint64 {
	return this.pMex.Term
}

// Encode messages.Rpc.
func (this *RequestVoteRPC) Encode() ([]byte, error) {
    var mess []byte
    var err error
	mess, err = proto.Marshal(&(*this).pMex)
	return mess, err
}

// Decode messages.Rpc.
func (this *RequestVoteRPC) Decode(b []byte) error {
    panic("unimplemented")
}

// GetCandidateId messages.Rpc.
func (this *RequestVoteRPC) GetCandidateId() string {
	return this.pMex.CandidateId
}

// GetLastLogIndex  messages.Rpc.
func (this *RequestVoteRPC) GetLastLogIndex() uint64 {
	return this.pMex.LastLogIndex
}

// GetLastLogTerm messages.Rpc.
func (this *RequestVoteRPC) GetLastLogTerm() uint64 {
	return this.pMex.LastLogTerm
}

// Manage implements messages.Rpc.
func (this *RequestVoteRPC) Execute(state *raftstate.State) *messages.Rpc{
    var myVote string = (*state).GetVoteFor()
    var sender = this.GetCandidateId()

    if ! (*state).CanVote(){
        return nil
    }
    if this.pMex.Term < (*state).GetTerm() {
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
    panic("non implemented")
}
