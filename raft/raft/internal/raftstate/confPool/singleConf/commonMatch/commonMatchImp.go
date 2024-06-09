package commonmatch

import (
	"log"
	nodestate "raft/internal/raftstate/confPool/NodeIndexPool/nodeState"
	"raft/internal/utiliy"

	"github.com/fatih/color"
)

type commonMatchImp struct {
	commitEntryC     chan int
	subs             []utiliy.Triple[nodestate.NodeState, <-chan int, int]
	numNodes         int
	numStable        uint
	commonMatchIndex int
    run bool
}

// StopNotify implements CommonMatch.
func (c *commonMatchImp) StopNotify() {
    close(c.commitEntryC)
    c.commitEntryC = nil
    c.run = false
}

// CommitNewEntryC implements CommonMatch.
func (c *commonMatchImp) CommitNewEntryC() <-chan int {
	return c.commitEntryC
}

// utility
func (c *commonMatchImp) updateCommonMatchIndex() {
	var halfNodeNum = c.numNodes / 2

	for _, v := range c.subs {
		go func() {
			for c.run{
				var newMatch = <-v.Snd
				color.Red("check if can increase commonMatchIdx: %v,%v,%v,%v",
					c.numNodes, newMatch, c.commonMatchIndex, v.Trd)
				if newMatch > c.commonMatchIndex && v.Trd < newMatch {
					color.Red("check passed: %v:%v:%v:%v\n",
						c.numNodes, newMatch, c.commonMatchIndex, v.Trd)
					c.numStable++
					if c.numStable > uint(halfNodeNum) {
                        if c.commitEntryC != nil{
                            c.commitEntryC <- c.commonMatchIndex
                        }
						c.commonMatchIndex++
						c.numStable = 1
					}
				}
				v.Trd = newMatch
			}
		}()
	}
}

func NewCommonMatchImp(initialCommonCommitIdx int, nodeSubs map[string]nodestate.NodeState) *commonMatchImp {
	var res = &commonMatchImp{
		subs:             nil,
		commonMatchIndex: initialCommonCommitIdx,
		numStable:        1, //INFO: Leader always stable
		numNodes:         len(nodeSubs),
		commitEntryC:     make(chan int),
        run:              true,
	}

	var nodeSubsNum = len(nodeSubs)

	res.subs = make([]utiliy.Triple[nodestate.NodeState, <-chan int, int], nodeSubsNum)

    var j = 0
    for _, v := range nodeSubs {
        if v.GetVoteRight(){
            var trp = &res.subs[j]

            (*trp).Fst = v
            _,(*trp).Snd = v.Subscribe(nodestate.MATCH)
            (*trp).Trd = v.FetchData(nodestate.MATCH)
        }else{
            //TODO: subs to see when the node is updated
            // go res.checkWhenNodeIsUpdated(v)
            res.numNodes--
        }
    }

	log.Println("list subs commmon Match: ", res.subs)

	go res.updateCommonMatchIndex()

	return res
}

//utility
func (c *commonMatchImp) checkWhenNodeIsUpdated(state nodestate.NodeState){
    var _,subC =  state.Subscribe(nodestate.MATCH)

    for {
        var match = <- subC
        if match >= c.commonMatchIndex{
            c.numNodes++
            break
        }
        //TODO: unsubscribe
    }

}
