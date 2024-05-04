package raft_log

import (
	"fmt"
	l "log"
	localfs "raft/internal/localFs"
	clusterconf "raft/internal/raftstate/clusterConf"
	p "raft/pkg/raft-rpcProtobuf-messages/rpcEncoding/out/protobuf"
	"strings"
)

type log struct {
	entries     []p.LogEntry
	lastApplied int
	commitIndex int64
	cConf       clusterconf.Configuration
	localFs     localfs.LocalFs
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
func (this *log) UpdateConfiguration(nodeIps []string) {
	this.cConf.UpdateConfiguration(nodeIps)
}

func (this *log) GetEntries() []*p.LogEntry {
	var e []*p.LogEntry = make([]*p.LogEntry, len(this.entries))
	for i, en := range this.entries {
		e[i] = &en
	}
	return e
}

func (this *log) LastLogIndex() int {
	return len(this.entries) - 1
}

func (this *log) More_recent_log(last_log_index int64, last_log_term uint64) bool {
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

func (this *log) AppendEntries(newEntries []*p.LogEntry, index int) {
	var indexEndQueue = len(this.GetEntries()) - 1
	if indexEndQueue == -1 {
		indexEndQueue = 0
	}
	fmt.Printf("log index last log: %v\n", indexEndQueue)

	this.entries = extend(this.entries, len(newEntries))
	for i, en := range newEntries {
		l.Printf("i: %d, i + index: %d, logEntry: %v", i, i+indexEndQueue, en.String())
		this.entries[indexEndQueue+i] = *en
	}

	l.Printf("my entries: %v, len: %d", this.GetEntries(), len(this.entries))
}

func (this *log) UpdateLastApplied() error {
	for int(this.commitIndex) > this.lastApplied {
		var entry *p.LogEntry = &this.entries[this.commitIndex]

		l.Printf("updating entry: %v", entry)
		switch entry.OpType {
		case p.Operation_JOIN_CONF:
			this.applyConf(entry)
		default:
			(*this).localFs.ApplyLogEntry(entry)
		}

		this.lastApplied++
	}
	return nil
}

func (this *log) GetCommitIndex() int64 {
	return this.commitIndex
}

func (this *log) SetCommitIndex(val int64) {
	this.commitIndex = val
}

// utility
func extend(slice []p.LogEntry, addedCapacity int) []p.LogEntry {
	//  l.Println("enxtend")
	n := len(slice)
	newSlice := make([]p.LogEntry, n+addedCapacity)
	for i := range slice {
		newSlice[i] = slice[i]
	}
	return newSlice
}

func (this *log) applyConf(entry *p.LogEntry) {
	var confUnfiltered string = string(entry.Payload)
	var confFiltered []string = strings.Split(confUnfiltered, " ")
	l.Printf("applying the new conf:%v\t%v\n", confUnfiltered, confFiltered)
	this.cConf.UpdateConfiguration(confFiltered)
}
