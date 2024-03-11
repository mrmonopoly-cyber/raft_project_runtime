package server

import (
	"log"
	"time"
  l "raft/raft_log"
  p "raft/protobuf"
)

type Role int

const (
	LEADER Role = iota
	FOLLOWER
	CANDIDATE
)

const (
	ELECTION_TIMEOUT time.Duration = 200
	H_TIMEOUT        time.Duration = 30
)

type raftStateImpl struct {
	id   string
	term uint64
	leaderId         string
	role Role
	voteFor          string
	//voting           bool
	serversID []string
	//updated_node     []net.IP
	//updating_node    []net.IP
	electionTimeout  *time.Timer
	heartbeatTimeout *time.Timer
  log              l.Log
}

type state interface {
	GetID() string
	GetTerm() uint64
	//GetLeaderIP() net.IP
	GetRole() Role
	StartElectionTimeout()
	StartHearthbeatTimeout()
	Leader() bool
	HeartbeatTimeout() *time.Timer
	ElectionTimeout() *time.Timer
	GetServersID() []string
	NewState() state
  GetVoteFor() string
  IncrementTerm()
  VoteFor(id string)
  GetEntries() []p.Entry
  GetCommitIndex() uint64
  SetRole(newRole Role)
}

func (_state *raftStateImpl) GetID() string {
	return _state.id
}

func (_state *raftStateImpl) GetTerm() uint64 {
	return _state.term
}

func (_state *raftStateImpl) GetRole() Role {
	return _state.role
}

func (_state * raftStateImpl) SetRole(newRole Role) {
  _state.role = newRole
}

func (_state *raftStateImpl) GetEntries() []p.Entry {
  return _state.log.GetEntries()
}

func (_state *raftStateImpl) GetCommitIndex() uint64 {
  return _state.log.GetCommitIndex()
}

func (_state *raftStateImpl) StartElectionTimeout() {
	log.Println("Resetting election timer from", _state.id)
	_state.electionTimeout.Reset(ELECTION_TIMEOUT)
}

func (_state *raftStateImpl) StartHearthbeatTimeout() {
	_state.heartbeatTimeout.Reset(H_TIMEOUT)
}

func (_state *raftStateImpl) Leader() bool {
	return _state.role == LEADER
}

func (_state *raftStateImpl) HeartbeatTimeout() *time.Timer {
	return _state.heartbeatTimeout
}

func (_state *raftStateImpl) ElectionTimeout() *time.Timer {
	return _state.electionTimeout
}

func (_state *raftStateImpl) GetServersID() []string {
	return _state.serversID
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

func NewState(term uint64, id string, role Role, serversId []string) *raftStateImpl {
	var s = new(raftStateImpl)
	s.id = id
	s.role = role
	s.term = term
	s.serversID = serversId
	s.electionTimeout = time.NewTimer(ELECTION_TIMEOUT)
	s.heartbeatTimeout = time.NewTimer(H_TIMEOUT)

	return s
}
