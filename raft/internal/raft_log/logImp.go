package raft_log

import (
	"log"
	l "log"
	"raft/pkg/raft-rpcProtobuf-messages/rpcEncoding/out/protobuf"
	"reflect"
	"sync"
)

type logEntryImp struct {
	lock sync.RWMutex

	entries     *[]LogInstance
	logSize     uint
	commitIndex int64
}

// GetEntriesRange implements LogEntry.
func (this *logEntryImp) GetEntriesRange(startIndex int) []*protobuf.LogEntry {
    this.lock.RLock()
    defer this.lock.RUnlock()

	var entrs = this.GetEntries()
	return entrs[startIndex:]
}



// AppendEntry implements LogEntry.
func (this *logEntryImp) AppendEntry(newEntrie *LogInstance) {
    this.lock.Lock()
    defer this.lock.Unlock()

    if this.isInLog(newEntrie.Entry,int(this.logSize)-1){
        log.Println("skipping already insert: ",newEntrie.Entry)
        return
    }

	l.Println("adding new entrie to the logEntryImp: ", *newEntrie)
	*this.entries = append(*this.entries, *newEntrie)
	this.logSize++
}

func (this *logEntryImp) GetEntries() []*protobuf.LogEntry {
	this.lock.RLock()
	defer this.lock.RUnlock()

	var lenEntries = len(*this.entries)
	var res []*protobuf.LogEntry = make([]*protobuf.LogEntry, lenEntries)

	for i, v := range *this.entries {
		res[i] = v.Entry
	}

	return res
}

// GetEntriAt implements LogEntry.
func (this *logEntryImp) GetEntriAt(index int64) *LogInstance {
    this.lock.RLock()
    defer this.lock.RUnlock()

	if index < 0 {
		index = 0
	}

    log.Println("GetEntriAt: log and size:",this.entries,this.logSize)
	if index < int64(this.logSize) {
		return &(*this.entries)[index]
	}
	log.Panicln("invald index GetEntrieAt: ", index)
	return nil
}

// DeleteFromEntry implements LogEntry.
func (this *logEntryImp) DeleteFromEntry(entryIndex uint) {
    this.lock.Lock()
    defer this.lock.Unlock()

	for i := int(entryIndex); i < len(*this.entries); i++ {
		(*this.entries)[i] = LogInstance{
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
    this.lock.RLock()
    defer this.lock.RUnlock()

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
	}
}

func (this *logEntryImp) NewLogInstanceBatch(entry []*protobuf.LogEntry, post []func()) []*LogInstance {
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

//utility

func (this *logEntryImp) isInLog(entry *protobuf.LogEntry, index int) bool{
    if index >= int(this.logSize){
        return false
    }

    var savedEntrie = this.GetEntriAt(int64(index))

    return  savedEntrie.Entry.Term == entry.Term &&
            savedEntrie.Entry.OpType == entry.OpType &&
            savedEntrie.Entry.Description == entry.Description &&
            reflect.DeepEqual(savedEntrie.Entry.Payload,entry.Payload)

}

func (this *logEntryImp) getLogState() *logEntryImp{
    return this
}
