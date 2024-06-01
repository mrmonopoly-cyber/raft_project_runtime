package raftstate

import (
	l "raft/internal/raft_log"
	nodematchidx "raft/internal/raftstate/nodeMatchIdx"
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
    currentNodeIp
    termMetadata
    voteMetadata
    roleMetadata
    l.LogEntry
    timeout.TimeoutPool

    leaderRoleExtra
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
    GetIdPrivate() string
    GetIdPublic() string
}

type leaderRoleExtra interface{
    leaderIpMetadata

    GetLeaderEntryChannel() *chan int64
    GetStatePool() nodematchidx.NodeCommonMatch

    IncreaseSupporters()
    IncreaseNotSupporters()
    GetNumSupporters() uint64
    GetNumNotSupporters() uint64
}

func NewState(idPrivate string, idPublic string, fsRootDir string) State {
	return newStateImplementation(idPrivate, idPublic, fsRootDir)
}
