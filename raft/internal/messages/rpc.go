package messages

import (
	p "raft/pkg/protobuf"

	"google.golang.org/protobuf/proto"
)


type MEX_TYPE int32

const (
    APPEND_ENTRY            MEX_TYPE = 0
    APPEND_ENTRY_RESPONSE        = 1
    REQUEST_VOTE                 = 2
    REQUEST_VOTE_RESPONSE        = 3
    COPY_STATE                   = 4
)

type Rpc interface {
  GetTerm() uint64
  Encode() ([]byte, error)
  Decode(b []byte) error
  ToString() string
  ToMessage() Message
}

type Message struct{
    Mex_type MEX_TYPE
    Payload []byte
}

func (this Message) ToByte() []byte  {
  var t p.Ty

  switch this.Mex_type {
    case APPEND_ENTRY:
      t = p.Ty_APPEND_ENTRY
    case APPEND_ENTRY_RESPONSE:
      t = p.Ty_APPEND_RESPONSE
    case REQUEST_VOTE:
      t = p.Ty_REQUEST_VOTE
    case REQUEST_VOTE_RESPONSE:
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
    Mex_type: MEX_TYPE(*mess.Kind),
    Payload: mess.Payload,
  }
}
