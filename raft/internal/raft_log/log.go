package raft_log

import (
	localfs "raft/internal/localFs"
	clusterconf "raft/internal/raftstate/clusterConf"
	p "raft/pkg/raft-rpcProtobuf-messages/rpcEncoding/out/protobuf"
)

type LogInstance struct {
	Entry             *p.LogEntry
	NotifyApplication chan int
}

type LogEntry interface {
	GetCommittedEntries() []LogInstance
	GetCommittedEntriesRange(startIndex int) []LogInstance

	GetEntries() []LogInstance
    GetEntriAt(index int64) (*LogInstance,error)
    AppendEntries(newEntries []*LogInstance)
    DeleteFromEntry(entryIndex uint)

    GetCommitIndex() int64
    IncreaseCommitIndex()
    MinimumCommitIndex(val uint)

	LastLogIndex() int
	LastLogTerm() uint

    NewLogInstanceBatch(entry []*p.LogEntry) []LogInstance
    NewLogInstance(entry *p.LogEntry) *LogInstance

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

func (this *log) NewLogInstance(entry *p.LogEntry) *LogInstance{
    return &LogInstance{
        Entry: entry,
        NotifyApplication: make(chan int),
    }
}

func (this *log) NewLogInstanceBatch(entry []*p.LogEntry) []LogInstance{
    var res []LogInstance = make([]LogInstance, len(entry))

    for i, v := range entry {
        res[i] = LogInstance{
            Entry: v,
            NotifyApplication: make(chan int),
        }
    }

    return res
}
