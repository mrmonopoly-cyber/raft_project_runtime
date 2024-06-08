package raft_log

import (
	"log"
	"sync"
)

type logEntryMasterImp struct {
	logEntryImp
	applyC chan int
}

// ApplyEntryC implements LogEntry.
func (this *logEntryMasterImp) ApplyEntryC() <-chan int {
	return this.applyC
}

// IncreaseCommitIndex implements LogEntry.
func (this *logEntryMasterImp) IncreaseCommitIndex() {
	log.Println("increasing commit Index, log imp")
	this.commitIndex++
	if this.applyC != nil {
		this.applyC <- int(this.commitIndex)
	}
	log.Println("increasing commit Index done, log imp")
}

func newLogImpMaster() *logEntryMasterImp {
    var masterEntries []LogInstance = nil

	var l = &logEntryMasterImp{
		logEntryImp: logEntryImp{
			commitIndex: -1,
			logSize:     0,
			entries:     &masterEntries,
			lock:        sync.RWMutex{},
		},
		applyC: make(chan int),
	}

	return l
}
