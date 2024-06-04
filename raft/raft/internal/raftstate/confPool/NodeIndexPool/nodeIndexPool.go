package nodeIndexPool

import "raft/internal/raftstate/confPool/NodeIndexPool/nodeState"

type OP uint8
const (
    ADD OP = iota
    REM OP = iota
)

type NodeIndexPool interface{
    UpdateStatusList(op OP, ip string)
    UpdateNodeInfo(info nodestate.INFO, ip string, val int) 
    FetchNodeInfo(info nodestate.INFO, ip string) (int,error)

}

func NewLeaederCommonIdx() NodeIndexPool{
    return newNodeIndexPoolImpl()
}
