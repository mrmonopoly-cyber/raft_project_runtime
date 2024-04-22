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
  var payload []byte = []byte(parameters["fileName"])

  switch parameters["operation"] {
    case "Read": 
      op = protobuf.Operation_READ
    case "Write":
      op = protobuf.Operation_WRITE
    case "Rename":
      op = protobuf.Operation_RENAME  
    case "Delete":
      op = protobuf.Operation_DELETE
    case "Create":
      op = protobuf.Operation_CREATE 
    default: 
      op = protobuf.Operation_READ
  }

  return &ClientRequest{
     pMex: protobuf.ClientReq{
      Op: op,
      Payload: payload,
    },
  }
}

func (this *ClientRequest) ToString() string {
  return " "
} 

func (this *ClientRequest) Encode() ([]byte, error) {
    var mex []byte
    var err error

    mex, err = proto.Marshal(&(this).pMex)
    if err != nil {
        panic("error encoding")
    }
    return mex, err
}

func (this *ClientRequest) Decode(rawMex []byte) error {
  err := proto.Unmarshal(rawMex, &this.pMex)
  if err != nil {
    panic("error decoding")
  }

  return err
  
} 


