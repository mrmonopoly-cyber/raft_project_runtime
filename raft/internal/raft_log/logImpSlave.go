package raft_log

type logEntrySlaveImp struct {
	*logEntryImp
	commitEntry chan int
}

//NotifyAppendEntryC implements LogEntrySlave
func (this *logEntrySlaveImp) NotifyAppendEntryC() chan int{
    return this.commitEntry
}

func (this *logEntrySlaveImp) CloseNotifyAppendEntryC(){
    close(this.commitEntry)
}

func newLogImpSlave(masterLog LogEntry) *logEntrySlaveImp {
	var l = &logEntrySlaveImp{
		logEntryImp: masterLog.getLogState(),
		commitEntry: make(chan int),
	}

	return l
}
