package raft_log

import (
	clusterconf "raft/internal/raftstate/clusterConf"
	p "raft/pkg/raft-rpcProtobuf-messages/rpcEncoding/out/protobuf"
)

type LogEntry interface {
	GetCommittedEntries() []*p.LogEntry
	GetEntries() []*p.LogEntry
    GetEntriAt(index int64) (*p.LogEntry,error)
    AppendEntries(newEntries []*p.LogEntry)
    DeleteFromEntry(entryIndex uint)

    GetCommitIndex() int64
    IncreaseCommitIndex()
    MinimumCommitIndex(val uint)

	LastLogIndex() int
	LastLogTerm() uint

    cConf 
}



func NewLogEntry(baseConf []string) LogEntry {
	var l = new(log)
	l.commitIndex = -1
	l.lastApplied = -1
    l.logSize = 0
	l.entries = make([]*p.LogEntry, 0)
    l.cConf = clusterconf.NewConf(baseConf)
    l.newEntryToApply = make(chan int)

    go l.updateLastApplied()

	return l
}
