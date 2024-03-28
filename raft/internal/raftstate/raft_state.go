package raftstate

import (
	"math/rand"
	"log"
	l "raft/internal/raft_log"
	p "raft/pkg/rpcEncoding/out/protobuf"
	"time"
)

type Role int

const (
	LEADER Role = iota
	FOLLOWER
	CANDIDATE
)

const (
  MIN_ELECTION_TIMEOUT time.Duration = 10000000000
  MAX_ELECTION_TIMEOUT time.Duration = 15000000000
	H_TIMEOUT        time.Duration = 3000000000
)

type raftStateImpl struct {
	id               string
  serverList       []string
	term             uint64
	leaderId         string
	role             Role
	voteFor          string
	voting           bool
  leaderState      *VolatileLeaderState
	electionTimeout  *time.Timer
	heartbeatTimeout *time.Timer
	log              l.LogEntry
	nSupporting      uint64
	nNotSupporting   uint64
	nNodeInCluster   uint64
}


type State interface {
	GetId() string
	GetTerm() uint64
	GetRole() Role
	StartElectionTimeout()
	StartHearthbeatTimeout()
  StopElectionTimeout()
  StopHearthbeatTimeout()
	Leader() bool
	HeartbeatTimeout() *time.Timer
	ElectionTimeout() *time.Timer
	GetVoteFor() string
	IncrementTerm()
	VoteFor(id string)
	CanVote() bool
	GetEntries() []p.LogEntry
	GetCommitIndex() uint64
    SetCommitIndex(val uint64)
    SetRole(newRole Role)
    AppendEntries(newEntries []*p.LogEntry, index int)
    SetTerm(newTerm uint64)
	MoreRecentLog(lastLogIndex int64, lastLogTerm uint64) bool
	IncreaseSupporters()
	IncreaseNotSupporters()
	IncreaseNodeInCluster()
	GetNumSupporters() uint64
	GetNumNotSupporters() uint64
	GetNumNodeInCluster() uint64
	ResetElection()
  BecomeFollower()
  SetNextIndex(id string, index int)
  SetMatchIndex(id string, index int)
  InitVolatileLeaderState()
  GetLastLogIndex() int
  UpdateLastApplied() int
}

func (this *raftStateImpl) GetId() string {
	return this.id
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

func (this *raftStateImpl) GetEntries() []p.LogEntry{
	return this.log.GetEntries()
}

func (this *raftStateImpl) AppendEntries(newEntries []*p.LogEntry, index int) {
  this.log.AppendEntries(newEntries, index)
}

func (this *raftStateImpl) GetCommitIndex() uint64 {
	return this.log.GetCommitIndex()
}

func (this *raftStateImpl) SetCommitIndex(val uint64) {
  this.log.SetCommitIndex(val)
}

func (this *raftStateImpl) StartElectionTimeout() {
  rand.New(rand.NewSource(time.Now().UnixNano()))
  t := rand.Intn((int(MAX_ELECTION_TIMEOUT)-int(MIN_ELECTION_TIMEOUT) + 1) + int(MIN_ELECTION_TIMEOUT))
	this.electionTimeout.Reset(time.Duration(t))
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

func (this *raftStateImpl) BecomeFollower() {
  this.role = FOLLOWER
  this.leaderState = nil
}

func (this *raftStateImpl) CanVote() bool {
	return this.voting
}

func (this *raftStateImpl) HeartbeatTimeout() *time.Timer {
    log.Println("timeout hearthbit")
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
func (this *raftStateImpl) MoreRecentLog(lastLogIndex uint64, lastLogTerm uint64) bool {
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

func (this *raftStateImpl) ResetElection() {
	this.nSupporting = 0
	this.nNotSupporting = 0
}

func (this *raftStateImpl) InitVolatileLeaderState() {
  this.leaderState = new(VolatileLeaderState)
  this.leaderState.InitMatchIndex(this.serverList)
  this.leaderState.InitNextIndex(this.serverList, this.log.LastLogIndex())
}

func (this *raftStateImpl) InitVolatileServer() {
  this.log.InitState()
}

func (this *raftStateImpl) SetNextIndex(id string, index int) {
  this.leaderState.SetNextIndex(id, index)
}

func (this *raftStateImpl) SetMatchIndex(id string, index int) {
  this.leaderState.SetMatchIndex(id, index)
}

func (this *raftStateImpl) GetLastLogIndex() int {
  return this.log.LastLogIndex()
}

func (this *raftStateImpl) UpdateLastApplied() int {
  return this.log.UpdateLastApplied()
  // TODO: apply log to state machine (?)
}

func NewState(term uint64, id string, role Role) State {
	var s = new(raftStateImpl)
	s.role = role
	s.term = term
  s.id = id
	s.electionTimeout = time.NewTimer(MAX_ELECTION_TIMEOUT)
	s.heartbeatTimeout = time.NewTimer(H_TIMEOUT)
	s.nNotSupporting = 0
	s.nSupporting = 0
	s.nNodeInCluster = 1
    s.voting = true
  s.InitVolatileServer()
	return s
}
