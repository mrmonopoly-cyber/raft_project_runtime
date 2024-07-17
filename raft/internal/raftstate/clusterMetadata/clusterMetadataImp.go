package clustermetadata

import (
	"log"
	"math/rand"
	"raft/internal/raftstate/timeout"
	"sync"
	"time"
)

type clusterMetadataImp struct {
	myIp     ipMetadata
	leaderIp ipMetadata

	term uint64
	role Role

	voteFor string
	voting  bool

    lockSupporter sync.RWMutex
	nSupporting    uint
	nNotSupporting uint

	timeout.TimeoutPool
}

// GetNumNotSupporters implements ClusterMetadata.
func (c *clusterMetadataImp) GetNumNotSupporters() uint {
    c.lockSupporter.RLock()
    defer c.lockSupporter.RUnlock()

	return c.nNotSupporting
}

// GetNumSupporters implements ClusterMetadata.
func (c *clusterMetadataImp) GetNumSupporters() uint {
    c.lockSupporter.RLock()
    defer c.lockSupporter.RUnlock()

    return c.nSupporting
}

// UpdateSupportersNum implements ClusterMetadata.
func (c *clusterMetadataImp) UpdateSupportersNum(op SUPP_OP) {
    c.lockSupporter.Lock()
    defer c.lockSupporter.Unlock()

	switch op {
	case INC:
		c.nSupporting++
	case DEC:
		c.nNotSupporting++
	}
}

// CanVote implements ClusterMetadata.
func (c *clusterMetadataImp) CanVote() bool {
	return c.voting
}

// GetLeaderIp implements ClusterMetadata.
func (c *clusterMetadataImp) GetLeaderIp(vis VISIBILITY) string {
	switch vis {
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
	switch vis {
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

// ResetElection implements ClusterMetadata.
func (c *clusterMetadataImp) ResetElection() {
    c.lockSupporter.Lock()
    defer c.lockSupporter.Unlock()

	c.nNotSupporting = 0
	c.nSupporting = 1
	c.voteFor = ""
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
		log.Println("becomming follower")
		c.StopTimeout(TIMER_HEARTHBIT)
		c.RestartTimeout(TIMER_ELECTION)
	case LEADER:
		log.Println("becomming leader")
		c.RestartTimeout(TIMER_HEARTHBIT)
		c.leaderIp = c.myIp
		c.ResetElection()
	}
	c.role = newRole
}

// SetTerm implements ClusterMetadata.
func (c *clusterMetadataImp) SetTerm(newTerm uint64) {
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
	rand.New(rand.NewSource(time.Now().UnixNano()))
	var randelection = rand.Intn((int(MAX_ELECTION_TIMEOUT) - int(MIN_ELECTION_TIMEOUT) + 1)) + int(MIN_ELECTION_TIMEOUT)
	var res = &clusterMetadataImp{
		myIp: ipMetadata{
			public:  idPublic,
			private: idPrivate,
		},
		role:           FOLLOWER,
		term:           0,
		voteFor:        "",
		voting:         true,
        lockSupporter: sync.RWMutex{},
		nSupporting:    0,
		nNotSupporting: 0,
		TimeoutPool:    timeout.NewTimeoutPool(),
	}

	res.TimeoutPool.AddTimeout(TIMER_ELECTION, time.Duration(randelection))
	res.TimeoutPool.AddTimeout(TIMER_HEARTHBIT, time.Duration(H_TIMEOUT))
	res.SetRole(FOLLOWER)
	return res
}
