package singleconf

import (
	"raft/internal/raft_log"
	nodeIndexPool "raft/internal/raftstate/confPool/NodeIndexPool"
	"sync"
)

type OP uint8
const(
    ADD OP = iota
    REM OP = iota
)

type SingleConf interface{
    GetConfig() []string
    raft_log.LogEntry
}

func NewSingleConf( fsRootDir string, 
                    conf []string,  
                    nodeList *sync.Map,
                    autoCommit *bool,
                    commonStatePool nodeIndexPool.NodeIndexPool) SingleConf{
    return newSingleConfImp(fsRootDir,conf,nodeList,autoCommit,commonStatePool)
}

