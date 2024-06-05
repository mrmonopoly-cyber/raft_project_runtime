package nodestate

import (
	"log"
	"sync"
)

type nodeStateImpl struct {
    counter uint
    subs sync.Map
	mathcIndex int
	nextIndex  int
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
	case NEXTT:
		n.nextIndex = val
    default:
        log.Panicln("case not managed: ", info)
	}
    log.Println("notifing all subscribers of the change")
    n.subs.Range(func(key, value any) bool {
        var C chan int = value.(chan int)
        C <- val
        return true
    })
}

// Substribe implements NodeState.
func (n *nodeStateImpl) Substribe() <-chan int {
    var notifC chan int = make(chan int)

    n.subs.Store(n.counter,notifC)
    n.counter++

    return notifC
}

func newNodeStateImpl() *nodeStateImpl {
	return &nodeStateImpl{
		mathcIndex: -1,
		nextIndex:  -1,
        subs: sync.Map{},
        counter: 0,
	}
}
