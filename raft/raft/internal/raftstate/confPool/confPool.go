package confpool

import (
	"raft/internal/node"
	"raft/internal/raft_log"
	clustermetadata "raft/internal/raftstate/clusterMetadata"
	nodeIndexPool "raft/internal/raftstate/confPool/NodeIndexPool"
	"sync"
)

type OP uint8

const (
	ADD OP = iota
	REM OP = iota
)

type ConfPool interface {
	UpdateNodeList(op OP, node node.Node)
	GetNodeList() *sync.Map
	GetNode(ip string) (node.Node, error)
    SendHearthBit()
	raft_log.LogEntry
	nodeIndexPool.NodeIndexPool
}

func NewConfPoll(rootDirFs string,commonMetadata clustermetadata.ClusterMetadata) ConfPool {
    return confPoolImpl(rootDirFs,commonMetadata)
}
