package raftstate

import (
	"math/rand"
	l "raft/internal/raft_log"
	"time"
)

const (
	LEADER Role = iota
	FOLLOWER
	CANDIDATE
)

const (
	MIN_ELECTION_TIMEOUT time.Duration = 10000000000
	MAX_ELECTION_TIMEOUT time.Duration = 15000000000
	H_TIMEOUT            time.Duration = 3000000000
)

type State interface {
	GetIdPrivate() string
	GetIdPublic() string
	GetTerm() uint64
	GetRole() Role
	StartElectionTimeout()
	StartHearthbeatTimeout()
	StopElectionTimeout()
	StopHearthbeatTimeout()
	Leader() bool
	HeartbeatTimeout() *time.Timer
	ElectionTimeout() *time.Timer
	GetVoteFor() string
	IncrementTerm()
	VoteFor(id string)
	CanVote() bool
    ToggleVoteRight()
	SetRole(newRole Role)
	SetTerm(newTerm uint64)
	IncreaseSupporters()
	IncreaseNotSupporters()
	GetNumSupporters() uint64
	GetNumNotSupporters() uint64
	ResetElection()

    CheckCommitIndex(idxList []int)
    l.LogEntry

	GetLeaderIpPrivate() string
	GetLeaderIpPublic() string
	SetLeaderIpPublic(ip string)
	SetLeaderIpPrivate(ip string)
}



func NewState(term uint64, idPrivate string, idPublic string, role Role, fsRootDir string) State {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	var s = new(raftStateImpl)
	s.role = role
	s.term = term
	s.idPrivate = idPrivate
	s.idPublic = idPublic
	s.electionTimeout = time.NewTimer(MAX_ELECTION_TIMEOUT)
	s.heartbeatTimeout = time.NewTimer(H_TIMEOUT)
	s.nNotSupporting = 0
	s.nSupporting = 0
	s.nNodeInCluster = 1
	s.voting = true
	s.log = l.NewLogEntry([]string{idPrivate})
	s.electionTimeoutRaw = rand.Intn((int(MAX_ELECTION_TIMEOUT) - int(MIN_ELECTION_TIMEOUT) + 1)) + int(MIN_ELECTION_TIMEOUT)
	return s
}
