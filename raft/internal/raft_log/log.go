package raft_log

import (
	clusterconf "raft/internal/raftstate/clusterConf"
	p "raft/pkg/raft-rpcProtobuf-messages/rpcEncoding/out/protobuf"
)

type LogEntry interface {
	GetEntries() []*p.LogEntry
	GetCommitIndex() int64
	More_recent_log(last_log_index int64, last_log_term uint64) bool
	SetCommitIndex(val int64)
	AppendEntries(newEntries []*p.LogEntry)
	LastLogIndex() int
	UpdateLastApplied() error
    clusterconf.Configuration
}



func NewLogEntry(baseConf []string) LogEntry {
	var l = new(log)
	l.commitIndex = -1
	l.lastApplied = -1
	l.entries = make([]*p.LogEntry, 0)
    l.cConf = clusterconf.NewConf(baseConf)

	return l
}
