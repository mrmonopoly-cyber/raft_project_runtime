package inforesponse

import (
	"raft/client/raft-rpcProtobuf-messages/rpcEncoding/out/protobuf"
	"google.golang.org/protobuf/proto"
	rpc "raft/client/src/internal/rpcs"
)

type InfoResponse struct {
  pMex protobuf.InfoResponse 
}

func NewInfoResonse() rpc.Rpc  {
  var req = &InfoResponse{
  }
  return req
}

func (this *InfoResponse) ToString() string {
  return " "
} 

func (this *InfoResponse) Encode() ([]byte, error) {
    return proto.Marshal(&(this).pMex)
}

func (this *InfoResponse) Decode(rawMex []byte) error {
  return proto.Unmarshal(rawMex, &this.pMex)
} 

