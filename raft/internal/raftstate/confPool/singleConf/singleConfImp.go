package singleconf

import (
	"raft/internal/raft_log"
	nodematchidx "raft/internal/raftstate/nodeMatchIdx"
	"sync"
)

type singleConfImp struct {
	nodeList *sync.Map
	conf     sync.Map
	numNodes uint
    autoCommit bool
	raft_log.LogEntry
    nodematchidx.NodeCommonMatch
}

func (s *singleConfImp) AppendEntry(entry *raft_log.LogInstance) {
    s.LogEntry.AppendEntry(entry)
    if s.autoCommit{ //INFO: FOLLOWER
        s.LogEntry.IncreaseCommitIndex()
        return
    }
    //INFO:LEADER
}

// AutoCommit implements SingleConf.
func (s *singleConfImp) AutoCommit(status bool) {
	s.autoCommit = true
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

func newSingleConfImp(fsRootDir string, conf []string, nodeList *sync.Map) *singleConfImp{
    var res = &singleConfImp{
        nodeList: nodeList,
        conf: sync.Map{},
        numNodes: 0,
        LogEntry: raft_log.NewLogEntry(fsRootDir),
        NodeCommonMatch: nodematchidx.NewNodeCommonMatch(),
    }

    for _, v := range conf {
        res.conf.Store(v,v)
        res.numNodes++
    }
    return res
    
}
