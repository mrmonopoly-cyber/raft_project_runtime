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
	entries     []*p.LogEntry
	lastApplied int
	commitIndex int64
	cConf       clusterconf.Configuration
	localFs     localfs.LocalFs
}

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

func (this *log) GetEntries() []*p.LogEntry {
	return this.entries
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

func (this *log) AppendEntries(newEntries []*p.LogEntry) {
	var indexEndQueue = len(this.GetEntries()) - 1
	if indexEndQueue == -1 {
		indexEndQueue = 0
	}

	fmt.Printf("log index last log: %v\n", indexEndQueue)
    this.entries = append(this.entries,newEntries...)
	l.Printf("my entries: %v, len: %d", this.GetEntries(), len(this.entries))

	if this.commitIndex == -1 {
		this.commitIndex = 0
	}
}

func (this *log) UpdateLastApplied() error {
	l.Printf("check if can apply some logEntry: commIndex:%v, lastApplied:%v\n", this.commitIndex, this.lastApplied)
	for int(this.commitIndex) > this.lastApplied {
		var entry *p.LogEntry = this.entries[this.commitIndex]

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
func (this *log) applyConf(entry *p.LogEntry) {
	var confUnfiltered string = string(entry.Payload)
	var confFiltered []string = strings.Split(confUnfiltered, " ")
	l.Printf("applying the new conf:%v\t%v\n", confUnfiltered, confFiltered)
    l.Println("for debugging reasong now it's only adding node to the conf")
	this.cConf.UpdateConfiguration(clusterconf.ADD, confFiltered)
    //WARN: only adding to conf
}
