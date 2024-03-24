package messages

import (
	"log"
	"raft/internal/raftstate"
)

type Rpc interface {
	GetId() string
	GetTerm() uint64
	ToString() string
	Execute(state *raftstate.State) *Rpc
    Encode() ([]byte, error)
}

func Decode(raw_mex []byte) *Rpc{
    	
    var num = int32(raw_mex[0])<<24 | int32(raw_mex[1])<<16 | int32(raw_mex[2])<<8 | int32(raw_mex[3])
    log.Println("message type :",num)

    panic("unimplemented")
}
