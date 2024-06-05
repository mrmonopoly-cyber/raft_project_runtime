package raft_log

import (
	"errors"
	l "log"
	"raft/pkg/raft-rpcProtobuf-messages/rpcEncoding/out/protobuf"
	"sync"
)

type logEntryImp struct {
	lock            sync.RWMutex
	newEntryToApply chan int

	entries     []LogInstance
	logSize     uint
	commitIndex int64
	lastApplied int
}

// AppendEntry implements LogEntry.
func (this *logEntryImp) AppendEntry(newEntrie *LogInstance) {
    l.Println("adding new entrie to the logEntryImp: ",*newEntrie)
    this.entries = append(this.entries, *newEntrie)
    this.logSize++

    go func ()  {
       <- newEntrie.Committed 
       this.commitIndex++
    }()
}

// logEntryImp
// GetCommittedEntries implements LogEntry.
func (this *logEntryImp) GetCommittedEntries() []LogInstance {
	return this.getEntries(0)
}

// GetCommittedEntriesRange implements LogEntry.
func (this *logEntryImp) GetCommittedEntriesRange(startIndex int) []LogInstance {
	return this.getEntries(startIndex)
}

func (this *logEntryImp) GetEntries() []*protobuf.LogEntry{
	this.lock.RLock()
	defer this.lock.RUnlock()

	var lenEntries = len(this.entries)
	var res []*protobuf.LogEntry = make([]*protobuf.LogEntry, lenEntries)

	for i, v := range this.entries {
		res[i] = v.Entry
	}

	return res
}

// GetEntriAt implements LogEntry.
func (this *logEntryImp) GetEntriAt(index int64) (*LogInstance, error) {
	if (index < int64(this.logSize)) {
		return &this.entries[index], nil
	}
	return nil, errors.New("invalid index: " + string(rune(index)))
}

// DeleteFromEntry implements LogEntry.
func (this *logEntryImp) DeleteFromEntry(entryIndex uint) {
	for i := int(entryIndex); i < len(this.entries); i++ {
		this.entries[i] = LogInstance{
			Entry:        nil,
			AtCompletion: func() {},
		}
		this.logSize--
	}
}

func (this *logEntryImp) GetCommitIndex() int64 {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.commitIndex
}

// MinimumCommitIndex implements LogEntry.
func (this *logEntryImp) MinimumCommitIndex(val uint) {
	this.lock.Lock()
	defer this.lock.Unlock()

	if val < this.logSize {
		this.commitIndex = int64(val)
		return
	}
	this.commitIndex = int64(this.logSize) - 1
	if this.commitIndex > int64(this.lastApplied) {
		this.newEntryToApply <- 1
	}
}

func (this *logEntryImp) LastLogIndex() int {
	this.lock.RLock()
	defer this.lock.RUnlock()

	var committedEntr = this.GetCommittedEntries()
	return len(committedEntr) - 1
}

// LastLogTerm implements LogEntry.
func (this *logEntryImp) LastLogTerm() uint {
	var committedEntr = this.GetCommittedEntries()
	var lasLogIdx = this.LastLogIndex()

	if lasLogIdx >= 0 {
		return uint(committedEntr[lasLogIdx].Entry.Term)
	}
	return 0

}

func (this *logEntryImp) NewLogInstance(entry *protobuf.LogEntry, post func()) *LogInstance{
    return &LogInstance{
        Entry: entry,
        AtCompletion: post,
        Committed: make(chan int),
    }
}

func (this *logEntryImp) NewLogInstanceBatch(entry []*protobuf.LogEntry, post []func()) []*LogInstance{
    var res []*LogInstance = make([]*LogInstance, len(entry))

    for i, v := range entry {
        res[i] = &LogInstance{
            Entry: v,
            Committed: make(chan int),
        }
        if post != nil && i < len(post) {
            res[i].AtCompletion = post[i]
        }
    }

    return res
}

func NewLogImp(oldEntries []*protobuf.LogEntry) *logEntryImp{
    var oldEntrLen = len(oldEntries)
    var oldInstance []LogInstance = make([]LogInstance, oldEntrLen)

    for i := 0; i < oldEntrLen; i++ {
        oldInstance[i].Entry = oldEntries[i]
    }

	var l = &logEntryImp{
        commitIndex: -1,
        lastApplied: -1,
        logSize: 0,
        entries: oldInstance,
        newEntryToApply: make(chan int),
        lock: sync.RWMutex{},
    }

	return l
}

// utility
func (this *logEntryImp) getEntries(startIndex int) []LogInstance {
	var committedEntries []LogInstance = nil

	if startIndex == -1 {
		startIndex = 0
	}

	for i := int(startIndex); i <= int(this.commitIndex); i++ {
		committedEntries = append(committedEntries, this.entries[i])
	}
	return committedEntries
}

