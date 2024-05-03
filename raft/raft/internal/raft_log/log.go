package raft_log

import (
	clusterconf "raft/internal/raftstate/clusterConf"
    "raft/internal/localFs"
	p "raft/pkg/raft-rpcProtobuf-messages/rpcEncoding/out/protobuf"
	"sync"
)

type LogEntry interface {
	GetEntries() []*p.LogEntry
	GetCommitIndex() int64
	More_recent_log(last_log_index int64, last_log_term uint64) bool
	AppendEntries(newEntries []*p.LogEntry)
	LastLogIndex() int

    GetClusterConfig() []string
    IsInConf(ip string) bool
    ConfStatus() bool
}



func NewLogEntry(baseConf []string, rootLFS string) LogEntry {
	var l = new(log)
	l.commitIndex = 0
	l.lastApplied = 0
    l.commitEntries = make([]*p.LogEntry, 0)
    l.cConf = clusterconf.NewConf(baseConf)
    l.localFs = localfs.NewFs(rootLFS)
    l.wg = sync.WaitGroup{}

    go func (){
        l.wg.Add(1)
        defer l.wg.Done()
        l.applyLogEntry()
    }()

	return l
}
