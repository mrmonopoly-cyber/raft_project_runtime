package raft_log

import (
	"raft/pkg/raft-rpcProtobuf-messages/rpcEncoding/out/protobuf"
	p "raft/pkg/raft-rpcProtobuf-messages/rpcEncoding/out/protobuf"
)

const SEPARATOR = "K"

type LogInstance struct {
	Entry             *p.LogEntry
    Committed       chan int
    AtCompletion    func() 
}

type LogEntry interface {
	GetCommittedEntries() []LogInstance
	GetCommittedEntriesRange(startIndex int) []LogInstance

	GetEntries() []*protobuf.LogEntry
    GetEntriAt(index int64) (*LogInstance,error)
    AppendEntry(newEntrie *LogInstance)
    DeleteFromEntry(entryIndex uint)

    GetCommitIndex() int64
    MinimumCommitIndex(val uint)

	LastLogIndex() int
	LastLogTerm() uint

    NewLogInstance(entry *p.LogEntry, post func()) *LogInstance
    NewLogInstanceBatch(entry []*p.LogEntry, post []func()) []*LogInstance
}

func NewLogEntry(oldEntries []*protobuf.LogEntry) LogEntry {
    return NewLogImp(oldEntries)
}

