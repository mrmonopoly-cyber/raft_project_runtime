package confpool

import (
	"raft/internal/node"
	"raft/internal/raft_log"
	clustermetadata "raft/internal/raftstate/clusterMetadata"
	"sync"
)

type OP uint8

const (
	ADD OP = iota
	REM OP = iota
)

type ConfPool interface {
	UpdateNodeList(op OP, node node.Node)
	GetConf() []string
	GetNodeList() *sync.Map
	GetNode(ip string) (node.Node, error)
	raft_log.LogEntry
}

func NewConfPoll(rootDirFs string,commonMetadata clustermetadata.ClusterMetadata) ConfPool {
    return confPoolImpl(rootDirFs,commonMetadata)
}
