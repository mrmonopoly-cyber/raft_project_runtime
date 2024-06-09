package commonmatch

import (
	nodestate "raft/internal/raftstate/confPool/NodeIndexPool/nodeState"
	"raft/internal/utiliy"

	"github.com/fatih/color"
)

type substriber = utiliy.Triple[nodestate.NodeState, <-chan int, int]


type commonMatchImp struct {
	commitEntryC     chan int
	subs             []substriber
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

func NewCommonMatchImp(initialCommonCommitIdx int, nodeSubs map[string]nodestate.NodeState) *commonMatchImp {
	var res = &commonMatchImp{
		subs:             nil,
		commonMatchIndex: initialCommonCommitIdx,
		numStable:        1, //INFO: Leader always stable
		numNodes:         0,
		commitEntryC:     make(chan int),
        run:              true,
	}

	var nodeSubsNum = len(nodeSubs)

	res.subs = make([]utiliy.Triple[nodestate.NodeState, <-chan int, int], nodeSubsNum)

    var j uint = 0
    for _, v := range nodeSubs {
        if v.GetVoteRight(){
            res.addUpdatedNode(j,v)
        }else{
            go res.checkWhenNodeIsUpdated(j,v)
        }

        res.numNodes++
        j++
    }
	return res
}

//utility

//INFO: subs to see when the node is updated
func (c *commonMatchImp) checkWhenNodeIsUpdated(indexToStoreSub uint, state nodestate.NodeState){
    var _,subC =  state.Subscribe(nodestate.MATCH)
    for {
        var match = <- subC
        if match >= c.commonMatchIndex{
            break           
        }
    }
    color.Green("node updated in commonMatch, adding it the majority count: \n")
    c.addUpdatedNode(indexToStoreSub,state)
    //TODO: unsubscribe

}

func (c *commonMatchImp) addUpdatedNode(indexTrpCell uint, state nodestate.NodeState){
    var trp = &c.subs[indexTrpCell]

    trp.Fst = state
    _,trp.Snd = state.Subscribe(nodestate.MATCH)
    trp.Trd = state.FetchData(nodestate.MATCH)

	go c.checkUpdateNewMatch(trp)
}
func (c *commonMatchImp) checkUpdateNewMatch(sub *substriber) {
    var halfNodeNum = c.numNodes/2

    for c.run{
        var newMatch = <-sub.Snd
        color.Red("check if can increase commonMatchIdx: %v,%v,%v,%v",
        c.numNodes, newMatch, c.commonMatchIndex, sub.Trd)
        if newMatch > c.commonMatchIndex && sub.Trd < newMatch {
            color.Red("check passed: %v:%v:%v:%v\n",
            c.numNodes, newMatch, c.commonMatchIndex, sub.Trd)
            c.numStable++
            if c.numStable > uint(halfNodeNum) {
                if c.commitEntryC != nil{
                    c.commitEntryC <- c.commonMatchIndex
                }
                c.commonMatchIndex++
                c.numStable = 1
            }
        }
        sub.Trd = newMatch
    }
}
