package raft_log

import p "raft/protobuf"

type Log struct {
	entries     []p.Entry
	commitIndex int
}
