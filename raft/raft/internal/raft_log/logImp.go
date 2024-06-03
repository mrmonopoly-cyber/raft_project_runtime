package raft_log

import (
	"errors"
	l "log"
	localfs "raft/internal/localFs"
	clusterconf "raft/internal/raftstate/clusterConf"
	"raft/pkg/raft-rpcProtobuf-messages/rpcEncoding/out/protobuf"
	p "raft/pkg/raft-rpcProtobuf-messages/rpcEncoding/out/protobuf"
	"strings"
	"sync"
)

type cConf interface {
	GetNumberNodesInCurrentConf() int
	IsInConf(ipNode string) bool
	NotifyChangeInConf() <-chan int
	GetConfig() []string
	GetRootFs() string
}

type log struct {
	lock            sync.RWMutex
	newEntryToApply chan int

	entries     []LogInstance
	logSize     uint
	commitIndex int64
	lastApplied int

	realClusterState
}

// GetConfig implements LogEntry.
// Subtle: this method shadows the method (realClusterState).GetConfig of log.realClusterState.
func (this *log) GetConfig() []string {
    return this.realClusterState.GetConfig()
}

// GetNumberNodesInCurrentConf implements LogEntry.
// Subtle: this method shadows the method (realClusterState).GetNumberNodesInCurrentConf of log.realClusterState.
func (this *log) GetNumberNodesInCurrentConf() int {
    return this.realClusterState.GetNumberNodesInCurrentConf()
}

// GetRootFs implements LogEntry.
func (this *log) GetRootFs() string {
    return this.realClusterState.GetRootDir()
}

// IsInConf implements LogEntry.
// Subtle: this method shadows the method (realClusterState).IsInConf of log.realClusterState.
func (this *log) IsInConf(ipNode string) bool {
    return this.realClusterState.IsInConf(ipNode)
}

// ResetConf implements LogEntry.
func (this *log) ResetLog() {
	this.realClusterState.Configuration = clusterconf.NewConf()
}

// NotifyChangeInConfChan implements LogEntry.
func (this *log) NotifyChangeInConf() <-chan int {
	return this.realClusterState.Configuration.NotifyChangeInConf()
}

type realClusterState struct {
	clusterconf.Configuration
	localfs.LocalFs
}

// log
// GetCommittedEntries implements LogEntry.
func (this *log) GetCommittedEntries() []LogInstance {
	return this.getEntries(0)
}

// GetCommittedEntriesRange implements LogEntry.
func (this *log) GetCommittedEntriesRange(startIndex int) []LogInstance {
	return this.getEntries(startIndex)
}

func (this *log) GetEntries() []LogInstance {
	this.lock.RLock()
	defer this.lock.RUnlock()

	var lenEntries = len(this.entries)
	var res []LogInstance = make([]LogInstance, lenEntries)

	for i, v := range this.entries {
		res[i] = v
	}

	return res
}

// GetEntriAt implements LogEntry.
func (this *log) GetEntriAt(index int64) (*LogInstance, error) {
	if (this.logSize == 1 && index == 0) || (index < int64(this.logSize)) {
		return &this.entries[index], nil
	}
	return nil, errors.New("invalid index: " + string(rune(index)))
}

func (this *log) AppendEntries(newEntries []*LogInstance) {
	this.lock.Lock()
	defer this.lock.Unlock()

	var lenEntries = len(this.entries)

	for _, v := range newEntries {
		l.Println("adding new logEntry: ", v)
		if int(this.logSize) < lenEntries {
			this.entries[this.logSize] = *v
		} else {
			this.entries = append(this.entries, *v)
			lenEntries++
		}
		this.logSize++
	}
}

// DeleteFromEntry implements LogEntry.
func (this *log) DeleteFromEntry(entryIndex uint) {
	for i := int(entryIndex); i < len(this.entries); i++ {
		this.entries[i] = LogInstance{
			Entry:        nil,
			AtCompletion: func() {},
		}
		this.logSize--
	}
}

func (this *log) GetCommitIndex() int64 {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.commitIndex
}

// IncreaseCommitIndex implements LogEntry.
func (this *log) IncreaseCommitIndex() {
	this.lock.Lock()
	defer this.lock.Unlock()
	if this.commitIndex < int64(this.logSize)-1 {
		this.commitIndex++
		this.newEntryToApply <- 1
	}
}

// MinimumCommitIndex implements LogEntry.
func (this *log) MinimumCommitIndex(val uint) {
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

func (this *log) LastLogIndex() int {
	this.lock.RLock()
	defer this.lock.RUnlock()

	var committedEntr = this.GetCommittedEntries()
	return len(committedEntr) - 1
}

// LastLogTerm implements LogEntry.
func (this *log) LastLogTerm() uint {
	var committedEntr = this.GetCommittedEntries()
	var lasLogIdx = this.LastLogIndex()

	if lasLogIdx >= 0 {
		return uint(committedEntr[lasLogIdx].Entry.Term)
	}
	return 0

}

// utility

func (this *log) updateLastApplied() error {
	for {
		select {
		case <-this.newEntryToApply:
			this.lastApplied++
			var entry *LogInstance = &this.entries[this.lastApplied]

			l.Printf("updating entry: %v", entry)
			switch entry.Entry.OpType {
			case p.Operation_JOIN_CONF_ADD, p.Operation_JOIN_CONF_DEL,
				p.Operation_COMMIT_CONFIG_REM, p.Operation_COMMIT_CONFIG_ADD:
				this.applyConf(entry.Entry.OpType, entry)
			default:
				(*this).ApplyLogEntry(entry.Entry)
			}

			go entry.AtCompletion()
		}

	}
}

func (this *log) applyConf(ope protobuf.Operation, entry *LogInstance) {
	var confUnfiltered string = string(entry.Entry.Payload)
	var confFiltered []string = strings.Split(confUnfiltered, "K")
	l.Printf("applying the new conf:%v\t%v\n", confUnfiltered, confFiltered)
	this.UpdateConfiguration(ope, confFiltered)
	//HACK: if you are follower this goroutine remain stuck forever
	//creating a zombie process
}

func (this *log) getEntries(startIndex int) []LogInstance {
	var committedEntries []LogInstance = nil

	if startIndex == -1 {
		startIndex = 0
	}

	for i := int(startIndex); i <= int(this.commitIndex); i++ {
		committedEntries = append(committedEntries, this.entries[i])
	}
	return committedEntries
}
