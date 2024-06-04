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
    GetNumNodes() uint
    SendMexToNode(ip string, mex []byte) error
    raft_log.LogEntry
}

func NewSingleConf(fsRootDir string, conf []string, nodeList *sync.Map) SingleConf{
    var res = &singleConfImp{
        nodeList: nodeList,
        conf: sync.Map{},
        numNodes: 0,
        LogEntry: raft_log.NewLogEntry(fsRootDir),
    }

    for _, v := range conf {
        res.conf.Store(v,v)
        res.numNodes++
    }
    return res
}

