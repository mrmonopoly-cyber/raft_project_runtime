package raft_log

import (
	localfs "raft/internal/localFs"
	clusterconf "raft/internal/raftstate/clusterConf"
	p "raft/pkg/raft-rpcProtobuf-messages/rpcEncoding/out/protobuf"
)

type LogInstance struct {
	Entry             *p.LogEntry
    AtCompletion    func() 
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

    NewLogInstance(entry *p.LogEntry, post func()) *LogInstance
    NewLogInstanceBatch(entry []*p.LogEntry, post []func()) []*LogInstance

    //conf info
    cConf
}



func NewLogEntry(fsRootDir string) LogEntry {
	var l = new(log)
	l.commitIndex = -1
	l.lastApplied = -1
    l.logSize = 0
	l.entries = nil
    l.Configuration = clusterconf.NewConf()
    l.newEntryToApply = make(chan int)
    l.LocalFs = localfs.NewFs(fsRootDir)

    go l.updateLastApplied()

	return l
}

func (this *log) NewLogInstance(entry *p.LogEntry, post func()) *LogInstance{
    return &LogInstance{
        Entry: entry,
        AtCompletion: post,
    }
}

func (this *log) NewLogInstanceBatch(entry []*p.LogEntry, post []func()) []*LogInstance{
    var res []*LogInstance = make([]*LogInstance, len(entry))

    for i, v := range entry {
        res[i] = &LogInstance{
            Entry: v,
        }
        if post != nil && i < len(post) {
            res[i].AtCompletion = post[i]
        }
    }

    return res
}
