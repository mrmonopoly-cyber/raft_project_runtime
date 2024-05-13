package raftstate

import (
	"log"
	l "raft/internal/raft_log"
	clusterconf "raft/internal/raftstate/clusterConf"
	nodematchidx "raft/internal/raftstate/nodeMatchIdx"
	p "raft/pkg/raft-rpcProtobuf-messages/rpcEncoding/out/protobuf"
	"time"
)

type Role int

type raftStateImpl struct {
	idPrivate       string
	idPublic        string
	leaderIdPrivate string
	leaderIdPublic  string

	term uint64
	role Role

	voteFor string
	voting  bool

	electionTimeout    *time.Timer
	heartbeatTimeout   *time.Timer
	electionTimeoutRaw int

	log l.LogEntry

	nSupporting    uint64
	nNotSupporting uint64
	nNodeInCluster uint64

	//LEADER
	leaderEntryToCommit chan int64
	statePool           nodematchidx.NodeCommonMatch
}

// GetStatePool implements State.
func (this *raftStateImpl) GetStatePool() nodematchidx.NodeCommonMatch {
    return this.statePool
}

// GetEntriAt implements State.
func (this *raftStateImpl) GetEntriAt(index int64) (*p.LogEntry, error) {
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

// CommitConfig implements State.
func (this *raftStateImpl) CommitConfig() {
	(*this).log.CommitConfig()
}

// ConfStatus implements State.
func (this *raftStateImpl) ConfChanged() bool {
	return this.log.ConfChanged()
}

// GetConfig implements State.
func (this *raftStateImpl) GetConfig() []string {
	return this.log.GetConfig()
}

// UpdateConfiguration implements State.
func (this *raftStateImpl) UpdateConfiguration(confOp clusterconf.CONF_OPE, nodeIps []string) {
	this.log.UpdateConfiguration(confOp, nodeIps)
}

func (this *raftStateImpl) GetIdPrivate() string {
	return this.idPrivate
}

func (this *raftStateImpl) GetIdPublic() string {
	return this.idPublic
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

func (this *raftStateImpl) GetEntries() []*p.LogEntry {
	return this.log.GetEntries()
}

func (this *raftStateImpl) AppendEntries(newEntries []*p.LogEntry) {
	this.log.AppendEntries(newEntries)
	if !this.Leader(){
        log.Println("not leader increasing commitIndex in AppendEntry (state)")
		this.log.IncreaseCommitIndex()
		return
	}
    log.Println("leader, request to send log Entry to follower")
	this.leaderEntryToCommit <- this.log.GetCommitIndex() + 1
}

func (this *raftStateImpl) GetLeaderEntryChannel() *chan int64 {
	return &this.leaderEntryToCommit
}

func (this *raftStateImpl) GetCommitIndex() int64 {
	return this.log.GetCommitIndex()
}

// LastLogIndex implements State.
func (this *raftStateImpl) LastLogIndex() int {
	return this.log.LastLogIndex()
}

// More_recent_log implements State.
func (this *raftStateImpl) More_recent_log(last_log_index int64, last_log_term uint64) bool {
	return this.log.More_recent_log(last_log_index, last_log_term)
}

func (this *raftStateImpl) StartElectionTimeout() {
	this.electionTimeout.Reset(time.Duration((*this).electionTimeoutRaw))
}

func (this *raftStateImpl) StopElectionTimeout() {
	this.electionTimeout.Stop()
}

func (this *raftStateImpl) StartHearthbeatTimeout() {
	this.heartbeatTimeout.Reset(H_TIMEOUT)
}

func (this *raftStateImpl) StopHearthbeatTimeout() {
	this.heartbeatTimeout.Stop()
}

func (this *raftStateImpl) Leader() bool {
	return this.role == LEADER
}

func (this *raftStateImpl) CanVote() bool {
	return this.voting
}

func (this *raftStateImpl) VoteRight(vote bool) {
	(*this).voting = vote
}

func (this *raftStateImpl) HeartbeatTimeout() *time.Timer {
	return this.heartbeatTimeout
}

func (this *raftStateImpl) ElectionTimeout() *time.Timer {
	return this.electionTimeout
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

// MoreRecentLog implements State.
func (this *raftStateImpl) MoreRecentLog(lastLogIndex int64, lastLogTerm uint64) bool {
	return this.log.More_recent_log(lastLogIndex, lastLogTerm)
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

// GetLeaderIpPrivate implements State.
func (this *raftStateImpl) GetLeaderIpPrivate() string {
	return (*this).leaderIdPrivate
}

// GetLeaderIpPublic implements State.
func (this *raftStateImpl) GetLeaderIpPublic() string {
	return (*this).leaderIdPublic
}

// SetLeaderIpPublic implements State.
func (this *raftStateImpl) SetLeaderIpPublic(ip string) {
	(*this).leaderIdPublic = ip
}

// SetLeaderIpPrivate implements State.
func (this *raftStateImpl) SetLeaderIpPrivate(ip string) {
	(*this).leaderIdPrivate = ip
}

func (this *raftStateImpl) LeaaderUpdateCommitIndex(){
    for{
        select{
        case <- this.statePool.GetNotifyChannel():
            this.IncreaseCommitIndex()
        }
    }
}
