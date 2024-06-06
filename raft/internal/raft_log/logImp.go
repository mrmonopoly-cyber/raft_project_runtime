package raft_log

import (
	"errors"
	l "log"
	"raft/pkg/raft-rpcProtobuf-messages/rpcEncoding/out/protobuf"
	"sync"
)

type logEntryImp struct {
	lock sync.RWMutex

	entries     []LogInstance
	logSize     uint
	commitIndex int64
	applyC      chan int
}

// GetEntriesRange implements LogEntry.
func (this *logEntryImp) GetEntriesRange(startIndex int) []*protobuf.LogEntry {
    var entrs = this.GetEntries()
    return entrs[startIndex:]
}

// ApplyEntryC implements LogEntry.
func (this *logEntryImp) ApplyEntryC() <-chan int {
	return this.applyC
}

// IncreaseCommitIndex implements LogEntry.
func (this *logEntryImp) IncreaseCommitIndex() {
	this.commitIndex++
	this.applyC <- int(this.commitIndex)
}

// AppendEntry implements LogEntry.
func (this *logEntryImp) AppendEntry(newEntrie *LogInstance) {
	l.Println("adding new entrie to the logEntryImp: ", *newEntrie)
	this.entries = append(this.entries, *newEntrie)
	this.logSize++
}

func (this *logEntryImp) GetEntries() []*protobuf.LogEntry {
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
	if index < int64(this.logSize) {
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
}

func (this *logEntryImp) LastLogIndex() int {
	this.lock.RLock()
	defer this.lock.RUnlock()

	return int(this.logSize) - 1
}

// LastLogTerm implements LogEntry.
func (this *logEntryImp) LastLogTerm() uint {
	var committedEntr = this.GetEntries()
	var lasLogIdx = this.LastLogIndex()

	if lasLogIdx >= 0 {
		return uint(committedEntr[lasLogIdx].Term)
	}
	return 0

}

func (this *logEntryImp) NewLogInstance(entry *protobuf.LogEntry, post func()) *LogInstance {
	return &LogInstance{
		Entry:        entry,
		AtCompletion: post,
		Committed:    make(chan int),
	}
}

func (this *logEntryImp) NewLogInstanceBatch(entry []*protobuf.LogEntry, post []func()) []*LogInstance {
	var res []*LogInstance = make([]*LogInstance, len(entry))

	for i, v := range entry {
		res[i] = &LogInstance{
			Entry:     v,
			Committed: make(chan int),
		}
		if post != nil && i < len(post) {
			res[i].AtCompletion = post[i]
		}
	}

	return res
}


func newLogImp(oldEntries []*protobuf.LogEntry) *logEntryImp {
	var oldEntrLen = len(oldEntries)
	var oldInstance []LogInstance = make([]LogInstance, oldEntrLen)

	for i := 0; i < oldEntrLen; i++ {
		oldInstance[i].Entry = oldEntries[i]
	}

	var l = &logEntryImp{
		commitIndex: -1,
		logSize:     0,
		entries:     oldInstance,
		lock:        sync.RWMutex{},
		applyC:      make(chan int),
	}

	return l
}
