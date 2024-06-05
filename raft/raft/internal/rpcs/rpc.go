package rpcs

import (
	"raft/internal/raft_log"
	clustermetadata "raft/internal/raftstate/clusterMetadata"
	nodestate "raft/internal/raftstate/confPool/NodeIndexPool/nodeState"
)

type Rpc interface {
	ToString() string
	Execute(intLog raft_log.LogEntry,
            metadata clustermetadata.ClusterMetadata,
            senderState nodestate.NodeState) Rpc
	Encode() ([]byte, error)
	Decode(rawMex []byte) error
}
