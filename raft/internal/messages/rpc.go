package messages

import "encoding/json"

type MEX_TYPE uint

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
    mex,_ := json.Marshal(this)
    return mex
}

func New_message(data []byte) *Message{
    var mex Message
    json.Unmarshal(data,&mex)
    return &mex
}
