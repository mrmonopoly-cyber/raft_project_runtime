package AppendResponse

import (
	"raft/internal/raftstate"
	"raft/internal/rpcs"
  "strconv"
	"raft/pkg/rpcEncoding/out/protobuf"
	//
	"google.golang.org/protobuf/proto"
)

type AppendResponse struct {
	pMex protobuf.AppendEntryResponse
}

func NewAppendResponseRPC(id string, success bool, term uint64, logIndexError int) rpcs.Rpc {
	var error *int32

	if logIndexError == -1 {
		error = nil
	} else {
		error = proto.Int32(int32(logIndexError))
	}

	return &AppendResponse{
		pMex: protobuf.AppendEntryResponse{
			Id:            &id,
			Term:          &term,
			Success:       &success,
			LogIndexError: error,
		},
	}
}

// GetId implements rpcs.Rpc.
func (this *AppendResponse) GetId() string {
	return this.pMex.GetId()
}

// Manage implements rpcs.Rpc.
func (this *AppendResponse) Execute(state *raftstate.State) *rpcs.Rpc {
	var resp *rpcs.Rpc = nil 

  if !this.pMex.GetSuccess() {
    if this.GetTerm() > (*state).GetTerm() {
      (*state).SetTerm(this.GetTerm())
      (*state).BecomeFollower()
    } else {
      (*state).SetNextIndex(this.GetId(), this.GetLogIndexError())
    }
  } else {
    (*state).SetNextIndex(this.GetId(), (*state).GetLastLogIndex()+1)
    (*state).SetMatchIndex(this.GetId(), (*state).GetLastLogIndex())
  }

  return resp
}

// ToString implements rpcs.Rpc.
func (this *AppendResponse) ToString() string {	
	return "{term : " + strconv.Itoa(int(this.pMex.GetTerm())) + ", id: " + this.pMex.GetId() + ", success: " + strconv.FormatBool(this.pMex.GetSuccess()) + ", error: " + strconv.Itoa(int(this.pMex.GetLogIndexError())) + "}"
}

func (this *AppendResponse) GetTerm() uint64 {
	return this.pMex.GetTerm()
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

func (this *AppendResponse) GetLogIndexError() int {
  return int(this.pMex.GetLogIndexError())
}
