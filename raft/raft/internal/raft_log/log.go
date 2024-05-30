package raft_log

import (
	localfs "raft/internal/localFs"
	clusterconf "raft/internal/raftstate/clusterConf"
	p "raft/pkg/raft-rpcProtobuf-messages/rpcEncoding/out/protobuf"
)


type LogEntry interface {
	GetCommittedEntries() []*p.LogEntry
	GetCommittedEntriesRange(startIndex int) []*p.LogEntry

	GetEntries() []*p.LogEntry
    GetEntriAt(index int64) (*p.LogEntry,error)
    AppendEntries(newEntries []*p.LogEntry)
    DeleteFromEntry(entryIndex uint)
    GetNotificationChanEntry(entry *p.LogEntry) (*chan int,error)

    GetCommitIndex() int64
    IncreaseCommitIndex()
    MinimumCommitIndex(val uint)

	LastLogIndex() int
	LastLogTerm() uint

    //conf info
    cConf
}



func NewLogEntry(fsRootDir string, baseConf []string) LogEntry {
	var l = new(log)
	l.commitIndex = -1
	l.lastApplied = -1
    l.logSize = 0
	l.entries = nil
    l.cConf = clusterconf.NewConf(baseConf)
    l.newEntryToApply = make(chan int)
    l.localFs = localfs.NewFs(fsRootDir)

    go l.updateLastApplied()

	return l
}
