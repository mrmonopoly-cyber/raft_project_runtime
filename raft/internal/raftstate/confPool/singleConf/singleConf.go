package singleconf

import (
	"raft/internal/raft_log"
	clustermetadata "raft/internal/raftstate/clusterMetadata"
	nodeIndexPool "raft/internal/raftstate/confPool/NodeIndexPool"
	confmetadata "raft/internal/raftstate/confPool/singleConf/confMetadata"
	"sync"
)

type OP uint8
const(
    ADD OP = iota
    REM OP = iota
)

type SingleConf interface{
    confmetadata.ConfMetadata
    SendHearthbit()
    CommiEntryC() <- chan int
    CloseCommitEntryC()
    raft_log.LogEntrySlave
}

func NewSingleConf( conf map[string]string,  
                    masterLog raft_log.LogEntry,
                    nodeList *sync.Map,
                    commonStatePool nodeIndexPool.NodeIndexPool,
                    commonMetadata clustermetadata.ClusterMetadata) SingleConf{
    return newSingleConfImp(conf,masterLog, nodeList,commonStatePool,commonMetadata)
}

