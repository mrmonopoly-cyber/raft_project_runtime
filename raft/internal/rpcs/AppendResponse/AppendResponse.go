package AppendResponse

import (
	"log"
	"raft/internal/node"
	"raft/internal/raft_log"
	clustermetadata "raft/internal/raftstate/clusterMetadata"
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
            sender node.Node)rpcs.Rpc {
    var resp rpcs.Rpc = nil
    var term uint64 = this.pMex.GetTerm()

    if !this.pMex.GetSuccess() {
        if term > metadata.GetTerm() {
            metadata.SetTerm(term)
            metadata.SetRole(clustermetadata.FOLLOWER)
        } else {
            log.Println("consistency fail")
            //log.Println(this.pMex.GetLogIndexError())
            // sender.SetNextIndex(int(this.pMex.GetLogIndexError()))
            //log.Println((*senderState).GetNextIndex())
        }
    } else {
        log.Printf("response ok increasing match and next index of node: %v\n", *this.pMex.Id)
        // sender.SetNextIndex(int(this.pMex.GetLogIndexError())+1)
        // sender.SetMatchIndex(int(this.pMex.GetLogIndexError()))
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
