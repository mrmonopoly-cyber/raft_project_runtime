package AppendResponse

import (
	"log"
	"raft/internal/raftstate"
	"raft/internal/rpcs"
	"raft/internal/node/nodeState"
	"raft/pkg/rpcEncoding/out/protobuf"
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
func (this *AppendResponse) Execute(state *raftstate.State, senderState *nodeState.VolatileNodeState) *rpcs.Rpc {
	var resp *rpcs.Rpc = nil
log.Println(senderState)
	var term uint64 = this.pMex.GetTerm()
	if !this.pMex.GetSuccess() {
		if term > (*state).GetTerm() {
			(*state).SetTerm(term)
			(*state).BecomeFollower()
		} else {
			(*senderState).SetNextIndex(int(this.pMex.GetLogIndexError()))
		}
	} else {
    log.Println((*senderState).GetNextIndex())
		(*senderState).SetNextIndex((*state).GetLastLogIndex()+1)
		(*senderState).SetMatchIndex((*state).GetLastLogIndex())
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
