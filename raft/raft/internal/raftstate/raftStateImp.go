package raftstate

import (
	"log"
	"math/rand"
	"raft/internal/raft_log"
	l "raft/internal/raft_log"
	nodematchidx "raft/internal/raftstate/nodeMatchIdx"
	"raft/internal/raftstate/timeout"
	"time"
)

type Role int

type raftStateImpl struct {
	myIp ipMetadata

	term uint64
	role Role

	voteFor string
	voting  bool

	l.LogEntry

	nSupporting    uint64
	nNotSupporting uint64

	leaderMetadata leader

	timeout.TimeoutPool

	nodeToUpdateChan chan NewNodeToUpdateInfo
}

// NotifyNodeToUpdate implements State.
func (this *raftStateImpl) NotifyNodeToUpdate(nodeIps []string) {
    go func() {
        this.nodeToUpdateChan <- NewNodeToUpdateInfo{
            NodeList: nodeIps,
            MatchToArrive: uint64(this.LastLogIndex()),
        }
    }()
}

// GetNewNodeToUpdate implements State.
func (this *raftStateImpl) GetNewNodeToUpdate() <-chan NewNodeToUpdateInfo {
	return this.nodeToUpdateChan
}

// GetLeaderIp implements State.
func (this *raftStateImpl) GetLeaderIp(vis VISIBILITY) string {
	switch vis {
	case PUB:
		return this.leaderMetadata.leaderIp.public
	case PRI:
		return this.leaderMetadata.leaderIp.private
	default:
		log.Println("invalid case ", vis)
		return ""
	}
}

// SetLeaderIp implements State.
func (this *raftStateImpl) SetLeaderIp(vis VISIBILITY, ip string) {
	switch vis {
	case PUB:
		this.leaderMetadata.leaderIp.public = ip
	case PRI:
		this.leaderMetadata.leaderIp.private = ip
	default:
		log.Panicln("unamanage case in setLeaderIp")
	}
}

type leader struct {
	leaderEntryToCommit chan int64
	statePool           nodematchidx.NodeCommonMatch
	leaderIp            ipMetadata
}

type ipMetadata struct {
	public  string
	private string
}

// GetStatePool implements State.
func (this *raftStateImpl) GetStatePool() nodematchidx.NodeCommonMatch {
	return this.leaderMetadata.statePool
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
		this.GetStatePool().InitCommonMatch(this.LastLogIndex())
	case LEADER:
		this.RestartTimeout(TIMER_HEARTHBIT)
		this.leaderMetadata.leaderIp = this.myIp
		this.ResetElection()
	}
	this.role = newRole
}

func (this *raftStateImpl) AppendEntries(newEntries []*l.LogInstance) {
    this.LogEntry.AppendEntries(newEntries)
	if this.role == FOLLOWER || this.GetNumberNodesInCurrentConf() <= 1{
		log.Println("auto commit entry: ", newEntries[0].Entry)
		for range newEntries {
			if this.role == LEADER {
				this.leaderMetadata.statePool.IncreaseCommonMathcIndex()
			}
			this.LogEntry.IncreaseCommitIndex()
		}
		return
	}
	log.Printf("leader, request to send log Entry to follower: ch %v, idx: %v\n",
		this.leaderMetadata.leaderEntryToCommit, this.LogEntry.GetCommitIndex()+1)
	for range newEntries {
		this.leaderMetadata.leaderEntryToCommit <- this.LogEntry.GetCommitIndex() + 1
	}
}

func (this *raftStateImpl) GetLeaderEntryChannel() *chan int64 {
	return &this.leaderMetadata.leaderEntryToCommit
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

func (this *raftStateImpl) leaaderUpdateCommitIndex() {
	for {
		select {
		case <-this.leaderMetadata.statePool.GetNotifyChannel():
			this.IncreaseCommitIndex()
		}
	}
}

func (this *raftStateImpl) checkPresenceInTheConf(){
    for{
        <- this.NotifyChangeInConf()
        var conf = this.GetConfig()
        var found = false
        for _, v := range conf {
            if v == this.myIp.private{
               found = true 
               continue
            }
        }
        
        if !found{
            this.SetRole(FOLLOWER)
            this.StopTimeout(TIMER_ELECTION)
            this.LogEntry = raft_log.NewLogEntry(this.GetRootDirFs())

            this.nNotSupporting = 0
            this.nSupporting = 0

            this.voting = true
            this.voteFor = ""

            this.leaderMetadata.leaderIp.public = ""
            this.leaderMetadata.leaderIp.private = ""
            this.leaderMetadata.statePool = nodematchidx.NewNodeCommonMatch()
        }

    }
}

func newStateImplementation(idPrivate string, idPublic string, fsRootDir string) *raftStateImpl {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	var randelection = rand.Intn((int(MAX_ELECTION_TIMEOUT) - int(MIN_ELECTION_TIMEOUT) + 1)) + int(MIN_ELECTION_TIMEOUT)
	var s = new(raftStateImpl)

	s.role = FOLLOWER
	s.term = 0

	s.myIp.private = idPrivate
	s.myIp.public = idPublic

	s.nNotSupporting = 0
	s.nSupporting = 0

	s.voting = true
	s.voteFor = ""

	s.leaderMetadata.leaderIp.public = ""
	s.leaderMetadata.leaderIp.private = ""
	s.leaderMetadata.statePool = nodematchidx.NewNodeCommonMatch()
	s.leaderMetadata.leaderEntryToCommit = make(chan int64)

	s.LogEntry = l.NewLogEntry(fsRootDir)

	s.TimeoutPool = timeout.NewTimeoutPool()
	s.TimeoutPool.AddTimeout(TIMER_ELECTION, time.Duration(randelection))
	s.TimeoutPool.AddTimeout(TIMER_HEARTHBIT, time.Duration(H_TIMEOUT))
	s.TimeoutPool.RestartTimeout(TIMER_HEARTHBIT)

    s.nodeToUpdateChan = make(chan NewNodeToUpdateInfo)

	go s.leaaderUpdateCommitIndex()

	return s
}
