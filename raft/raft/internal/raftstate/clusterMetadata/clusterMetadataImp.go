package clustermetadata

import (
	"log"
	"raft/internal/raftstate/timeout"
)

type clusterMetadataImp struct {
	myIp     ipMetadata
	leaderIp ipMetadata

	term uint64
	role Role

	voteFor string
	voting  bool

	nSupporting    uint64
	nNotSupporting uint64

    timeout.TimeoutPool
}

// CanVote implements ClusterMetadata.
func (c *clusterMetadataImp) CanVote() bool {
    return c.voting
}

// GetLeaderIp implements ClusterMetadata.
func (c *clusterMetadataImp) GetLeaderIp(vis VISIBILITY) string {
    switch vis{
    case PUB:
        return c.leaderIp.public
    case PRI:
        return c.leaderIp.private
    default:
        log.Panicln("unmanaged case: ", vis)
        return ""
    }
}

// GetMyIp implements ClusterMetadata.
func (c *clusterMetadataImp) GetMyIp(vis VISIBILITY) string {
    switch vis{
    case PUB:
        return c.myIp.public
    case PRI:
        return c.myIp.private
    default:
        log.Panicln("unmanaged case: ", vis)
        return ""
    }
}

// GetRole implements ClusterMetadata.
func (c *clusterMetadataImp) GetRole() Role {
    return c.role
}

// GetTerm implements ClusterMetadata.
func (c *clusterMetadataImp) GetTerm() uint64 {
    return c.term
}

// GetVoteFor implements ClusterMetadata.
func (c *clusterMetadataImp) GetVoteFor() string {
    return c.voteFor
}

// IncrementTerm implements ClusterMetadata.
func (c *clusterMetadataImp) IncrementTerm() {
    c.term++
}

// ResetElection implements ClusterMetadata.
func (c *clusterMetadataImp) ResetElection() {
    c.nNotSupporting = 0
    c.nSupporting = 0
}

// SetLeaderIp implements ClusterMetadata.
func (c *clusterMetadataImp) SetLeaderIp(vis VISIBILITY, ip string) {
	switch vis {
	case PUB:
		c.leaderIp.public = ip
	case PRI:
		c.leaderIp.private = ip
	default:
		log.Panicln("unamanage case in setLeaderIp")
	}
}

// SetRole implements ClusterMetadata.
func (c *clusterMetadataImp) SetRole(newRole Role) {
	if c.role == newRole {
		return
	}
	switch newRole {
	case FOLLOWER:
		c.StopTimeout(TIMER_HEARTHBIT)
	case LEADER:
		c.RestartTimeout(TIMER_HEARTHBIT)
		c.leaderIp = c.myIp
		c.ResetElection()
	}
	c.role = newRole
}

// SetTerm implements ClusterMetadata.
func (c *clusterMetadataImp) SetTerm(newTerm uint64) {
    //HACK: to remove legacy code
	c.term = newTerm
}

// VoteFor implements ClusterMetadata.
func (c *clusterMetadataImp) VoteFor(id string) {
    c.voteFor = id
}

// VoteRight implements ClusterMetadata.
func (c *clusterMetadataImp) VoteRight(vote bool) {
    c.voting = vote
}

type ipMetadata struct {
	public  string
	private string
}

func newClusterMetadataImp(idPrivate string, idPublic string) *clusterMetadataImp {
	return &clusterMetadataImp{
        myIp: ipMetadata{
            public: idPublic,
            private: idPrivate,
        },
        TimeoutPool: timeout.NewTimeoutPool(),
    }
}
