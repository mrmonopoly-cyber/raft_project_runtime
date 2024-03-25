package RequestVoteRPC

import (
	"log"
	"raft/internal/rpcs"
	"raft/internal/raftstate"
	"strconv"
    "raft/pkg/rpcEncoding/out/protobuf"

	"google.golang.org/protobuf/proto"
)

type RequestVoteRPC struct {
    pMex protobuf.RequestVote
}

func NewRequestVoteRPC(term uint64, candidateId string,
lastLogIndex uint64, lastLogTerm uint64) rpcs.Rpc {
    log.Println("creating new request vote with candidateId: ", candidateId)
    return &RequestVoteRPC{
        protobuf.RequestVote{
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
	var mex string = "{term : " + strconv.Itoa(int(this.pMex.Term)) + ", leaderId: " + this.pMex.CandidateId + ",lastLogIndex: " + strconv.Itoa(int(this.pMex.LastLogIndex)) + ", lastLogTerm: " + strconv.Itoa(int(this.pMex.LastLogTerm)) + "}"

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
	return mess, err
}

// Decode rpcs.Rpc.
func (this *RequestVoteRPC) Decode(b []byte) error {
    pb := new(protobuf.RequestVote)
	err := proto.Unmarshal(b, pb)

	if err != nil {
        this.pMex.Term = pb.GetTerm()
		this.pMex.CandidateId = pb.GetCandidateId()
		this.pMex.LastLogTerm = pb.GetLastLogTerm()
		this.pMex.LastLogIndex = pb.GetLastLogIndex()
	}

	return err
}

// GetCandidateId rpcs.Rpc.
func (this *RequestVoteRPC) GetCandidateId() string {
	return this.pMex.CandidateId
}

// GetLastLogIndex  rpcs.Rpc.
func (this *RequestVoteRPC) GetLastLogIndex() uint64 {
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
    if this.pMex.Term < (*state).GetTerm() {
        log.Printf("request vote: not valid term vote false: my term %v, other therm %v", 
                    (*state).GetTerm(), this.GetTerm())
        return this.respondeVote(state, &sender, false)
    }

    if ! (*state).MoreRecentLog(this.GetLastLogIndex(), this.GetLastLogTerm()) {
        log.Printf("request vote: log not recent enough")
        return this.respondeVote(state, &sender,false)
    }else if myVote == "" || myVote == this.GetCandidateId(){
        log.Printf("request vote: vote accepted, voting for: %v", this.GetCandidateId())
        this.respondeVote(state,&sender,true)
        (*state).VoteFor(sender)
    }

    return this.respondeVote(state,&sender,true)
}

func (this *RequestVoteRPC) respondeVote(state *raftstate.State, sender *string, vote bool) *rpcs.Rpc{
    panic("non implemented")
}
