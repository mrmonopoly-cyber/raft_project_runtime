package nodeIndexPool

import (
	"errors"
	"log"
	"raft/internal/raftstate/confPool/NodeIndexPool/nodeState"
	"sync"
)

type nodeIndexPoolImpl struct {
	nodeStateList sync.Map
}

// FetchNodeInfo implements NodeIndexPool.
func (n *nodeIndexPoolImpl) FetchNodeInfo(ip string) (nodestate.NodeState, error) {
    var state,err = n.getNode(ip)

    if err != nil{
        var newState = nodestate.NewNodeState(ip)
		n.nodeStateList.Store(ip, newState)
        return newState,err
    }
    return state,nil
}

// UpdateStatusList implements NodeIndexPool.
func (n *nodeIndexPoolImpl) UpdateStatusList(op OP, ip string) {
	switch op {
	case ADD:
		var _, f = n.nodeStateList.Load(ip)
		if !f {
			var newState = nodestate.NewNodeState(ip)
			n.nodeStateList.Store(ip, newState)
		}
        return
	case REM:
		n.nodeStateList.Delete(ip)
        return
	}
    log.Panicln("unmanaged case: ", op)
}

func newNodeIndexPoolImpl() *nodeIndexPoolImpl {
	return &nodeIndexPoolImpl{
		nodeStateList: sync.Map{},
	}

}

//utility
func (n *nodeIndexPoolImpl) getNode(ip string) (nodestate.NodeState,error){
    var v,f= n.nodeStateList.Load(ip)

    if !f{
        return nil,errors.New("state for node " + ip + " not found")
    }

    return v.(nodestate.NodeState),nil
}
