package RequestVoteRPC

import (
	"log"
	"raft/internal/rpcs"
	"raft/internal/raftstate"
	"strconv"
    "raft/pkg/rpcEncoding/out/protobuf"
    "raft/internal/rpcs/RequestVoteResponse"

	"google.golang.org/protobuf/proto"
)

type RequestVoteRPC struct {
    pMex protobuf.RequestVote
}

func NewRequestVoteRPC(term uint64, candidateId string,
lastLogIndex int64, lastLogTerm uint64) rpcs.Rpc {
    return &RequestVoteRPC{
        pMex : protobuf.RequestVote{
            Term: term,
            CandidateId: candidateId,
            LastLogIndex: lastLogIndex,
            LastLogTerm: lastLogTerm,
        },
    }
}

// GetId rpcs.Rpc.
func (this *RequestVoteRPC) GetId() string {
  return this.pMex.CandidateId
}

// ToString rpcs.Rpc.
func (this *RequestVoteRPC) ToString() string {
	var mex string = "{term : " + strconv.Itoa(int(this.GetTerm())) + ", leaderId: " + this.GetCandidateId()+ ",lastLogIndex: " + strconv.Itoa(int(this.GetLastLogIndex())) + ", lastLogTerm: " + strconv.Itoa(int(this.GetLastLogTerm())) + "}"

    log.Println("rpc RequestVote :", mex)

    return mex
}

// GetTerm rpcs.Rpc.
func (this *RequestVoteRPC) GetTerm() uint64 {
	return this.pMex.Term
}

// Encode rpcs.Rpc.
func (this *RequestVoteRPC) Encode() ([]byte, error) {
    var mess []byte
    var err error

    mess, err = proto.Marshal(&(*this).pMex)
    if err != nil {
        log.Panicln("error in Encoding Request Vote: ", err)
    }

	return mess, err
}

// Decode rpcs.Rpc.
func (this *RequestVoteRPC) Decode(b []byte) error {
	err := proto.Unmarshal(b,&this.pMex)
    if err != nil {
        log.Panicln("error in Encoding Request Vote: ", err)
    }
	return err
}

// GetCandidateId rpcs.Rpc.
func (this *RequestVoteRPC) GetCandidateId() string {
	return this.pMex.CandidateId
}

// GetLastLogIndex  rpcs.Rpc.
func (this *RequestVoteRPC) GetLastLogIndex() int64 {
	return this.pMex.LastLogIndex
}

// GetLastLogTerm rpcs.Rpc.
func (this *RequestVoteRPC) GetLastLogTerm() uint64 {
	return this.pMex.LastLogTerm
}

// Manage implements rpcs.Rpc.
func (this *RequestVoteRPC) Execute(state *raftstate.State) *rpcs.Rpc{
    var myVote string = (*state).GetVoteFor()
    var sender = this.GetCandidateId()

    if ! (*state).CanVote(){
        log.Printf("request vote: this node cannot vote right now")
        return nil
    }
    if this.GetTerm() < (*state).GetTerm() {
        log.Printf("request vote: not valid term vote false: my term %v, other therm %v", 
                    (*state).GetTerm(), this.GetTerm())
        return this.respondeVote(state, &sender, false)
    }

    if ! (*state).MoreRecentLog(this.GetLastLogIndex(), this.GetLastLogTerm()) {
        log.Printf("request vote: log not recent enough")
        return this.respondeVote(state, &sender,false)
    }else if myVote == "" || myVote == this.GetCandidateId(){
        log.Printf("request vote: vote accepted, voting for: %v", this.GetCandidateId())
        (*state).VoteFor(sender)
        return this.respondeVote(state,&sender,true)
    }

    return this.respondeVote(state,&sender,true)
}

func (this *RequestVoteRPC) respondeVote(state *raftstate.State, sender *string, vote bool) *rpcs.Rpc{
    var resp = RequestVoteResponse.NewRequestVoteResponseRPC(*sender,vote, (*state).GetTerm())
    return &resp
}
