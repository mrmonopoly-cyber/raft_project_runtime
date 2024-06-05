package singleconf

import (
	"log"
	"raft/internal/node"
	"raft/internal/raft_log"
	clustermetadata "raft/internal/raftstate/clusterMetadata"
	nodeIndexPool "raft/internal/raftstate/confPool/NodeIndexPool"
	"sync"
)

type singleConfImp struct {
	nodeList *sync.Map
	conf     sync.Map
	numNodes uint
    autoCommit *bool
	raft_log.LogEntry
    nodeIndexPool.NodeIndexPool
    clustermetadata.ClusterMetadata
}

func (s *singleConfImp) AppendEntry(entry *raft_log.LogInstance) {
    s.LogEntry.AppendEntry(entry)
    if *s.autoCommit{ //INFO: FOLLOWER
        s.LogEntry.IncreaseCommitIndex()
        return
    }

    //INFO:LEADER
    //Propagate to all nodes in this conf
    s.conf.Range(func(key, value any) bool {
        var v,f = s.nodeList.Load(key)
        var fNode node.Node


        if !f{
            log.Println("node not yet connected or crashes, skipping send")
            return false
        }
        fNode = v.(node.Node)
        fNode.GetIp() //INFO: just to remove the compiler error

        //TODO: generate an AppendEntry for this specific node and send it
        //The AppendEntry will contain all the missed entry + the new one

        return true
    })
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

func newSingleConfImp(  fsRootDir string, 
                        conf []string, 
                        nodeList *sync.Map,
                        autoCommit *bool,
                        commonStatePool nodeIndexPool.NodeIndexPool,
                        commonMetadata clustermetadata.ClusterMetadata) *singleConfImp{
    var res = &singleConfImp{
        nodeList: nodeList,
        conf: sync.Map{},
        numNodes: 0,
        LogEntry: raft_log.NewLogEntry(fsRootDir),
        autoCommit: autoCommit,
        NodeIndexPool: commonStatePool,
        ClusterMetadata: commonMetadata,
    }

    for _, v := range conf {
        res.conf.Store(v,v)
        res.numNodes++
    }
    return res
    
}
