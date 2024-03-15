package AppendEntryResponse

import (
	"google.golang.org/protobuf/proto"
	"raft/internal/messages"
	p "raft/pkg/protobuf"
)

type AppendEntryResponse struct {
	success       bool
	term          uint64
	logIndexError uint64
}

// ToMessage implements messages.Rpc.
func (this *AppendEntryResponse) ToMessage() messages.Message {
	panic("unimplemented")
}

// ToString implements messages.Rpc.
func (this *AppendEntryResponse) ToString() string {
	panic("unimplemented")
}

func new_AppendEntryResponse(success bool, term uint64,
	logIndexError uint64) messages.Rpc {
	return &AppendEntryResponse{
		success:       success,
		term:          term,
		logIndexError: logIndexError,
	}
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
