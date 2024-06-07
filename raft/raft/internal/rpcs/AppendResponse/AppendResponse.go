package AppendResponse

import (
	"log"
	"raft/internal/raft_log"
	clustermetadata "raft/internal/raftstate/clusterMetadata"
	nodestate "raft/internal/raftstate/confPool/NodeIndexPool/nodeState"
	"raft/internal/rpcs"
	"raft/pkg/raft-rpcProtobuf-messages/rpcEncoding/out/protobuf"
	"strconv"

	//
	"google.golang.org/protobuf/proto"
)

type AppendResponse struct {
	pMex protobuf.AppendEntryResponse
}

func NewAppendResponseRPC(id string, success bool, term uint64, logIndexError int) rpcs.Rpc {
	var error int32 = int32(logIndexError) 

  return &AppendResponse{
		pMex: protobuf.AppendEntryResponse{
			Id:            &id,
			Term:          &term,
			Success:       &success,
			LogIndexError: &error,
		},
	}
}

// Manage implements rpcs.Rpc.
func (this *AppendResponse) Execute( 
            intLog raft_log.LogEntry,
            metadata clustermetadata.ClusterMetadata,
            senderState nodestate.NodeState)rpcs.Rpc {
    var resp rpcs.Rpc = nil
    var term uint64 = this.pMex.GetTerm()
    var errorIndex = this.pMex.GetLogIndexError()

    if !this.pMex.GetSuccess() {
        if term > metadata.GetTerm() {
            metadata.SetTerm(term)
            metadata.SetRole(clustermetadata.FOLLOWER)
        } else {
            log.Println("consistency fail: ", errorIndex)
            senderState.UpdateNodeState(nodestate.NEXTT,int(errorIndex))
            log.Println("setting next index to: ",senderState.FetchData(nodestate.NEXTT))
        }
    } else {
        log.Printf("response ok increasing match and next index of node: %v\n", *this.pMex.Id)
        senderState.UpdateNodeState(nodestate.NEXTT,int(errorIndex)+1)
        senderState.UpdateNodeState(nodestate.MATCH,int(errorIndex))
    }

    return resp
}

// ToString implements rpcs.Rpc.
func (this *AppendResponse) ToString() string {
	return "{term : " + strconv.Itoa(int(this.pMex.GetTerm())) + ", id: " + this.pMex.GetId() + ", success: " + strconv.FormatBool(this.pMex.GetSuccess()) + ", error: " + strconv.Itoa(int(this.pMex.GetLogIndexError())) + "}"
}

func (this *AppendResponse) Encode() ([]byte, error) {
	var mess []byte
	var err error
	mess, err = proto.Marshal(&(*this).pMex)
	return mess, err
}

func (this *AppendResponse) Decode(b []byte) error {
	err := proto.Unmarshal(b, &this.pMex)
    if err != nil {
        log.Panicln("error in Decoding Append Response: ", err)
    }
	return err
}
