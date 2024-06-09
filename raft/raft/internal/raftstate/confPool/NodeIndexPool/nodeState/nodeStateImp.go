package nodestate

import (
	"log"
	"sync"

	"github.com/fatih/color"
)

type nodeStateImpl struct {
	counter    uint
	subsNxt    sync.Map
	subsMtc    sync.Map
	nodeIp     string
	mathcIndex int
	nextIndex  int
	voteRight  bool
}

// SetVoteRight implements NodeState.
func (n *nodeStateImpl) SetStatusUpdated() {
    n.voteRight = true
}

// GetVoteRight implements NodeState.
func (n *nodeStateImpl) GetVoteRight() bool{
    return n.voteRight 
}

// FetchData implements NodeState.
func (n *nodeStateImpl) FetchData(info INFO) int {
	switch info {
	case MATCH:
		return n.mathcIndex
	case NEXTT:
		return n.nextIndex
	}
	log.Panicln("case not managed: ", info)
	return -1
}

// UpdateNodeState implements NodeState.
func (n *nodeStateImpl) UpdateNodeState(info INFO, val int) {
	switch info {
	case MATCH:
		n.mathcIndex = val
		color.Cyan("updating match index: ", val)
		n.subsMtc.Range(func(key, value any) bool {
			//HACK: until i resolve the main bug with the stuck node
			go func() {
				var C chan int = value.(chan int)
				C <- val
			}()
			return true
		})
	case NEXTT:
		n.nextIndex = val
		color.Cyan("updating next index: ", val)
		n.subsNxt.Range(func(key, value any) bool {
			//HACK: until i resolve the main bug with the stuck node
			go func() {
				var C chan int = value.(chan int)
				C <- val
			}()
			return true
		})

	default:
		log.Panicln("case not managed: ", info)
	}

}

// Substribe implements NodeState.
func (n *nodeStateImpl) Subscribe(info INFO) (int, <-chan int) {
	var notifC chan int = make(chan int)

	switch info {
	case MATCH:
		n.subsMtc.Store(n.counter, notifC)
	case NEXTT:
		n.subsNxt.Store(n.counter, notifC)
	default:
		log.Panicln("unmanaged case: ", info)
	}
	n.counter++

	return int(n.counter - 1), notifC
}

func newNodeStateImpl(nodeIp string) *nodeStateImpl {
	return &nodeStateImpl{
		mathcIndex: -1,
		nextIndex:  0,
		subsNxt:    sync.Map{},
		subsMtc:    sync.Map{},
		counter:    0,
		nodeIp:     nodeIp,
        voteRight: false,
	}
}
