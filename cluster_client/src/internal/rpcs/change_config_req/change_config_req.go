package changeconfigreq

import (
	"raft/client/raft-rpcProtobuf-messages/rpcEncoding/out/protobuf"
	"google.golang.org/protobuf/proto"
	rpc "raft/client/src/internal/rpcs"
)

type ChangeConfigRequest struct {
  pMex protobuf.ChangeConfReq 
}

func NewConfigRequest(config protobuf.ClusterConf, op protobuf.AdminOp) rpc.Rpc  {
  var req = &ChangeConfigRequest{
    pMex: protobuf.ChangeConfReq{
      Op: op,
      Conf: &config,
    },
  }
  return req
}

func (this *ChangeConfigRequest) ToString() string {
  return " "
} 

func (this *ChangeConfigRequest) Encode() ([]byte, error) {
    return proto.Marshal(&(this).pMex)
}

func (this *ChangeConfigRequest) Decode(rawMex []byte) error {
  return proto.Unmarshal(rawMex, &this.pMex)
} 


