package singleconf

import (
	"raft/internal/raft_log"
	"sync"
)

type OP uint8
const(
    ADD OP = iota
    REM OP = iota
)

type SingleConf interface{
    GetConfig() []string
    AutoCommit(status bool)
    raft_log.LogEntry
}

func NewSingleConf(fsRootDir string, conf []string, nodeList *sync.Map) SingleConf{
    return newSingleConfImp(fsRootDir,conf,nodeList)
}

