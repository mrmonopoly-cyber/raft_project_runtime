package confpool

import (
	"raft/internal/node"
	"raft/internal/raft_log"
	clustermetadata "raft/internal/raftstate/clusterMetadata"
	nodeIndexPool "raft/internal/raftstate/confPool/NodeIndexPool"
	confmetadata "raft/internal/raftstate/confPool/singleConf/confMetadata"
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
    confmetadata.ConfMetadata
}

func NewConfPoll(rootDirFs string,commonMetadata clustermetadata.ClusterMetadata) ConfPool {
    return confPoolImpl(rootDirFs,commonMetadata)
}
