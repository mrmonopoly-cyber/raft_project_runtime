package clientrequest

import (
	"raft/client/raft-rpcProtobuf-messages/rpcEncoding/out/protobuf"
	"google.golang.org/protobuf/proto"
	rpc "raft/client/src/internal/rpcs"
)

type ClientRequest struct {
  pMex protobuf.ClientReq 
}

func NewClientRequest(parameters map[string]string) rpc.Rpc  {
  var op protobuf.Operation = protobuf.Operation_READ
  var fileName []byte = []byte(parameters["fileName"])
  var others []byte = nil

  switch parameters["operation"] {
    case "Read": 
      op = protobuf.Operation_READ
    case "Write":
      op = protobuf.Operation_WRITE
      others = []byte(parameters["addParam"])
    case "Rename":
      op = protobuf.Operation_RENAME  
      others = []byte(parameters["addParam"])
    case "Delete":
      op = protobuf.Operation_DELETE
    case "Create":
      op = protobuf.Operation_CREATE 
      others = []byte(parameters["addParam"])
    default: 
      op = protobuf.Operation_READ
  }

  return &ClientRequest{
     pMex: protobuf.ClientReq{
      Op: op,
      FileName: fileName,
      Others: others,
    },
  }
}

func (this *ClientRequest) ToString() string {
  return " "
} 

func (this *ClientRequest) Encode() ([]byte, error) {
    return proto.Marshal(&(this).pMex)
}

func (this *ClientRequest) Decode(rawMex []byte) error {
  return proto.Unmarshal(rawMex, &this.pMex)
} 


