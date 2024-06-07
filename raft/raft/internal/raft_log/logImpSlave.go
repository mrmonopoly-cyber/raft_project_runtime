package raft_log

type LogEntrySlaveImp struct {
	*logEntryImp
	commitEntry chan int
	commitIndex int64
}

//NotifyAppendEntryC implements LogEntrySlave
func (this *LogEntrySlaveImp) NotifyAppendEntryC() chan int{
    return this.commitEntry
}

func newLogImpSlave(masterLog LogEntry) *LogEntrySlaveImp {
	var l = &LogEntrySlaveImp{
		logEntryImp: masterLog.getLogState(),
		commitEntry: make(chan int),
        commitIndex: masterLog.GetCommitIndex(),
	}

	return l
}

func (this *LogEntrySlaveImp) GetCommitIndex() int64 {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.commitIndex
}
