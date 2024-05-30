package nodematchidx

import (
	"errors"
	"log"
	"raft/internal/node/nodeState"
	"sync"
)

type commonMatchNode struct {
	lock                sync.RWMutex
	notifyChannNewEntry chan int
	allNodeStates       sync.Map
	behindNode          sync.Map
	numNode             uint
	futureCommonIdx           int
	numStable           int
}


// InitCommonMatch implements NodeCommonMatch.
func (c *commonMatchNode) InitCommonMatch(commonMatchIndex int) {
	c.futureCommonIdx = commonMatchIndex
	c.allNodeStates.Range(func(key, value any) bool {
		var nodeState = value.(nodeState.VolatileNodeState)
		nodeState.InitVolatileState(commonMatchIndex)
		return true
	})
}

// IncreaseCommonMathcIndex implements NodeCommonMatch.
func (c *commonMatchNode) IncreaseCommonMathcIndex() {
	c.futureCommonIdx++
}

// DoneUpdating implements NodeCommonMatch.
func (c *commonMatchNode) DoneUpdating(ip string) {
	var nodeState, err = c.findNode(ip)
	if err != nil {
		return
	}
	nodeState.NodeUpdated()
}

// Updated implements NodeCommonMatch.
func (c *commonMatchNode) Updated(ip string) bool {
	var nodeState, err = c.findNode(ip)
	if err != nil {
		return false
	}
	return nodeState.Updated()
}

// InitVolatileState implements NodeCommonMatch.
func (c *commonMatchNode) InitVolatileState(ip string, lastLogIndex int) {
	var nodeState, err = c.findNode(ip)
	if err != nil {
		return
	}
	nodeState.InitVolatileState(lastLogIndex)
}

// GetNodeState implements NodeCommonMatch.
func (c *commonMatchNode) GetNodeState(ip string) (nodeState.VolatileNodeState, error) {
	var nodeStatePriv, err = c.findNode(ip)
	if err != nil {
		return nil, err
	}
	return nodeStatePriv, nil
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
	var numNodeHalf = c.numNode / 2

	nodeStatePriv, err = c.findNode(ip)
	if err != nil {
		return err
	}

	switch indexType {
	case NEXT:
		nodeStatePriv.SetNextIndex(value)
	case MATCH:
		log.Printf("updating match index with value %v, numNodes %v, stable %v\n", value, c.numNode, c.numStable)
		/*
		   TODO: when you want to update the match index of a node:
		   1- check if, before updating, the node has at least the futureCommonIdx as match index
		       1.1- if true: update the match index of node and return
		       1.2- if false:
		           1.2.1- update the match index of node
		           1.2.2- check if the new match index is < the futureCommonIdx:
		               1.2.2.1- if true: return
		               1.2.2.2- if false:
		                   numStable++
		                   for numStable > numNode/2:
		                       notifyChannNewEntry <- futureCommonIdx
		                       futureCommonIdx++
		                       numStable=1
		                       foreach nodestate ns:
		                           if ns.matchIndex > futureCommonIdx:
		                               numStable++
		*/
		matchIdx = nodeStatePriv.GetMatchIndex()
		nodeStatePriv.SetMatchIndex(value)
      //   if matchIdx == c.futureCommonIdx {
		    // log.Panicf("check mathc index, current: %v, common %v\n", matchIdx, c.futureCommonIdx)
      //   }
        log.Printf("check mathc index, current: %v, common %v\n", matchIdx, c.futureCommonIdx)
        //HACK: i don't know why this if else works in this way but it's working
        // at least seems like it, probably it's right but i don't know why
		if matchIdx >= c.futureCommonIdx {
            c.futureCommonIdx++
			return nil
		}

        if value < c.futureCommonIdx {
            return nil
        }
		c.numStable++
		for c.numStable > int(numNodeHalf) {
			c.notifyChannNewEntry <- c.futureCommonIdx
			c.futureCommonIdx++
			c.numStable = 1
			c.allNodeStates.Range(func(key, value any) bool {
				var nodeStateP = value.(nodeState.VolatileNodeState)
				if nodeStateP.GetMatchIndex() > c.futureCommonIdx {
					c.numStable++
				}
				return true
			})
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

	return c.notifyChannNewEntry
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
