package raft_log

import (
	"raft/internal/utiliy"
	"raft/pkg/raft-rpcProtobuf-messages/rpcEncoding/out/protobuf"
	p "raft/pkg/raft-rpcProtobuf-messages/rpcEncoding/out/protobuf"
)

const SEPARATOR = "K"

type LogInstance struct {
	Entry             *p.LogEntry
    ReturnValue       chan utiliy.Pair[[]byte,error]
}

type logEntryWrite interface {
    AppendEntry(newEntrie []*LogInstance, prevLogIndex int) uint
    DeleteFromEntry(entryIndex uint)

    MinimumCommitIndex(val uint)
    IncreaseCommitIndex()

    ApplyEntryC() chan int
}

type LogEntryRead interface {
	GetEntries() []*protobuf.LogEntry
    GetEntriAt(index int64) *LogInstance
    GetEntriesRange(startIndex int) []*protobuf.LogEntry

    GetCommitIndex() int64
    GetLogSize() uint

	LastLogIndex() int
	LastLogTerm() uint

    NewLogInstance(entry *p.LogEntry) *LogInstance
    NewLogInstanceBatch(entry []*p.LogEntry) []*LogInstance

}

type LogEntrySlave interface{
    LogEntryRead
    NotifyAppendEntryC() chan int
}

type LogEntry interface {
    LogEntryRead
    logEntryWrite
    getLogState() *logEntryImp
}

func NewLogEntry() LogEntry {
    return newLogImpMaster()
}

func NewLogEntrySlave(masterLog LogEntry) LogEntrySlave{
    return newLogImpSlave(masterLog)
}

