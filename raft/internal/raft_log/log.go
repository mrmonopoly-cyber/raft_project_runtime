package raft_log

import p "raft/pkg/rpcEncoding/out/protobuf"

type Log struct {
	entries     []p.LogEntry
  lastApplied  int
  commitIndex  uint64
}

func (this *Log) GetEntries() []p.LogEntry{
  return this.entries
}

func (this *Log) LastLogIndex() int {
  return len(this.entries)-1
}

func (this *Log) More_recent_log(last_log_index uint64, last_log_term uint64) bool {
    //TODO implement More recent log
    return false
}

func (this *Log) AppendEntries(newEntries []*p.LogEntry, index int) {
  for i, en := range newEntries {
    this.entries[index + i + 1] = *en
  }
}

func (this *Log) UpdateLastApplied() int {
  if int(this.commitIndex) > this.lastApplied {
    this.lastApplied++
    return this.lastApplied
  }
  return -1
}

func (this *Log) GetCommitIndex() uint64 {
  return this.commitIndex
}

func (this *Log) SetCommitIndex(val uint64) {
  this.commitIndex = val
}

func (this *Log) InitState() {
  this.commitIndex = 0
  this.lastApplied = 0
}

