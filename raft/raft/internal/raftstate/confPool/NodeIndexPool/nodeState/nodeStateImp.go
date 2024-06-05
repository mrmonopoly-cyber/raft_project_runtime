package nodestate

import (
	"log"
	"sync"
)

type nodeStateImpl struct {
    counter uint
    subsNxt sync.Map
    subsMtc sync.Map
    nodeIp string
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
        log.Println("notifing all subscribers of the change in match index")
        go n.subsMtc.Range(func(key, value any) bool {
            var C chan int = value.(chan int)
            C <- val
            return true
        })
	case NEXTT:
		n.nextIndex = val
        log.Println("notifing all subscribers of the change in next index")
        go n.subsNxt.Range(func(key, value any) bool {
            var C chan int = value.(chan int)
            C <- val
            return true
        })
    default:
        log.Panicln("case not managed: ", info)
	}

}

// Substribe implements NodeState.
func (n *nodeStateImpl) Subscribe(info INFO) <-chan int {
    var notifC chan int = make(chan int)

    switch info{
    case MATCH:
        n.subsMtc.Store(n.counter,notifC)
    case NEXTT:
        n.subsNxt.Store(n.counter,notifC)
    default:
        log.Panicln("unmanaged case: ",info)
    }
    n.counter++

    return notifC
}

func newNodeStateImpl(nodeIp string) *nodeStateImpl {
	return &nodeStateImpl{
		mathcIndex: -1,
		nextIndex:  -1,
        subsNxt: sync.Map{},
        subsMtc: sync.Map{},
        counter: 0,
        nodeIp: nodeIp,
	}
}
