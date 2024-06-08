package commonmatch

import (
	"log"
	nodestate "raft/internal/raftstate/confPool/NodeIndexPool/nodeState"
	"raft/internal/utiliy"

	"github.com/fatih/color"
)


type commonMatchImp struct {
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
            for{
                var newMatch = <- v.Snd
                color.Red("check if can increase commonMatchIdx: %v,%v,%v,%v",
                    c.numNodes,newMatch,c.commonMatchIndex, v.Trd)
                if newMatch >= c.commonMatchIndex && v.Trd < newMatch{
                    color.Red("check passed: %v:%v:%v:%v\n",
                        c.numNodes, newMatch, c.commonMatchIndex, v.Trd)
                    c.numStable++
                    if c.numStable > uint(halfNodeNum){
                        c.commitEntryC <- c.commonMatchIndex
                        c.commonMatchIndex++
                        c.numStable=1
                    }
                }
                v.Trd = newMatch
            }
        }()
    }
}

func NewCommonMatchImp(initialCommonCommitIdx int, nodeSubs []nodestate.NodeState) *commonMatchImp {
	var res = &commonMatchImp{
        subs: nil,
        commonMatchIndex: initialCommonCommitIdx,
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
