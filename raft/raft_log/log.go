package raft_log

import p "raft/protobuf"

type Log struct {
	entries     []p.Entry
	commitIndex uint64
}

type LogInterface interface {
  GetEntries() []p.Entry
  GetCommitIndex() uint64
}

func (l *Log) GetEntries() []p.Entry {
  return l.entries
}

func (l *Log) GetCommitIndex() uint64 {
  return l.commitIndex
}
