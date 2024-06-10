package RequestVoteResponse

import (
	"log"
	"raft/internal/raft_log"
	clustermetadata "raft/internal/raftstate/clusterMetadata"
	nodestate "raft/internal/raftstate/confPool/NodeIndexPool/nodeState"
	confmetadata "raft/internal/raftstate/confPool/singleConf/confMetadata"
	"raft/internal/rpcs"
	"raft/pkg/raft-rpcProtobuf-messages/rpcEncoding/out/protobuf"

	"github.com/fatih/color"
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
func (this *RequestVoteResponse) Execute(
                    intLog raft_log.LogEntry,
                    metadata clustermetadata.ClusterMetadata,
                    confMetadata confmetadata.ConfMetadata,
                    senderState nodestate.NodeState) rpcs.Rpc{
    if this.pMex.Term > metadata.GetTerm(){
        metadata.SetRole(clustermetadata.FOLLOWER)
        return nil
    }

    if this.GetVote() {
        log.Println("received positive vote");
        metadata.UpdateSupportersNum(clustermetadata.INC)
    }else {
        log.Println("received negative vote");
        metadata.UpdateSupportersNum(clustermetadata.DEC)
    }
    
    var nodeInCluster = confMetadata.GetNumNodesInConf()
    var nVictory = nodeInCluster/2
    var supp = metadata.GetNumSupporters()
    var notSupp = metadata.GetNumNotSupporters()

    if supp > nVictory {
        color.Green("election won")
        metadata.SetRole(clustermetadata.LEADER)
        return nil
    }
    if supp + notSupp == nodeInCluster{
        color.Red("election lost")
        metadata.ResetElection()
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

