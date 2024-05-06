package ClientResponse

import (
	"raft/internal/raftstate"
	"raft/internal/rpcs"
	"raft/internal/node/nodeState"
  "raft/pkg/raft-rpcProtobuf-messages/rpcEncoding/out/protobuf"

	"google.golang.org/protobuf/proto"
)

type ClientResponse struct {
  pMex protobuf.ClientResp
}

func NewClientResponseRPC(success bool, content []byte) rpcs.Rpc {
    return &ClientResponse{
      pMex: protobuf.ClientResp{
          Content: content,
          Success: success,
      },
    }
}

// Manage implements rpcs.Rpc.
func (this *ClientResponse) Execute(state *raftstate.State, senderState *nodeState.VolatileNodeState) *rpcs.Rpc {
    panic("dummy implementation")
}

// ToString implements rpcs.Rpc.
func (this *ClientResponse) ToString() string {
    panic("dummy implementation")
}

func (this *ClientResponse) Encode() ([]byte, error) {
    return proto.Marshal(&(*this).pMex)
}
func (this *ClientResponse) Decode(b []byte) error {
	return proto.Unmarshal(b,&this.pMex)
}
