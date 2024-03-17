package message

import (
	"encoding/json"
	"raft/internal/messages"
	"raft/internal/messages/AppendEntryRPC"
	"raft/internal/messages/AppendEntryResponse"
	"raft/internal/messages/CopyStateRPC"
	"raft/internal/messages/RequestVoteRPC"
	"raft/internal/messages/RequestVoteResponse"
	p "raft/pkg/protobuf"

	"google.golang.org/protobuf/proto"
)

type EnumType int

const (
  APPEND_ENTRY EnumType = iota
  REQUEST_VOTE 
  APPEND_RESPONSE
  VOTE_RESPONSE
  COPY_STATE
)

type Message struct {
  Ty EnumType
  Payload []byte
}

func (m *Message) ToMessage(b []byte) {
  json.Unmarshal(b, m)
}

func (this *Message) ToByte() []byte  {
  var t p.Ty

  switch this.Ty{
    case APPEND_ENTRY:
      t = p.Ty_APPEND_ENTRY
    case APPEND_RESPONSE:
      t = p.Ty_APPEND_RESPONSE
    case REQUEST_VOTE:
      t = p.Ty_REQUEST_VOTE
    case VOTE_RESPONSE:
      t = p.Ty_VOTE_RESPONSE
    case COPY_STATE:
      t = p.Ty_COPY_STATE
  }
  
  mess := &p.Message{
    Kind: t.Enum(),
    Payload: this.Payload,
  } 
    mex,_ := proto.Marshal(mess)
    return mex
}

func NewMessage(data []byte) *Message{
  mess := new(p.Message)
  proto.Unmarshal(data,mess)
    
  return &Message{
    Ty: EnumType(*mess.Kind),
    Payload: mess.Payload,
  }
}


func (this *Message) ToRpc() *messages.Rpc{

    var recRpc messages.Rpc
    switch this.Ty{
    case APPEND_ENTRY:
        var rpc AppendEntryRPC.AppendEntryRPC
        recRpc = &rpc
    case APPEND_RESPONSE:
        var rpc AppendEntryResponse.AppendEntryResponse 
        recRpc = &rpc
    case REQUEST_VOTE:
        var rpc RequestVoteRPC.RequestVoteRPC
        recRpc = &rpc
    case VOTE_RESPONSE:
        var rpc RequestVoteResponse.RequestVoteResponse
        recRpc = &rpc
    case COPY_STATE:
        var rpc CopyStateRPC.CopyStateRPC
        recRpc = &rpc
    default:
        panic("impossible case invalid type RPC " + string(rune(this.Ty)))
    }
    recRpc.Decode(this.Payload)

    return &recRpc
}
