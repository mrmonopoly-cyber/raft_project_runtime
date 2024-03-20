package raftstate

import (
	l "raft/internal/raft_log"
	p "raft/pkg/protobuf"
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
	log              l.Log
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
	GetEntries() []p.Entry
	GetCommitIndex() uint64
	SetRole(newRole Role)
	SetTerm(newTerm uint64)
	MoreRecentLog(lastLogIndex uint64, lastLogTerm uint64) bool
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

func (_state *raftStateImpl) GetEntries() []p.Entry {
	return _state.log.GetEntries()
}

func (_state *raftStateImpl) GetCommitIndex() uint64 {
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
func (_state *raftStateImpl) MoreRecentLog(lastLogIndex uint64, lastLogTerm uint64) bool {
    return _state.log.More_recent_log(lastLogIndex,lastLogTerm)
}

func NewState(term uint64, id string, role Role) State {
	var s = new(raftStateImpl)
	s.role = role
	s.term = term
	s.id = id
	// s.serversIP = serversIp
	s.electionTimeout = time.NewTimer(ELECTION_TIMEOUT)
	s.heartbeatTimeout = time.NewTimer(H_TIMEOUT)

	return s
}
