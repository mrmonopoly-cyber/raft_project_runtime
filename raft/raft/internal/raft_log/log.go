package raft_log

import p "raft/pkg/rpcEncoding/out/protobuf"

type LogEntry interface{
    GetEntries() []p.LogEntry
    GetCommitIndex() int64
    More_recent_log(last_log_index int64, last_log_term uint64) bool 
}

type log struct {
	entries     []p.LogEntry
	commitIndex int64
}

func (l *log) GetEntries() []p.LogEntry{
  return l.entries
}

func (l *log) GetCommitIndex() int64 {
  return l.commitIndex
}

func (l *log) More_recent_log(last_log_index int64, last_log_term uint64) bool {
    if l.GetCommitIndex() == -1 {
        return true
    }
    if last_log_index >= l.GetCommitIndex(){
        var lastEntry *p.LogEntry = &l.GetEntries()[last_log_index]
        if last_log_term >= *lastEntry.Term {
            return true
        }
    }
    return false
}

func NewLogEntry() LogEntry{
    return &log{
        commitIndex: -1,
    }
}
