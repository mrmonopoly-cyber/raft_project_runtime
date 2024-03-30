package AppendResponse

import (
	"raft/internal/raftstate"
	"raft/internal/rpcs"
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
func (this *AppendResponse) Execute(state *raftstate.State) *rpcs.Rpc {
	var resp *rpcs.Rpc = nil

	var id string = this.pMex.GetId()
	var term uint64 = this.pMex.GetTerm()
	if !this.pMex.GetSuccess() {
		if term > (*state).GetTerm() {
			(*state).SetTerm(term)
			(*state).BecomeFollower()
		} else {
			(*state).SetNextIndex(id, int(this.pMex.GetLogIndexError()))
		}
	} else {
		(*state).SetNextIndex(id, (*state).GetLastLogIndex()+1)
		(*state).SetMatchIndex(id, (*state).GetLastLogIndex())
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
	var pb = new(protobuf.AppendEntryResponse)
	err := proto.Unmarshal(b, pb)

	if err != nil {
		this.pMex.Term = pb.Term
		this.pMex.Id = pb.Id
		this.pMex.Success = pb.Success
		this.pMex.LogIndexError = pb.LogIndexError
	}

	return err
}
