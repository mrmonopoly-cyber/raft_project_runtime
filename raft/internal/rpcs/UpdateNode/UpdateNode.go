package UpdateNode

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

type UpdateNode struct {
    pMex protobuf.UpdateNode
}

func ChangeVoteRightNode(voteAble bool) rpcs.Rpc {
    return &UpdateNode{
        pMex: protobuf.UpdateNode{
            Votante: voteAble,
        },
    }
}

// Manage implements rpcs.Rpc.
func (this *UpdateNode) Execute(
            intLog raft_log.LogEntry,
            metadata clustermetadata.ClusterMetadata,
            confMetadata confmetadata.ConfMetadata,
            senderState nodestate.NodeState) rpcs.Rpc {
    metadata.VoteRight(this.pMex.Votante)
    return nil
}

// ToString implements rpcs.Rpc.
func (this *UpdateNode) ToString() string {
    if this.pMex.Votante {
        return "votante: TRUE"
    }
    return "votante: FALSE"
}

func (this *UpdateNode) Encode() ([]byte, error) {
    var mess []byte
    var err error

    mess, err = proto.Marshal(&(*this).pMex)
    if err != nil {
        log.Panicln("error in Encoding Request Vote: ", err)
    }

	return mess, err
}
func (this *UpdateNode) Decode(b []byte) error {
	err := proto.Unmarshal(b,&this.pMex)
    if err != nil {
        log.Panicln("error in Encoding Request Vote: ", err)
    }
	return err
}
