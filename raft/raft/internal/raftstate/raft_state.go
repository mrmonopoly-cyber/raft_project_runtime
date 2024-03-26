package raftstate

import (
	"log"
	l "raft/internal/raft_log"
	p "raft/pkg/rpcEncoding/out/protobuf"
	"time"
)

type Role int

const (
	LEADER Role = iota
	FOLLOWER
	CANDIDATE
)

const (
	ELECTION_TIMEOUT time.Duration = 10000000000
	H_TIMEOUT        time.Duration = 3000000000
)

type raftStateImpl struct {
	id               string
	term             uint64
	leaderId         string
	role             Role
	voteFor          string
	voting           bool
	electionTimeout  *time.Timer
	heartbeatTimeout *time.Timer
	log              l.LogEntry
	nSupporting      uint64
	nNotSupporting   uint64
	nNodeInCluster   uint64
}


type State interface {
	GetId() string
	GetTerm() uint64
	GetRole() Role
	StartElectionTimeout()
	StartHearthbeatTimeout()
	Leader() bool
	HeartbeatTimeout() *time.Timer
	ElectionTimeout() *time.Timer
	GetVoteFor() string
	IncrementTerm()
	VoteFor(id string)
	CanVote() bool
	GetEntries() []p.LogEntry
	GetCommitIndex() int64
	SetRole(newRole Role)
	SetTerm(newTerm uint64)
	MoreRecentLog(lastLogIndex int64, lastLogTerm uint64) bool
	IncreaseSupporters()
	IncreaseNotSupporters()
	IncreaseNodeInCluster()
	GetNumSupporters() uint64
	GetNumNotSupporters() uint64
	GetNumNodeInCluster() uint64
	ResetElection()
}

func (_state *raftStateImpl) GetId() string {
	return _state.id
}

func (_state *raftStateImpl) GetTerm() uint64 {
	return _state.term
}

func (_state *raftStateImpl) SetTerm(newTerm uint64) {
	_state.term = newTerm
}

func (_state *raftStateImpl) GetRole() Role {
	return _state.role
}

func (_state *raftStateImpl) SetRole(newRole Role) {
	_state.role = newRole
}

func (_state *raftStateImpl) GetEntries() []p.LogEntry{
	return _state.log.GetEntries()
}

func (_state *raftStateImpl) GetCommitIndex() int64 {
	return _state.log.GetCommitIndex()
}

func (_state *raftStateImpl) StartElectionTimeout() {
	_state.electionTimeout.Reset(ELECTION_TIMEOUT)
}

func (_state *raftStateImpl) StartHearthbeatTimeout() {
	_state.heartbeatTimeout.Reset(H_TIMEOUT)
}

func (_state *raftStateImpl) Leader() bool {
	return _state.role == LEADER
}

func (_state *raftStateImpl) CanVote() bool {
	return _state.voting
}

func (_state *raftStateImpl) HeartbeatTimeout() *time.Timer {
    log.Println("timeout hearthbit")
	return _state.heartbeatTimeout
}

func (_state *raftStateImpl) ElectionTimeout() *time.Timer {
	return _state.electionTimeout
}

func (_state *raftStateImpl) GetVoteFor() string {
	return _state.voteFor
}

func (_state *raftStateImpl) IncrementTerm() {
	_state.term += 1
}

func (_state *raftStateImpl) VoteFor(id string) {
	_state.voteFor = id
}

// MoreRecentLog implements State.
func (_state *raftStateImpl) MoreRecentLog(lastLogIndex int64, lastLogTerm uint64) bool {
	return _state.log.More_recent_log(lastLogIndex, lastLogTerm)
}

// GetNumSupporters implements State.
func (_state *raftStateImpl) GetNumSupporters() uint64 {
	return _state.nSupporting
}

// IncreaseNotSupporters implements State.
func (_state *raftStateImpl) IncreaseNotSupporters() {
	_state.nNotSupporting++
}

// IncreaseSupporters implements State.
func (_state *raftStateImpl) IncreaseSupporters() {
	_state.nSupporting++
}

// GetNumNotSupporters implements State.
func (_state *raftStateImpl) GetNumNotSupporters() uint64 {
	return _state.nNotSupporting
}

// GetNumNodeInCluster implements State.
func (_state *raftStateImpl) GetNumNodeInCluster() uint64 {
	return _state.nNodeInCluster
}

// IncreaseNodeInCluster implements State.
func (_state *raftStateImpl) IncreaseNodeInCluster() {
	_state.nNodeInCluster++
}

func (_state *raftStateImpl) ResetElection() {
	_state.nSupporting = 0
	_state.nNotSupporting = 0
}

func NewState(term uint64, id string, role Role) State {
	var s = new(raftStateImpl)
	s.role = role
	s.term = term
	s.id = id
	// s.serversIP = serversIp
	s.electionTimeout = time.NewTimer(ELECTION_TIMEOUT)
	s.heartbeatTimeout = time.NewTimer(H_TIMEOUT)
	s.nNotSupporting = 0
	s.nSupporting = 0
	s.nNodeInCluster = 1
    s.voting = true
    s.log = l.NewLogEntry()

	return s
}
