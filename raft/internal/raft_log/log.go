package raft_log

import (
	p "raft/pkg/raft-rpcProtobuf-messages/rpcEncoding/out/protobuf"
)

type LogEntry interface {
	GetEntries() []*p.LogEntry
	GetCommitIndex() int64
	More_recent_log(last_log_index int64, last_log_term uint64) bool
	SetCommitIndex(val int64)
	AppendEntries(newEntries []*p.LogEntry, index int)
	LastLogIndex() int
	UpdateLastApplied() (bool, *p.LogEntry)
	InitState()
}



func NewLogEntry() LogEntry {
	var l = new(log)
	l.commitIndex = 0
	l.lastApplied = 0
	l.entries = make([]p.LogEntry, 0)

	return l
}


