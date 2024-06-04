package confpool

import (
	"raft/internal/node"
	"raft/internal/raft_log"
	"sync"
)

type OP uint8
const(
    ADD OP = iota
    REM OP = iota
)

type ConfPool interface{
    AutoCommitSet(status bool)
    UpdateNodeList(op OP, node node.Node)
    GetConf() []string
    GetNodeList() *sync.Map
    GetNode(ip string) (node.Node,error)
    raft_log.LogEntry
}

func NewConfPoll(rootDirFs string) ConfPool{
    return confPoolImpl(rootDirFs)
}
