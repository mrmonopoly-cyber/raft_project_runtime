package RequestVoteResponse

import (
	"log"
	"raft/internal/raftstate"
	"raft/internal/rpcs"
	"raft/pkg/rpcEncoding/out/protobuf"
	"google.golang.org/protobuf/proto"
)

type RequestVoteResponse struct {
    pMex protobuf.RequestVoteResponse
}

func NewRequestVoteResponseRPC(id string, vote bool, term uint64) rpcs.Rpc {
    return &RequestVoteResponse{
        pMex: protobuf.RequestVoteResponse{
            Id: id,
            VoteGranted: vote,
            Term: term,
        },
    }
}

// GetId implements rpcs.Rpc.
func (this *RequestVoteResponse) GetId() string {
    return this.pMex.Id
}


// Manage implements rpcs.Rpc.
func (this *RequestVoteResponse) Execute(state *raftstate.State) *rpcs.Rpc {
    if this.GetVote() {
        log.Println("received positive vote");
        (*state).IncreaseSupporters()
    }else {
        log.Println("received negative vote");
        (*state).IncreaseNotSupporters()
    }
    
    var nodeInCluster = (*state).GetNumNodeInCluster()
    var nVictory = nodeInCluster/2
    var supp = (*state).GetNumSupporters()
    var notSupp = (*state).GetNumNotSupporters()

    if supp > nVictory {
        log.Println("election won");
        (*state).SetRole(raftstate.LEADER)
        (*state).ResetElection()
        return nil
    }
    if supp + notSupp == nodeInCluster{
        log.Println("election lost");
        (*state).ResetElection()
    }

    return nil
}

// ToString implements rpcs.Rpc.
func (this *RequestVoteResponse) ToString() string {
    var vote string = "false"
    if this.GetVote() {
        vote = "true"
    }
    var mex string = "{"+ "Id: " + this.GetId() + 
                     ", VoteGranted: " + vote + 
                     ", Term: " + string(rune(this.GetTerm())) +
                     "}"
    log.Println("rpc RequestVoteResponse :", mex)

    return mex
}

func (this *RequestVoteResponse) GetTerm() uint64 {
   return this.pMex.Term
}

func (this *RequestVoteResponse) GetVote() bool{
    return this.pMex.VoteGranted
}

func (this *RequestVoteResponse) Encode() ([]byte, error) {
    var mess []byte
    var err error

    mess, err = proto.Marshal(&(*this).pMex)
    if err != nil {
        log.Panicln("error in Encoding Request Vote Response: ", err)
    }

	return mess, err
}
func (this *RequestVoteResponse) Decode(b []byte) error {
	err := proto.Unmarshal(b,&this.pMex)
    if err != nil {
        log.Panicln("error in Decoding Request Vote Response: ", err)
    }
	return err
}

