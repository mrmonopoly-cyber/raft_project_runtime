package nodematchidx

import (
	"raft/internal/node/nodeState"
	"sync"
)

type INDEX uint8
const(
    MATCH INDEX = iota
    NEXT INDEX = iota
)

type OPERATION uint8
const(
    INC OPERATION = iota
    DEC OPERATION = iota
)

type NodeCommonMatch interface{
    GetNotifyChannel() chan int
    AddNode(ip string)
    RemNode(ip string)
    UpdateNodeState(ip string, indexType INDEX, value int) error
    GetNodeState(ip string) (nodeState.VolatileNodeState,error)
    GetNodeIndex(ip string, indexType INDEX) (int,error)
    InitVolatileState(ip string, lastLogIndex int)
    Updated(ip string) bool
    DoneUpdating(ip string)
    IncreaseCommonMathcIndex()
    InitCommonMatch(commonMatchIndex int)
    ChangeNnuNodes(op OPERATION)
}

func NewNodeCommonMatch() NodeCommonMatch{
    return &commonMatchNode{
    	notifyChannNewEntry: make(chan int),
    	allNodeStates:   sync.Map{},
    	numNode:     0,
    	futureCommonIdx:   -1,
    	numStable:   0,
    }
}
