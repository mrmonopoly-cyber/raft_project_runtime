package extra

import (
    p "raft/pkg/protobuf"
	"raft/internal/messages"
	"google.golang.org/protobuf/proto"
)

type NEW_RPC struct
{
    term uint64
}

func new_NEW_RPC(term uint64) messages.Rpc{
    return &NEW_RPC{term:term}
}

func (this NEW_RPC) GetTerm() uint64{
    return this.term
}
func (this NEW_RPC) Encode() ([]byte, error){
    return nil,nil
}
func (this NEW_RPC) Decode(b []byte) error{
    return nil
}

