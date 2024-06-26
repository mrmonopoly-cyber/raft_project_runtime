package nodeIndexPool

import "raft/internal/raftstate/confPool/NodeIndexPool/nodeState"

type OP uint8
const (
    ADD OP = iota
    REM OP = iota
)

type NodeIndexPool interface{
    UpdateStatusList(op OP, ip string)
    FetchNodeInfo(ip string) (nodestate.NodeState,error)
}

func NewLeaederCommonIdx() NodeIndexPool{
    return newNodeIndexPoolImpl()
}
