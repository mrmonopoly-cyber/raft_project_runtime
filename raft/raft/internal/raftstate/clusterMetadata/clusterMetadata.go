package clustermetadata

import (
	"raft/internal/raftstate/timeout"
)


type ClusterMetadata interface{
    leaderIpMetadata
    currentNodeIp

    termMetadata
    voteMetadata
    roleMetadata

    timeout.TimeoutPool
}

const (
    TIMER_ELECTION = "election"
    TIMER_HEARTHBIT = "hearthbit"
)

type Role uint
const (
	LEADER Role = iota
	FOLLOWER
	CANDIDATE
)

type VISIBILITY int
const (
    PUB VISIBILITY = iota
    PRI VISIBILITY = iota
)

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

func NewClusterMetadata(idPrivate string, idPublic string) ClusterMetadata{
    return newClusterMetadataImp(idPrivate,idPublic)
}
