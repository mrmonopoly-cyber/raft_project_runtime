package Redirection

import (
	"log"
	"raft/internal/raft_log"
	clustermetadata "raft/internal/raftstate/clusterMetadata"
	nodestate "raft/internal/raftstate/confPool/NodeIndexPool/nodeState"
	confmetadata "raft/internal/raftstate/confPool/singleConf/confMetadata"
	"raft/internal/rpcs"
	"raft/pkg/raft-rpcProtobuf-messages/rpcEncoding/out/protobuf"

	"google.golang.org/protobuf/proto"
)

type Redirection struct {
    pMex protobuf.Redirect
}

func NewredirectionRPC(leaderIp string) rpcs.Rpc {
    return &Redirection{
        pMex: protobuf.Redirect{
            Leaderip: leaderIp,
        },
    }
}

// Manage implements rpcs.Rpc.
func (this *Redirection) Execute( 
            intLog raft_log.LogEntry,
            metadata clustermetadata.ClusterMetadata,
            confMetadata confmetadata.ConfMetadata,
            senderState nodestate.NodeState)rpcs.Rpc {
    panic("dummy implementation")
}

// ToString implements rpcs.Rpc.
func (this *Redirection) ToString() string {
    panic("dummy implementation")
}

func (this *Redirection) Encode() ([]byte, error) {
    var mess []byte
    var err error

    mess, err = proto.Marshal(&(*this).pMex)
    if err != nil {
        log.Panicln("error in Encoding Request Vote: ", err)
    }

	return mess, err
}
func (this *Redirection) Decode(b []byte) error {
	err := proto.Unmarshal(b,&this.pMex)
    if err != nil {
        log.Panicln("error in Encoding Request Vote: ", err)
    }
	return err
}
