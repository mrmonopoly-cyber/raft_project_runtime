package commonmatch

import (
	nodestate "raft/internal/raftstate/confPool/NodeIndexPool/nodeState"
	"sync"
)

type commonMatchImp struct {
    wg  sync.WaitGroup
    commitEntryC chan int
    subs []nodestate.NodeState
    numNodes int
    numStable uint
    commonMatchIndex int
}

// CommitNewEntryC implements CommonMatch.
func (c *commonMatchImp) CommitNewEntryC() <-chan int {
    return c.commitEntryC
}

//utility
func (c *commonMatchImp) updateCommonMatchIndex()  {
}

func NewCommonMatchImp(nodeSubs []nodestate.NodeState) *commonMatchImp {
	var res = &commonMatchImp{
        wg: sync.WaitGroup{},
        subs: nodeSubs,
        commonMatchIndex: -1,
        numStable: 0,
        numNodes: len(nodeSubs),
        commitEntryC: make(chan int),
    }

    go res.updateCommonMatchIndex()
    
    return res
}
