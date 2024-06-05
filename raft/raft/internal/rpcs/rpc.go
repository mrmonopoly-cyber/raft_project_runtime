package rpcs

import (
	"raft/internal/node"
	"raft/internal/raft_log"
	clustermetadata "raft/internal/raftstate/clusterMetadata"
)

type Rpc interface {
	ToString() string
	Execute(intLog raft_log.LogEntry,
            metadata clustermetadata.ClusterMetadata,
            sender node.Node) Rpc
	Encode() ([]byte, error)
	Decode(rawMex []byte) error
}
