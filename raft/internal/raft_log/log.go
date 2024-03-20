package raft_log

import p "raft/pkg/protobuf"

type Log struct {
	entries     []p.Entry
	commitIndex uint64
}

func (l *Log) GetEntries() []p.Entry {
  return l.entries
}

func (l *Log) GetCommitIndex() uint64 {
  return l.commitIndex
}

func (l *Log) More_recent_log(last_log_index uint64, last_log_term uint64) bool {
    //TODO implement More recent log
    return false
}

func (l *Log) AppendEntries(newEntries []*p.Entry) {
  for _, en := range newEntries {
    l.entries = append(l.entries, *en)
  }
}
