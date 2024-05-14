package raft_log

import (
	"errors"
	l "log"
	localfs "raft/internal/localFs"
	clusterconf "raft/internal/raftstate/clusterConf"
	p "raft/pkg/raft-rpcProtobuf-messages/rpcEncoding/out/protobuf"
	"strings"
	"sync"
)

type cConf interface{
    GetNumberNodesInCurrentConf() int
    IsInConf(ipNode string) bool
    ConfChanged() bool
    GetConfig() []string
}

type log struct {
    lock            sync.RWMutex
    newEntryToApply chan int

    entries     []*p.LogEntry
    logSize     uint
    commitIndex int64
    lastApplied int

    realClusterState
}

type realClusterState struct{
    cConf   clusterconf.Configuration
    localFs localfs.LocalFs
}


// GetCommittedEntries implements LogEntry.
func (this *log) GetCommittedEntries() []*p.LogEntry {
    var committedEntries []*p.LogEntry = make([]*p.LogEntry, this.commitIndex+1)
    for i := range committedEntries {
        committedEntries[i] = this.entries[i]
    }
    return committedEntries
}

// GetEntriAt implements LogEntry.
func (this *log) GetEntriAt(index int64) (*p.LogEntry, error) {
	l.Printf("GetEntiesAt: logSize %v, index: %v", this.logSize, index)
	if (this.logSize == 1 && index == 0) || (index < int64(this.logSize)) {
		return this.entries[index], nil
	}
	return nil, errors.New("invalid index: " + string(rune(index)))
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

// IncreaseCommitIndex implements LogEntry.
func (this *log) IncreaseCommitIndex() {
	this.lock.Lock()
	defer this.lock.Unlock()
	if this.commitIndex < int64(this.logSize)-1 {
		this.commitIndex++
		this.newEntryToApply <- 1
	}
}

// DeleteFromEntry implements LogEntry.
func (this *log) DeleteFromEntry(entryIndex uint) {
	for i := int(entryIndex); i < len(this.entries); i++ {
		this.entries[i] = nil
		this.logSize--
	}
}

//clusterConf

// GetNumberNodesInCurrentConf implements LogEntry.
func (this *log) GetNumberNodesInCurrentConf() int {
	return this.cConf.GetNumberNodesInCurrentConf()
}

// IsInConf implements LogEntry.
func (this *log) IsInConf(nodeIp string) bool {
	return this.cConf.IsInConf(nodeIp)
}

// CommitConfig implements LogEntry.
func (this *log) CommitConfig() {
	this.cConf.CommitConfig()
}

// ConfStatus implements LogEntry.
func (this *log) ConfChanged() bool {
	return this.cConf.ConfChanged()
}

// GetConfig implements LogEntry.
func (this *log) GetConfig() []string {
	return this.cConf.GetConfig()
}

// UpdateConfiguration implements LogEntry.
func (this *log) UpdateConfiguration(confOp clusterconf.CONF_OPE, nodeIps []string) {
	this.cConf.UpdateConfiguration(confOp, nodeIps)
}

//log

func (this *log) GetEntries() []*p.LogEntry {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.entries
}

func (this *log) LastLogIndex() int {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return len(this.entries) - 1
}

func (this *log) More_recent_log(last_log_index int64, last_log_term uint64) bool {
	this.lock.RLock()
	defer this.lock.RUnlock()

	if last_log_index >= this.commitIndex {
		var entries []*p.LogEntry = this.GetEntries()
		if len(entries) <= int(last_log_index) {
			return true
		}
		if last_log_term >= (entries[last_log_index]).Term {
			return true
		}
	}
	return false

}

func (this *log) AppendEntries(newEntries []*p.LogEntry) {
	this.lock.Lock()
	defer this.lock.Unlock()

	var lenEntries = len(this.entries)

	l.Println("Append Entries, before: ", this.entries)
	for _, v := range newEntries {
		if int(this.logSize) < lenEntries {
			this.entries[this.logSize] = v
		} else {
			this.entries = append(this.entries, v)
			lenEntries++
		}
		this.logSize++
	}
	l.Println("Append Entries, after: ", this.entries)
}

func (this *log) GetCommitIndex() int64 {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.commitIndex
}

// utility

func (this *log) updateLastApplied() error {
	for {
		select {
		case <-this.newEntryToApply:
			this.lastApplied++
			var entry *p.LogEntry = this.entries[this.lastApplied]

			l.Printf("updating entry: %v", entry)
			switch entry.OpType {
			case p.Operation_JOIN_CONF_ADD:
				this.applyConf(clusterconf.ADD, entry)
			case p.Operation_JOIN_CONF_DEL:
				this.applyConf(clusterconf.DEL, entry)
			default:
				(*this).localFs.ApplyLogEntry(entry)
			}
		}

	}
}

func (this *log) applyConf(ope clusterconf.CONF_OPE, entry *p.LogEntry) {
	var confUnfiltered string = string(entry.Payload)
	var confFiltered []string = strings.Split(confUnfiltered, " ")
	l.Printf("applying the new conf:%v\t%v\n", confUnfiltered, confFiltered)
	this.cConf.UpdateConfiguration(ope, confFiltered)
}
