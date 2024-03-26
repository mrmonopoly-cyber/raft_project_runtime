package raft_log

import p "raft/pkg/rpcEncoding/out/protobuf"

type Log struct {
	entries     []p.LogEntry
    entriesNum uint64
	commitIndex uint64
}

func (l *Log) GetEntries() []p.LogEntry{
  return l.entries
}

func (l *Log) GetCommitIndex() uint64 {
  return l.commitIndex
}

func (l *Log) More_recent_log(last_log_index uint64, last_log_term uint64) bool {
    if last_log_index >= l.entriesNum {
        var lastEntry *p.LogEntry = &l.GetEntries()[last_log_index]
        if last_log_term >= *lastEntry.Term {
            return true
        }
    }
    return false
}
