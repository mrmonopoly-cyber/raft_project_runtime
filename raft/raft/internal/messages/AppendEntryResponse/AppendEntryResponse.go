package AppendEntryResponse

import (
	"raft/internal/messages"
	"raft/internal/raftstate"
	p "raft/pkg/protobuf"
	"strconv"

	"google.golang.org/protobuf/proto"
)

type AppendEntryResponse struct {
	id            string
	success       bool
	term          uint64
	logIndexError uint64
}

// Manage implements messages.Rpc.
func (this *AppendEntryResponse) Execute(state *raftstate.State) *messages.Rpc {
	if !this.success {

	} 
    return nil
}

// ToString implements messages.Rpc.
func (this *AppendEntryResponse) ToString() string {
	return "{term : " + strconv.Itoa(int(this.term)) + ",\nlogIndexError: " + strconv.Itoa(int(this.logIndexError)) + ", \nsuccess: " + strconv.FormatBool(this.success) + "}"
}

func NewAppendEntryResponse(id string, success bool, term uint64,
	logIndexError uint64) messages.Rpc {
	return &AppendEntryResponse{
		id:            id,
		success:       success,
		term:          term,
		logIndexError: logIndexError,
	}
}

func (this AppendEntryResponse) GetId() string {
	return this.id
}

func (this AppendEntryResponse) GetTerm() uint64 {
	return this.term
}
func (this AppendEntryResponse) Encode() ([]byte, error) {

	response := &p.AppendEntryResponse{
		Success:       proto.Bool(this.success),
		Term:          proto.Uint64(this.term),
		LogIndexError: proto.Uint64(this.logIndexError),
	}

	return proto.Marshal(response)
}
func (this AppendEntryResponse) Decode(b []byte) error {

	pb := new(p.AppendEntryResponse)
	err := proto.Unmarshal(b, pb)

	if err != nil {
		this.term = pb.GetTerm()
		this.success = pb.GetSuccess()
		this.logIndexError = pb.GetLogIndexError()
	}

	return err
}
func (this AppendEntryResponse) HasSucceded() bool {
	return this.success
}
func (this AppendEntryResponse) GetIndex() uint64 {
	return this.logIndexError
}
