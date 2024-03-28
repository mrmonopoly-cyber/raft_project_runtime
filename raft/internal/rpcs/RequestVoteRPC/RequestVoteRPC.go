package RequestVoteRPC

import (
	"log"
	"raft/internal/raftstate"
	"raft/internal/rpcs"
	"raft/pkg/rpcEncoding/out/protobuf"
	"strconv"

	"google.golang.org/protobuf/proto"
)

type RequestVoteRPC struct {
	pMex protobuf.RequestVote
}

func NewRequestVoteRPC(term uint64, candidateId string,
	lastLogIndex uint64, lastLogTerm uint64) rpcs.Rpc {
	return &RequestVoteRPC{
		pMex: protobuf.RequestVote{
			Term:         term,
			CandidateId:  candidateId,
			LastLogIndex: lastLogIndex,
			LastLogTerm:  lastLogTerm,
		},
	}
}

// ToString rpcs.Rpc.
func (this *RequestVoteRPC) ToString() string {
	var mex string = "{term : " + strconv.Itoa(int(this.pMex.GetTerm())) + ", leaderId: " + this.pMex.GetCandidateId() + ",lastLogIndex: " + strconv.Itoa(int(this.pMex.GetLastLogIndex())) + ", lastLogTerm: " + strconv.Itoa(int(this.pMex.GetLastLogTerm())) + "}"

	log.Println("rpc RequestVote :", mex)

	return mex
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

// Manage implements rpcs.Rpc.
func (this *RequestVoteRPC) Execute(state *raftstate.State) *rpcs.Rpc {
	var myVote string = (*state).GetVoteFor()
	var sender = this.pMex.GetCandidateId()

	if !(*state).CanVote() {
		log.Printf("request vote: this node cannot vote right now")
		return nil
	}
	if this.pMex.GetTerm() < (*state).GetTerm() {
		log.Printf("request vote: not valid term vote false: my term %v, other therm %v",
			(*state).GetTerm(), this.pMex.GetTerm())
		return this.respondeVote(state, &sender, false)
	}

	if !(*state).MoreRecentLog(this.pMex.GetLastLogIndex(), this.pMex.GetLastLogTerm()) {
		log.Printf("request vote: log not recent enough")
		return this.respondeVote(state, &sender, false)
	} else if myVote == "" || myVote == this.pMex.GetCandidateId() {
		log.Printf("request vote: vote accepted, voting for: %v", this.pMex.GetCandidateId())
		this.respondeVote(state, &sender, true)
		(*state).VoteFor(sender)
	}

	return this.respondeVote(state, &sender, true)
}

func (this *RequestVoteRPC) respondeVote(state *raftstate.State, sender *string, vote bool) *rpcs.Rpc {
	panic("non implemented")
}
