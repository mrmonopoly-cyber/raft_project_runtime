package message

import (
	"encoding/json"
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

func (m *Message) ToByte() []byte {
  result, _ := json.Marshal(m)
  return result
}

func (m *Message) ToMessage(b []byte) {
  json.Unmarshal(b, m)
}
