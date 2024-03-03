package raft_log

type Log struct {
	entries     []Entry
	commitIndex int
}
