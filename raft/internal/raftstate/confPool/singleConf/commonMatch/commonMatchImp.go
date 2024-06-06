package commonmatch

import (
	"log"
	nodestate "raft/internal/raftstate/confPool/NodeIndexPool/nodeState"
	"raft/internal/utiliy"
	"sync"
)


type commonMatchImp struct {
    lock sync.RWMutex
    commitEntryC chan int
    subs []utiliy.Triple[nodestate.NodeState,<- chan int,int]
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
    var halfNodeNum = c.numNodes/2
    
    for _,v  := range c.subs {
        go func(){
            <- v.Snd
            c.lock.Lock()
            var newMatch = v.Fst.FetchData(nodestate.MATCH)
            if newMatch >= c.commonMatchIndex && v.Trd < c.commonMatchIndex {
                c.numStable++
                for c.numStable > uint(halfNodeNum){
                    c.commitEntryC <- 1
                    c.numStable=1
                    //INFO: it's possible that a node has a mach greater match index 
                    //and so every time i increment the current common i have to check if
                    //the other nodes has already a greater common match index
                    for _, v := range c.subs {
                        if v.Trd >= c.commonMatchIndex{
                            c.numStable++
                        }
                        
                    }
                }
            }
            v.Trd = newMatch
            c.lock.Unlock()
        }()
    }
}

func NewCommonMatchImp(nodeSubs []nodestate.NodeState) *commonMatchImp {
	var res = &commonMatchImp{
        lock: sync.RWMutex{},
        subs: nil,
        commonMatchIndex: -1,
        numStable: 1, //INFO: Leader always stable
        numNodes: len(nodeSubs),
        commitEntryC: make(chan int),
    }

    var nodeSubsNum = len(nodeSubs)
    
    res.subs = make([]utiliy.Triple[nodestate.NodeState,<- chan int, int], nodeSubsNum)
    for i := 0; i < nodeSubsNum; i++ {
        res.subs[i].Fst = nodeSubs[i]
        res.subs[i].Snd = nodeSubs[i].Subscribe(nodestate.MATCH)
        res.subs[i].Trd = nodeSubs[i].FetchData(nodestate.MATCH)
    }
    log.Println("list subs commmon Match: ", res.subs)

    go res.updateCommonMatchIndex()
    
    return res
}
