package AppendResponse

import (
	"log"
	"raft/internal/node"
	"raft/internal/raftstate"
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
func (this *AppendResponse) Execute(state raftstate.State, sender node.Node) rpcs.Rpc {
    var resp rpcs.Rpc = nil
    var term uint64 = this.pMex.GetTerm()
    var err error 

    if !this.pMex.GetSuccess() {
        if term > state.GetTerm() {
            state.SetTerm(term)
            state.SetRole(raftstate.FOLLOWER)
        } else {
            log.Println("consistency fail")
            //log.Println(this.pMex.GetLogIndexError())
            sender.SetNextIndex(int(this.pMex.GetLogIndexError()))
            //log.Println((*senderState).GetNextIndex())
        }
    } else {
        log.Printf("response ok setting match index of node %v : %v\n", this.pMex.GetId(), this.pMex.GetLogIndexError())
        err = sender.SetNextIndex(int(this.pMex.GetLogIndexError())+1)
        if err != nil {
            log.Panicln("error updating nextIndex: ", err)
        }
        sender.SetMatchIndex(int(this.pMex.GetLogIndexError()))
        if err != nil {
            log.Panicln("error updating matchIndex: ", err)
        }
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
