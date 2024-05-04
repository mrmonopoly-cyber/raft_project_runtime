package raft_log

import (
	l "log"
	localfs "raft/internal/localFs"
	clusterconf "raft/internal/raftstate/clusterConf"
	p "raft/pkg/raft-rpcProtobuf-messages/rpcEncoding/out/protobuf"
	"strings"
	"sync"
)

type log struct {
	wg                 sync.WaitGroup
	uncommittedEntries chan *p.LogEntry
	commitEntries      []*p.LogEntry
	lastApplied        int
	commitIndex        int64
	cConf              clusterconf.Configuration
	localFs            localfs.LocalFs
}

// ConfStatus implements LogEntry.
func (this *log) ConfStatus() bool {
    return this.cConf.ConfStatus()
}

// IsInConf implements LogEntry.
func (this *log) IsInConf(ip string) bool {
	for _, v := range this.cConf.GetConfig() {
		if v == ip {
			return true
		}
	}
	return false
}

// GetClusterConfig implements LogEntry.
func (this *log) GetClusterConfig() []string {
	return this.cConf.GetConfig()
}

func (this *log) GetEntries() []*p.LogEntry {
	var e []*p.LogEntry = make([]*p.LogEntry, len(this.commitEntries))
	for i, en := range this.commitEntries {
		e[i] = en
	}
	return e
}

func (this *log) LastLogIndex() int {
	return len(this.commitEntries) - 1
}

func (this *log) More_recent_log(last_log_index int64, last_log_term uint64) bool {
	if last_log_index >= this.commitIndex {
		var commitEntries []*p.LogEntry = this.GetEntries()
		if len(commitEntries) <= int(last_log_index) {
			return true
		}
		if last_log_term >= (commitEntries[last_log_index]).Term {
			return true
		}
	}
	return false

}

func (this *log) AppendEntries(newEntries []*p.LogEntry) {
	for _, en := range newEntries {
        l.Printf("append new Log Entry: %v\n",en)
		this.uncommittedEntries <- en
	}
}

func (this *log) GetCommitIndex() int64 {
	return this.commitIndex
}

// utility
func (this *log) applyLogEntry() {
	for {
		var entry *p.LogEntry = <-(*this).uncommittedEntries
        l.Printf("applying logEntry %v\n",entry)
		switch entry.GetOpType() {
		case p.Operation_JOIN_CONF:
			var confUnfiltered string = string(entry.Payload)
			var confFiltered []string = strings.Split(confUnfiltered, " ")
			l.Printf("applying the new conf:%v\t%v\n", confUnfiltered, confFiltered)
			this.cConf.UpdateConfiguration(confFiltered)
		case p.Operation_CREATE, p.Operation_READ, p.Operation_WRITE, p.Operation_DELETE:
			this.localFs.ApplyLogEntry(entry)
		}
		this.commitEntries = append(this.commitEntries, entry)
		this.lastApplied++
		this.commitIndex++
	}
}

func extend(slice []*p.LogEntry, addedCapacity int) []*p.LogEntry {
	//  l.Println("enxtend")
	n := len(slice)
	newSlice := make([]*p.LogEntry, n+addedCapacity)
	for i := range slice {
		newSlice[i] = slice[i]
	}
	return newSlice
}
