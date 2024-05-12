package nodematchidx

import (
	"errors"
	"log"
	"raft/internal/node/nodeState"
	"sync"
)

type commonMatchNode struct {
	lock          sync.RWMutex
	notifyChann   chan int
	allNodeStates sync.Map
	numNode       uint
	commonIdx     int
	numStable     int
}

// GetNodeState implements NodeCommonMatch.
func (c *commonMatchNode) GetNodeState(ip string) (nodeState.VolatileNodeState,error) {
    var nodeStatePriv,err = c.findNode(ip)
    if err != nil {
        return nil,err
    }
    return nodeStatePriv,nil
}

// GetNodeState implements NodeCommonMatch.
func (c *commonMatchNode) GetNodeIndex(ip string, indexType INDEX) (int, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	var err error
	var nodeStatePriv nodeState.VolatileNodeState
	nodeStatePriv, err = c.findNode(ip)
	if err != nil {
		return -10, err
	}

	switch indexType {
	case MATCH:
		return int(nodeStatePriv.GetMatchIndex()), nil
	case NEXT:
		return nodeStatePriv.GetNextIndex(), nil
	}
	return -10, errors.New("invalid operation type " + string(indexType))
}

// UpdateNodeState implements NodeCommonMatch.
func (c *commonMatchNode) UpdateNodeState(ip string, indexType INDEX, value int) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	var err error
	var nodeStatePriv nodeState.VolatileNodeState
	var matchIdx int
	nodeStatePriv, err = c.findNode(ip)
	if err != nil {
		return err
	}

	switch indexType {
	case NEXT:
		nodeStatePriv.SetNextIndex(value)
	case MATCH:
        //FIX: incomplete implementation
		matchIdx = nodeStatePriv.GetMatchIndex()
		nodeStatePriv.SetMatchIndex(value)
		if matchIdx < c.commonIdx && value >= c.commonIdx {
			c.numStable++
			if c.numStable > int(c.numNode)/2 {
				c.commonIdx++
				c.numStable = 0
				c.notifyChann <- c.commonIdx
			}
		}
	}

	return nil
}

// AddNode implements NodeCommonMatch.
func (c *commonMatchNode) AddNode(ip string) {
	c.lock.Lock()
	defer c.lock.Unlock()
    
    var newState nodeState.VolatileNodeState = nodeState.NewVolatileState()

	c.allNodeStates.Store(ip, newState)
	c.numNode++
}

// RemoNode implements NodeCommonMatch.
func (c *commonMatchNode) RemNode(ip string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.allNodeStates.Delete(ip)
	c.numNode--
}

// GetNotifyChannel implements NodeCommonMatch.
func (c *commonMatchNode) GetNotifyChannel() chan int {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.notifyChann
}

// utility
func (c *commonMatchNode) findNode(ip string) (nodeState.VolatileNodeState, error) {
	var v any
	var found bool
	var nodeStatePriv nodeState.VolatileNodeState

	v, found = c.allNodeStates.Load(ip)
	if !found {
		log.Panicln("state not found for node ", ip)
		return nil, errors.New("state not found for node " + ip)
	}

	nodeStatePriv = v.(nodeState.VolatileNodeState)

	return nodeStatePriv, nil
}
