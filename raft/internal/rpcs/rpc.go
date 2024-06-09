package rpcs

import (
	"raft/internal/raft_log"
	clustermetadata "raft/internal/raftstate/clusterMetadata"
	nodestate "raft/internal/raftstate/confPool/NodeIndexPool/nodeState"
	confmetadata "raft/internal/raftstate/confPool/singleConf/confMetadata"
)

type Rpc interface {
	ToString() string
	Execute(
            intLog raft_log.LogEntry,
            metadata clustermetadata.ClusterMetadata,
            confMetadata confmetadata.ConfMetadata,
            senderState nodestate.NodeState) Rpc
	Encode() ([]byte, error)
	Decode(rawMex []byte) error
}
