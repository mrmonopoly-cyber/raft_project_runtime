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
    var numNodeHalf = c.numNode/2

	nodeStatePriv, err = c.findNode(ip)
	if err != nil {
		return err
	}

	switch indexType {
	case NEXT:
		nodeStatePriv.SetNextIndex(value)
	case MATCH:
        /*
        TODO: when you want to update the match index of a node:
        1- check if, before updating, the node has at least the commonIdx as match index
            1.1- if true: update the match index of node and return
            1.2- if false:
                1.2.1- update the match index of node
                1.2.2- check if the new match index is < the commonIdx:
                    1.2.2.1- if true: return
                    1.2.2.2- if false:
                        numStable++
                        for numStable > numNode/2:
                            notifyChann <- commonIdx
                            commonIdx++
                            numStable=0
                            foreach nodestate ns:
                                if ns.matchIndex > commonIdx:
                                    numStable++
        */
        matchIdx = nodeStatePriv.GetMatchIndex()
        nodeStatePriv.SetMatchIndex(value)
        if matchIdx >= c.commonIdx || value < c.commonIdx{
            return nil
        }
        c.numStable++
        for c.numStable > int(numNodeHalf){
            c.notifyChann <- c.commonIdx
            c.commonIdx++
            c.numStable=0
            c.allNodeStates.Range(func(key, value any) bool {
                var nodeStateP = value.(nodeState.VolatileNodeState)
                if nodeStateP.GetMatchIndex() > c.commonIdx{
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
