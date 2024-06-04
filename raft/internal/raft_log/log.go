package raft_log

import (
	localfs "raft/internal/localFs"
	p "raft/pkg/raft-rpcProtobuf-messages/rpcEncoding/out/protobuf"
	"sync"
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

	GetEntries() []LogInstance
    GetEntriAt(index int64) (*LogInstance,error)
    AppendEntry(newEntrie *LogInstance)
    DeleteFromEntry(entryIndex uint)

    GetCommitIndex() int64
    IncreaseCommitIndex()
    MinimumCommitIndex(val uint)

	LastLogIndex() int
	LastLogTerm() uint

    NewLogInstance(entry *p.LogEntry, post func()) *LogInstance
    NewLogInstanceBatch(entry []*p.LogEntry, post []func()) []*LogInstance

}



func NewLogEntry(fsRootDir string) LogEntry {
	var l = &log{
        commitIndex: -1,
        lastApplied: -1,
        logSize: 0,
        entries: nil,
        newEntryToApply: make(chan int),
        LocalFs: localfs.NewFs(fsRootDir),
        lock: sync.RWMutex{},
    }

    go l.updateLastApplied()

	return l
}

func (this *log) NewLogInstance(entry *p.LogEntry, post func()) *LogInstance{
    return &LogInstance{
        Entry: entry,
        AtCompletion: post,
        Committed: make(chan int),
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
