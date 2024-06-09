package RequestVoteRPC

import (
	"log"
	"raft/internal/raft_log"
	clustermetadata "raft/internal/raftstate/clusterMetadata"
	nodestate "raft/internal/raftstate/confPool/NodeIndexPool/nodeState"
	"raft/internal/rpcs"
	"raft/internal/rpcs/RequestVoteResponse"
	"raft/pkg/raft-rpcProtobuf-messages/rpcEncoding/out/protobuf"
	"strconv"

	"google.golang.org/protobuf/proto"
)

type RequestVoteRPC struct {
	pMex protobuf.RequestVote
}

func NewRequestVoteRPC(metadata clustermetadata.ClusterMetadata, intLog raft_log.LogEntry) rpcs.Rpc {
    return &RequestVoteRPC{
        pMex : protobuf.RequestVote{
            Term: metadata.GetTerm(),
            CandidateId: metadata.GetMyIp(clustermetadata.PRI),
            LastLogIndex: int64(intLog.LastLogIndex()),
            LastLogTerm: uint64(intLog.LastLogTerm()),
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

func (this *RequestVoteRPC) Execute(intLog raft_log.LogEntry,
            metadata clustermetadata.ClusterMetadata,
            senderState nodestate.NodeState) rpcs.Rpc{
	var myVote string = metadata.GetVoteFor()
	var senderIp = this.pMex.GetCandidateId()

	if !metadata.CanVote() {
		log.Printf("request vote: this node cannot vote right now")
		return nil
	}

    /*check validity of the vote:
        term >= this.term
        votedFor == NULL OR votedFor == candidateId,
        lastLogIdx >= this.lastLogIdx,
        lastLogTerm >= this.lastLogTerm,
    */

    log.Printf("my vote: %v, candidate id: %v\n", myVote, this.pMex.CandidateId)
    log.Printf("my term: %v, candidate term: %v\n", metadata.GetTerm(), this.pMex.Term)
    log.Printf("my laslogidx: %v, candidate lastlogidx: %v\n", 
    intLog.LastLogIndex(), this.pMex.LastLogIndex)
    log.Printf("my lastlogTerm: %v, candidate lastlotTerm: %v\n", 
    intLog.LastLogTerm(), this.pMex.GetLastLogTerm())

    
    if (this.pMex.Term >= metadata.GetTerm()) && (myVote == "" || myVote == this.pMex.CandidateId) &&
        (this.pMex.LastLogIndex >= int64(intLog.LastLogIndex())) && 
        (this.pMex.LastLogTerm >= uint64(intLog.LastLogTerm())){
            metadata.VoteFor(this.pMex.CandidateId)
            log.Println("vote accepted")
            return this.respondeVote(metadata,&senderIp,true)
    }

    log.Println("vote rejected")
	return this.respondeVote(metadata, &senderIp, false)
}

func (this *RequestVoteRPC)respondeVote(state clustermetadata.ClusterMetadata, sender *string, vote bool) rpcs.Rpc{
    var resp = RequestVoteResponse.NewRequestVoteResponseRPC(*sender,vote, state.GetTerm())
    return resp
}
