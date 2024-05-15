package nodematchidx

import (
    "sync"
    "raft/internal/node/nodeState"
)

type INDEX uint8
const(
    MATCH INDEX = iota
    NEXT INDEX = iota
)

type NodeCommonMatch interface{
    GetNotifyChannel() chan int
    AddNode(ip string)
    RemNode(ip string)
    UpdateNodeState(ip string, indexType INDEX, value int) error
    GetNodeState(ip string) (nodeState.VolatileNodeState,error)
    GetNodeIndex(ip string, indexType INDEX) (int,error)
}

func NewNodeCommonMatch() NodeCommonMatch{
    return &commonMatchNode{
    	notifyChann: make(chan int),
    	allNodeStates:   sync.Map{},
    	numNode:     1,
    	commonIdx:   -1,
    	numStable:   1,
    }
}
