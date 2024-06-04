package singleconf

import (
	"errors"
	"raft/internal/node"
	"raft/internal/raft_log"
	"sync"
)

type singleConfImp struct {
	nodeList *sync.Map
	conf     sync.Map
	numNodes uint
	raft_log.LogEntry
}


// GetConfig implements SingleConf.
func (s *singleConfImp) GetConfig() []string {
	var res []string = nil

    s.conf.Range(func(key, value any) bool {
        res = append(res, value.(string))
        return true
    })

    return res
}

// GetNumNodes implements SingleConf.
func (s *singleConfImp) GetNumNodes() uint {
	return s.numNodes
}

// SendMexToNode implements SingleConf.
func (s *singleConfImp) SendMexToNode(ip string, mex []byte) error {
	var v,f = s.nodeList.Load(ip)
    var nodeS node.Node
    if !f{
        return errors.New("node not found")
    }

    nodeS = v.(node.Node)
    return nodeS.Send(mex)
}
