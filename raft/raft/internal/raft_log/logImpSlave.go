package raft_log

type LogEntrySlaveImp struct {
	*logEntryImp
	commitEntry chan int
}

//NotifyAppendEntryC implements LogEntrySlave
func (this *LogEntrySlaveImp) NotifyAppendEntryC() chan int{
    return this.commitEntry
}

func newLogImpSlave(masterLog LogEntry) *LogEntrySlaveImp {
	var l = &LogEntrySlaveImp{
		logEntryImp: masterLog.getLogState(),
		commitEntry: make(chan int),
	}

	return l
}
