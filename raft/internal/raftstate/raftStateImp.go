package raftstate

import (
	l "raft/internal/raft_log"
	p "raft/pkg/raft-rpcProtobuf-messages/rpcEncoding/out/protobuf"
	"time"
)

type Role int

type raftStateImpl struct {
	idPrivate          string
	idPublic           string
	serverList         []string
	term               uint64
	leaderIdPrivate    string
	leaderIdPublic     string
	role               Role
	voteFor            string
	voting             bool
	electionTimeout    *time.Timer
	heartbeatTimeout   *time.Timer
	log                l.LogEntry
	nSupporting        uint64
	nNotSupporting     uint64
	nNodeInCluster     uint64
	electionTimeoutRaw int
}

// ConfStatus implements State.
func (this *raftStateImpl) ConfStatus() bool {
    return this.ConfStatus()
}

// IsInConf implements State.
func (this *raftStateImpl) IsInConf(ip string) bool {
    return this.IsInConf(ip)
}

// GetClusterConfig implements State.
func (this *raftStateImpl) GetClusterConfig() []string {
	return this.log.GetClusterConfig()
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

func (this *raftStateImpl) ToggleVoteRight() {
	(*this).voting = !(*this).voting
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

// GetNumNodeInCluster implements State.
func (this *raftStateImpl) GetNumNodeInCluster() uint64 {
	return this.nNodeInCluster
}

// IncreaseNodeInCluster implements State.
func (this *raftStateImpl) IncreaseNodeInCluster() {
	this.nNodeInCluster++
}

// DecreaseNodeInCluster implements State.
func (this *raftStateImpl) DecreaseNodeInCluster() {
	this.nNodeInCluster--
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
