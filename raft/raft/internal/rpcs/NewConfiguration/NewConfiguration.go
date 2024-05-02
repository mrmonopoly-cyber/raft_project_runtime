package NewConfiguration

import (
	"log"
	"raft/internal/node/nodeState"
	"raft/internal/raftstate"
	"raft/internal/rpcs"
	"raft/pkg/raft-rpcProtobuf-messages/rpcEncoding/out/protobuf"

	"google.golang.org/protobuf/proto"
)

type NewConfiguration struct {
    pMex protobuf.NewConf
}

func NewNewConfigurationRPC(newConf []string) rpcs.Rpc {
    return &NewConfiguration{
        pMex: protobuf.NewConf{
            NewConfiguration: newConf,
        },
    }
}

// Manage implements rpcs.Rpc.
func (this *NewConfiguration) Execute(state *raftstate.State, senderState *nodeState.VolatileNodeState) *rpcs.Rpc {
    (*state).OverwriteConf(this.pMex.NewConfiguration)
    return nil
}

// ToString implements rpcs.Rpc.
func (this *NewConfiguration) ToString() string {
    panic("dummy implementation")
}

func (this *NewConfiguration) Encode() ([]byte, error) {
    var mess []byte
    var err error

    mess, err = proto.Marshal(&(*this).pMex)
    if err != nil {
        log.Panicln("error in Encoding Request Vote: ", err)
    }

	return mess, err
}
func (this *NewConfiguration) Decode(b []byte) error {
	err := proto.Unmarshal(b,&this.pMex)
    if err != nil {
        log.Panicln("error in Encoding Request Vote: ", err)
    }
	return err
}
