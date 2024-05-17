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

type NodeCommonMatch interface{
    GetNotifyChannel() chan int
    GetNotifyChannelOldEntry() chan EntryToSend
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
}

func NewNodeCommonMatch() NodeCommonMatch{
    return &commonMatchNode{
    	notifyChannNewEntry: make(chan int),
    	notifyChannOldEntry: make(chan EntryToSend),
    	allNodeStates:   sync.Map{},
    	numNode:     1,
    	commonIdx:   -1,
    	numStable:   1,
    }
}
