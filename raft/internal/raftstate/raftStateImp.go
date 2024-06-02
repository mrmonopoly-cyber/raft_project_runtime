package raftstate

import (
	"log"
	"math/rand"
	l "raft/internal/raft_log"
	nodematchidx "raft/internal/raftstate/nodeMatchIdx"
	"raft/internal/raftstate/timeout"
	"time"

	"raft/pkg/raft-rpcProtobuf-messages/rpcEncoding/out/protobuf"
)

type Role int

type raftStateImpl struct {
	myIp ipMetadata

	term uint64
	role Role

	voteFor string
	voting  bool

	log l.LogEntry

	nSupporting    uint64
	nNotSupporting uint64

	leaderMetadata leader

	timeout.TimeoutPool
}

// GetLeaderIp implements State.
func (this *raftStateImpl) GetLeaderIp(vis VISIBILITY) string {
    switch vis{
    case PUB:
        return this.leaderMetadata.leaderIp.public
    case PRI:      
        return this.leaderMetadata.leaderIp.private
    default:
        log.Println("invalid case ", vis)
        return ""
    }
}

// GetTimeoutNotifycationChan implements State.
// Subtle: this method shadows the method (TimeoutPool).GetTimeoutNotifycationChan of raftStateImpl.TimeoutPool.
func (this *raftStateImpl) GetTimeoutNotifycationChan(name string) (chan time.Time, error) {
    return this.TimeoutPool.GetTimeoutNotifycationChan(name)
}

// SetLeaderIp implements State.
func (this *raftStateImpl) SetLeaderIp(vis VISIBILITY, ip string) {
    switch vis{
    case PUB:
        this.myIp.public = ip
    case PRI:
        this.myIp.private = ip
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

// NewLogInstance implements State.
func (this *raftStateImpl) NewLogInstance(entry *protobuf.LogEntry) *l.LogInstance {
	return this.log.NewLogInstance(entry)
}

// NewLogInstanceBatch implements State.
func (this *raftStateImpl) NewLogInstanceBatch(entry []*protobuf.LogEntry) []*l.LogInstance {
	return this.log.NewLogInstanceBatch(entry)
}

// GetCommittedEntriesRange implements State.
func (this *raftStateImpl) GetCommittedEntriesRange(startIndex int) []l.LogInstance {
	return this.log.GetCommittedEntriesRange(startIndex)
}

// LastLogTerm implements State.
func (this *raftStateImpl) LastLogTerm() uint {
	return this.log.LastLogTerm()
}

// GetCommittedEntries implements State.
func (this *raftStateImpl) GetCommittedEntries() []l.LogInstance {
	return this.log.GetCommittedEntries()
}

// GetStatePool implements State.
func (this *raftStateImpl) GetStatePool() nodematchidx.NodeCommonMatch {
	return this.leaderMetadata.statePool
}

// GetEntriAt implements State.
func (this *raftStateImpl) GetEntriAt(index int64) (*l.LogInstance, error) {
	return this.log.GetEntriAt(index)
}

// IncreaseCommitIndex implements State.
func (this *raftStateImpl) IncreaseCommitIndex() {
	this.log.IncreaseCommitIndex()
}

// MinimumCommitIndex implements State.
func (this *raftStateImpl) MinimumCommitIndex(val uint) {
	this.log.MinimumCommitIndex(val)
}

// DeleteFromEntry implements State.
func (this *raftStateImpl) DeleteFromEntry(entryIndex uint) {
	this.log.DeleteFromEntry(entryIndex)
}

// GetNumberNodesInCurrentConf implements State.
func (this *raftStateImpl) GetNumberNodesInCurrentConf() int {
	return this.log.GetNumberNodesInCurrentConf()
}

// IsInConf implements State.
func (this *raftStateImpl) IsInConf(nodeIp string) bool {
	return this.log.IsInConf(nodeIp)
}

// ConfStatus implements State.
func (this *raftStateImpl) ConfChanged() bool {
	return this.log.ConfChanged()
}

// GetConfig implements State.
func (this *raftStateImpl) GetConfig() []string {
	return this.log.GetConfig()
}

func (this *raftStateImpl) GetIdPrivate() string {
	return this.myIp.private
}

func (this *raftStateImpl) GetIdPublic() string {
	return this.myIp.private
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
	this.role = newRole
}

func (this *raftStateImpl) GetEntries() []l.LogInstance {
	return this.log.GetEntries()
}

func (this *raftStateImpl) AppendEntries(newEntries []*l.LogInstance) {
	var numNodeInConf = this.GetNumberNodesInCurrentConf()

	this.log.AppendEntries(newEntries)
	//INFO: thir case of OR happens if there are two nodes in the cluster, the follower drops
	// and so the leader does not have to ask the dropped follower the ack the deletion of that
	// node
	if this.role == LEADER ||
		numNodeInConf == 1 ||
		(this.role == LEADER && numNodeInConf == 2) {
		log.Println("auto commit entry: ", newEntries[0].Entry)
		for range newEntries {
			if this.role == LEADER {
				this.leaderMetadata.statePool.IncreaseCommonMathcIndex()
			}
			this.log.IncreaseCommitIndex()
		}
		return
	}
	log.Printf("leader, request to send log Entry to follower: ch %v, idx: %v\n", this.leaderMetadata.leaderEntryToCommit, this.log.GetCommitIndex()+1)
	for range newEntries {
		this.leaderMetadata.leaderEntryToCommit <- this.log.GetCommitIndex() + 1
	}
}

func (this *raftStateImpl) GetLeaderEntryChannel() *chan int64 {
	return &this.leaderMetadata.leaderEntryToCommit
}

func (this *raftStateImpl) GetCommitIndex() int64 {
	return this.log.GetCommitIndex()
}

// LastLogIndex implements State.
func (this *raftStateImpl) LastLogIndex() int {
	return this.log.LastLogIndex()
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
    s.leaderMetadata.leaderIp.private= ""
	s.leaderMetadata.statePool = nodematchidx.NewNodeCommonMatch()
	s.leaderMetadata.leaderEntryToCommit = make(chan int64)

    s.log = l.NewLogEntry(fsRootDir)

	s.TimeoutPool = timeout.NewTimeoutPool()
	s.TimeoutPool.AddTimeout(TIMER_ELECTION, time.Duration(randelection))
	s.TimeoutPool.AddTimeout(TIMER_HEARTHBIT, time.Duration(H_TIMEOUT))
    s.TimeoutPool.StopTimeout(TIMER_ELECTION)
    s.TimeoutPool.RestartTimeout(TIMER_HEARTHBIT)

	go s.leaaderUpdateCommitIndex()

	return s
}
