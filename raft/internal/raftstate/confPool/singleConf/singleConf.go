package singleconf

import (
	"raft/internal/raft_log"
	clustermetadata "raft/internal/raftstate/clusterMetadata"
	nodeIndexPool "raft/internal/raftstate/confPool/NodeIndexPool"
	"raft/pkg/raft-rpcProtobuf-messages/rpcEncoding/out/protobuf"
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
                    oldEntries []*protobuf.LogEntry,
                    nodeList *sync.Map,
                    commonStatePool nodeIndexPool.NodeIndexPool,
                    commonMetadata clustermetadata.ClusterMetadata) SingleConf{
    return newSingleConfImp(fsRootDir,conf,oldEntries, nodeList,commonStatePool,commonMetadata)
}

