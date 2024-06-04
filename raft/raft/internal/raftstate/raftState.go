package raftstate

import (
	confpool "raft/internal/raftstate/confPool"
	"raft/internal/raftstate/timeout"
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

type VISIBILITY int
const (
    PUB VISIBILITY = iota
    PRI VISIBILITY = iota
)

const (
    TIMER_ELECTION = "election"
    TIMER_HEARTHBIT = "hearthbit"
)

type State interface {
    leaderIpMetadata
    currentNodeIp

    termMetadata
    voteMetadata
    roleMetadata

    timeout.TimeoutPool

    confpool.ConfPool
}

type leaderIpMetadata interface{
    SetLeaderIp(vis VISIBILITY, ip string)
    GetLeaderIp(vis VISIBILITY) string
}

type roleMetadata interface{
    GetRole() Role
    SetRole(newRole Role)
}

type voteMetadata interface{
    GetVoteFor() string
    VoteFor(id string)
    CanVote() bool
    VoteRight(vote bool)
    ResetElection()
}

type termMetadata interface{
    GetTerm() uint64
    IncrementTerm()
    SetTerm(newTerm uint64)
}

type currentNodeIp interface{
    GetMyIp(vis VISIBILITY) string
}


func NewState(idPrivate string, idPublic string, fsRootDir string) State {
	return newStateImplementation(idPrivate, idPublic, fsRootDir)
}
