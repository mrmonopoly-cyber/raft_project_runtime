package raft_log

import p "raft/pkg/rpcEncoding/out/protobuf"

type LogEntry interface{
    GetEntries() []p.LogEntry
    GetCommitIndex() int64
    More_recent_log(last_log_index int64, last_log_term uint64) bool 
}

type log struct {
    entries     []p.LogEntry
    lastApplied  int
    commitIndex  uint64
}

func (this *log) GetEntries() []p.LogEntry{
  return this.entries
}

func (this *log) LastLogIndex() int {
  return len(this.entries)-1
}

func (this *log) More_recent_log(last_log_index uint64, last_log_term uint64) bool {
    //TODO implement More recent log
    return false
}

func (this *log) AppendEntries(newEntries []*p.LogEntry, index int) {
  for i, en := range newEntries {
    this.entries[index + i + 1] = *en
  }
}

func (this *log) UpdateLastApplied() int {
  if int(this.commitIndex) > this.lastApplied {
    this.lastApplied++
    return this.lastApplied
  }
  return -1
}

func (this *log) GetCommitIndex() uint64 {
  return this.commitIndex
}

func (this *log) SetCommitIndex(val uint64) {
  this.commitIndex = val
}

func (this *log) InitState() {
  this.commitIndex = 0
  this.lastApplied = 0
}
