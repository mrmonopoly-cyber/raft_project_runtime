package raft_log

import "sync"

type LogEntrySlaveImp struct {
	logEntryImp
	commitEntry chan int
}

//NotifyAppendEntryC implements LogEntrySlave
func (this *LogEntrySlaveImp) NotifyAppendEntryC() chan int{
    return this.commitEntry
}

func newLogImpSlave(masterLog LogEntry) *LogEntrySlaveImp {
	var l = &LogEntrySlaveImp{
		logEntryImp: logEntryImp{
			commitIndex: masterLog.GetCommitIndex(),
			logSize:     masterLog.getLogSize(),
			entries:     masterLog.getEntriesRaw(),
			lock:        sync.RWMutex{},
		},
		commitEntry: make(chan int),
	}

	return l
}
