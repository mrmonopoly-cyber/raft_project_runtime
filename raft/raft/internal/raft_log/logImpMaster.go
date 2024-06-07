package raft_log

import (
	"log"
	"sync"
)

type LogEntryMasterImp struct {
	logEntryImp
	applyC chan int

	commitIndex int64
}

// ApplyEntryC implements LogEntry.
func (this *LogEntryMasterImp) ApplyEntryC() <-chan int {
	return this.applyC
}

// IncreaseCommitIndex implements LogEntry.
func (this *LogEntryMasterImp) IncreaseCommitIndex() {
	log.Println("increasing commit Index, log imp")
	this.commitIndex++
	if this.applyC != nil {
		this.applyC <- int(this.commitIndex)
	}
	log.Println("increasing commit Index done, log imp")
}

func newLogImpMaster() *LogEntryMasterImp {
    var masterEntries []LogInstance = nil

	var l = &LogEntryMasterImp{
		logEntryImp: logEntryImp{
			logSize:     0,
			entries:     &masterEntries,
			lock:        sync.RWMutex{},
		},
        commitIndex: -1,
		applyC: make(chan int),
	}

	return l
}

// MinimumCommitIndex implements LogEntry.
func (this *LogEntryMasterImp) MinimumCommitIndex(val uint) {
	this.lock.Lock()
	defer this.lock.Unlock()

	if val < this.logSize {
		this.commitIndex = int64(val)
		return
	}
	this.commitIndex = int64(this.logSize) - 1
}

func (this *LogEntryMasterImp) GetCommitIndex() int64 {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.commitIndex
}
