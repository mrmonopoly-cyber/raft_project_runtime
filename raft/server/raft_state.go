package server

import (
	"log"
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
	id   int
	term int
	//leaderIP         net.IP
	role Role
	//voteFor          int
	//voting           bool
	serversID []int
	//updated_node     []net.IP
	//updating_node    []net.IP
	electionTimeout  *time.Timer
	heartbeatTimeout *time.Timer
}

type state interface {
	GetID() int
	GetTerm() int
	//GetLeaderIP() net.IP
	GetRole() Role
	StartElectionTimeout()
	StartHearthbeatTimeout()
	Leader() bool
	HeartbeatTimeout() *time.Timer
	ElectionTimeout() *time.Timer
	GetServersID() []int
	NewState() state
}

func (_state *raftStateImpl) GetID() int {
	return _state.id
}

func (_state *raftStateImpl) GetTerm() int {
	return _state.term
}

func (_state *raftStateImpl) GetRole() Role {
	return _state.role
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

func (_state *raftStateImpl) GetServersID() []int {
	return _state.serversID
}

func NewState(term int, id int, role Role, serversId []int) *raftStateImpl {
	var s = new(raftStateImpl)
	s.id = id
	s.role = role
	s.term = term
	s.serversID = serversId
	s.electionTimeout = time.NewTimer(ELECTION_TIMEOUT)
	s.heartbeatTimeout = time.NewTimer(H_TIMEOUT)

	return s
}
