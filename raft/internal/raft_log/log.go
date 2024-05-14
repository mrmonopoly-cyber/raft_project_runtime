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

    More_recent_log(last_log_index int64, last_log_term uint64) bool
    clusterConf 
}



func NewLogEntry(baseConf []string) LogEntry {
	var l = new(log)
	l.commitIndex = -1
	l.lastApplied = -1
    l.logSize = 0
	l.entries = make([]*p.LogEntry, 0)
    l.cConf = clusterconf.NewConf(baseConf)

    go l.updateLastApplied()

	return l
}
