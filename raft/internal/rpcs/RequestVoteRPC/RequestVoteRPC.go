package RequestVoteRPC

import (
	"log"
	"raft/internal/node"
	"raft/internal/raftstate"
	"raft/internal/rpcs"
	"raft/internal/rpcs/RequestVoteResponse"
	"raft/pkg/raft-rpcProtobuf-messages/rpcEncoding/out/protobuf"
	"strconv"

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

// GetCandidateId rpcs.Rpc.
func (this *RequestVoteRPC) GetCandidateId() string {
  return this.pMex.CandidateId
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
	err := proto.Unmarshal(b,&this.pMex)
    if err != nil {
        log.Panicln("error in Encoding Request Vote: ", err)
    }
	return err
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
func (this *RequestVoteRPC) Execute(state raftstate.State, sender node.Node) *rpcs.Rpc {
	var myVote string = state.GetVoteFor()
	var senderIp = this.pMex.GetCandidateId()

	if !state.CanVote() {
		log.Printf("request vote: this node cannot vote right now")
		return nil
	}

    /*TODO: check validity of the vote:
        term >= this.term
        votedFor == NULL OR votedFor == candidateId,
        lastLogIdx >= this.lastLogIdx,
        lastLogTerm >= this.lastLogTerm,


    */
    
    if (this.pMex.Term >= state.GetTerm()) && (myVote == "" || myVote == this.pMex.CandidateId) &&
        (this.pMex.LastLogIndex >= int64(state.LastLogIndex())) && (this.pMex.LastLogTerm >= uint64(state.LastLogTerm())){
            log.Println("vote accepted")
            state.VoteFor(this.pMex.CandidateId)
            return this.respondeVote(state,&senderIp,true)
    }

    log.Println("vote rejected")
	return this.respondeVote(state, &senderIp, false)
}

func (this *RequestVoteRPC) respondeVote(state raftstate.State, sender *string, vote bool) *rpcs.Rpc{
    var resp = RequestVoteResponse.NewRequestVoteResponseRPC(*sender,vote, state.GetTerm())
    return &resp
}
