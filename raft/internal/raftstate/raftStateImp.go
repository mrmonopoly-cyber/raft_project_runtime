package raftstate

import (
	"log"
	"math/rand"
	confpool "raft/internal/raftstate/confPool"
	"raft/internal/raftstate/timeout"
	"time"
)

type Role int

type ipMetadata struct {
	public  string
	private string
}

type raftStateImpl struct {
	myIp ipMetadata
	leaderIp ipMetadata

	term uint64
	role Role

	voteFor string
	voting  bool

	nSupporting    uint64
	nNotSupporting uint64

    timeout.TimeoutPool

    confpool.ConfPool
}

// GetLeaderIp implements State.
func (this *raftStateImpl) GetLeaderIp(vis VISIBILITY) string {
	switch vis {
	case PUB:
		return this.leaderIp.public
	case PRI:
		return this.leaderIp.private
	default:
		log.Panicln("invalid case ", vis)
		return ""
	}
}

// SetLeaderIp implements State.
func (this *raftStateImpl) SetLeaderIp(vis VISIBILITY, ip string) {
	switch vis {
	case PUB:
		this.leaderIp.public = ip
	case PRI:
		this.leaderIp.private = ip
	default:
		log.Panicln("unamanage case in setLeaderIp")
	}
}

// GetMyIp implements State.
func (this *raftStateImpl) GetMyIp(vis VISIBILITY) string {
	switch vis {
	case PUB:
		return this.myIp.public
	case PRI:
		return this.myIp.private
	default:
		log.Panicln("invalid case ip: ", vis)
		return ""
	}
}

func (this *raftStateImpl) GetTerm() uint64 {
	return this.term
}

func (this *raftStateImpl) SetTerm(newTerm uint64) {
	this.term = newTerm
}

func (this *raftStateImpl) GetRole() Role {
	return this.role
}

func (this *raftStateImpl) SetRole(newRole Role) {
	if this.role == newRole {
		return
	}
	switch newRole {
	case FOLLOWER:
		this.StopTimeout(TIMER_HEARTHBIT)
        this.ConfPool.AutoCommitSet(true)
	case LEADER:
		this.RestartTimeout(TIMER_HEARTHBIT)
		this.leaderIp = this.myIp
		this.ResetElection()
        this.ConfPool.AutoCommitSet(false)
	}
	this.role = newRole
}

func (this *raftStateImpl) CanVote() bool {
	return this.voting
}

func (this *raftStateImpl) VoteRight(vote bool) {
	(*this).voting = vote
}

func (this *raftStateImpl) GetVoteFor() string {
	return this.voteFor
}

func (this *raftStateImpl) IncrementTerm() {
	this.term += 1
}

func (this *raftStateImpl) VoteFor(id string) {
	this.voteFor = id
}

// GetNumSupporters implements State.
func (this *raftStateImpl) GetNumSupporters() uint64 {
	return this.nSupporting
}

// IncreaseNotSupporters implements State.
func (this *raftStateImpl) IncreaseNotSupporters() {
	this.nNotSupporting++
}

// IncreaseSupporters implements State.
func (this *raftStateImpl) IncreaseSupporters() {
	this.nSupporting++
}

// GetNumNotSupporters implements State.
func (this *raftStateImpl) GetNumNotSupporters() uint64 {
	return this.nNotSupporting
}

func (this *raftStateImpl) ResetElection() {
	this.nSupporting = 0
	this.nNotSupporting = 0
}


func newStateImplementation(idPrivate string, idPublic string, fsRootDir string) *raftStateImpl {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	var randelection = rand.Intn((int(MAX_ELECTION_TIMEOUT) - int(MIN_ELECTION_TIMEOUT) + 1)) + int(MIN_ELECTION_TIMEOUT)
	var s = new(raftStateImpl)

	s.term = 0

	s.myIp.private = idPrivate
	s.myIp.public = idPublic

	s.nNotSupporting = 0
	s.nSupporting = 0

	s.voting = true
	s.voteFor = ""

	s.leaderIp.public = ""
	s.leaderIp.private = ""

	s.TimeoutPool = timeout.NewTimeoutPool()
	s.TimeoutPool.AddTimeout(TIMER_ELECTION, time.Duration(randelection))
	s.TimeoutPool.AddTimeout(TIMER_HEARTHBIT, time.Duration(H_TIMEOUT))
	s.TimeoutPool.RestartTimeout(TIMER_HEARTHBIT)

    s.ConfPool = confpool.NewConfPoll(fsRootDir)
    s.SetRole(FOLLOWER)
    s.ConfPool.AutoCommitSet(true)

	return s
}
