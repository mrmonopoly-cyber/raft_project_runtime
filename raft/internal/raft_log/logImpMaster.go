package raft_log

import (
	"log"
	"sync"
)

type LogEntryMasterImp struct {
	logEntryImp
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
			commitIndex: -1,
			logSize:     0,
			entries:     &masterEntries,
			lock:        sync.RWMutex{},
            applyC: make(chan int),
		},
	}

	return l
}
