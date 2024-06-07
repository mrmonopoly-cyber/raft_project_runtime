package singleconf

import (
	"raft/internal/raft_log"
	clustermetadata "raft/internal/raftstate/clusterMetadata"
	nodeIndexPool "raft/internal/raftstate/confPool/NodeIndexPool"
	"sync"
)

type OP uint8
const(
    ADD OP = iota
    REM OP = iota
)

type SingleConf interface{
    GetConfig() map[string]string
    SendHearthbit()
    CommiEntryC() <- chan int
    raft_log.LogEntrySlave
}

func NewSingleConf( conf []string,  
                    masterLog raft_log.LogEntry,
                    nodeList *sync.Map,
                    commonStatePool nodeIndexPool.NodeIndexPool,
                    commonMetadata clustermetadata.ClusterMetadata) SingleConf{
    return newSingleConfImp(conf,masterLog, nodeList,commonStatePool,commonMetadata)
}

