package raft_log

import p "raft/pkg/rpcEncoding/out/protobuf"

type Log struct {
	entries     []p.LogEntry
}

type LogInterface interface {
  GetEntries() []p.Entry
}

func (l *Log) GetEntries() []p.LogEntry{
  return l.entries
}

func (l *Log) LastLogIndex() int {
  return len(l.entries)-1
}

func (l *Log) More_recent_log(last_log_index uint64, last_log_term uint64) bool {
    //TODO implement More recent log
    return false
}

func (l *Log) AppendEntries(newEntries []*p.LogEntry, index int) {
  for i, en := range newEntries {
    l.entries[index + i + 1] = *en
  }
}
