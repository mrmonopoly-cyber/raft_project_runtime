package nodestate

import "log"

type nodeStateImpl struct {
	mathcIndex int
	nextIndex  int
}

// FetchData implements NodeState.
func (n *nodeStateImpl) FetchData(info INFO) int {
    switch info{
    case MATCH:
        return n.mathcIndex
    case NEXTT:
        return n.nextIndex
    }
    log.Panicln("case not managed: ",info)
    return -1
}

// UpdateNodeState implements NodeState.
func (n *nodeStateImpl) UpdateNodeState(info INFO, val int) {
    switch info{
    case MATCH:
        n.mathcIndex = val
        return
    case NEXTT:
        n.nextIndex = val
        return
    }
    log.Panicln("case not managed: ",info)
}

func newNodeStateImpl() *nodeStateImpl {
	return &nodeStateImpl{
		mathcIndex: -1,
		nextIndex:  -1,
	}
}
