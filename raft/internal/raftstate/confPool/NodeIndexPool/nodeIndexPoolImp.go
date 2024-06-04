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
func (n *nodeIndexPoolImpl) FetchNodeInfo(info nodestate.INFO, ip string) (int, error) {
    var state,err = n.getNode(ip)
    var res int

    if err != nil{
        return -99,err
    }
    res = state.FetchData(info)
    return res,nil
}

// UpdateNodeInfo implements NodeIndexPool.
func (n *nodeIndexPoolImpl) UpdateNodeInfo(info nodestate.INFO, ip string, val int) {
    var state,err = n.getNode(ip)

    if err != nil{
        return 
    }

    state.UpdateNodeState(info,val)
}

// UpdateStatusList implements NodeIndexPool.
func (n *nodeIndexPoolImpl) UpdateStatusList(op OP, ip string) {
	switch op {
	case ADD:
		var _, f = n.nodeStateList.Load(ip)
		if !f {
			var newState = nodestate.NewNodeState()
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
