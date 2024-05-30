package raft_log

import (
	"errors"
	l "log"
	localfs "raft/internal/localFs"
	clusterconf "raft/internal/raftstate/clusterConf"
	"raft/pkg/raft-rpcProtobuf-messages/rpcEncoding/out/protobuf"
	p "raft/pkg/raft-rpcProtobuf-messages/rpcEncoding/out/protobuf"
	"strings"
	"sync"
)

type cConf interface {
	GetNumberNodesInCurrentConf() int
	IsInConf(ipNode string) bool
	ConfChanged() bool
	GetConfig() []string
}

type logInstance struct {
    entry *p.LogEntry
    notifyApplication chan int
}

type log struct {
	lock            sync.RWMutex
	newEntryToApply chan int

    entries     []logInstance
	logSize     uint
	commitIndex int64
	lastApplied int

	realClusterState
}


type realClusterState struct {
	cConf   clusterconf.Configuration
	localFs localfs.LocalFs
}

//log
// GetCommittedEntries implements LogEntry.
func (this *log) GetCommittedEntries() []*p.LogEntry{
    return this.getEntries(0)
}

// GetCommittedEntriesRange implements LogEntry.
func (this *log) GetCommittedEntriesRange(startIndex int) []*p.LogEntry{
    return this.getEntries(startIndex)
}

func (this *log) GetEntries() []*p.LogEntry{
    this.lock.RLock()
    defer this.lock.RUnlock()

    var lenEntries = len(this.entries)
    var res []*p.LogEntry = make([]*p.LogEntry, lenEntries)

    for i, v := range this.entries {
        res[i] = v.entry
    }


	return res
}

// GetEntriAt implements LogEntry.
func (this *log) GetEntriAt(index int64) (*p.LogEntry, error) {
	if (this.logSize == 1 && index == 0) || (index < int64(this.logSize)) {
		return this.entries[index].entry, nil
	}
	return nil, errors.New("invalid index: " + string(rune(index)))
}

func (this *log) AppendEntries(newEntries []*p.LogEntry) []chan int {
	this.lock.Lock()
	defer this.lock.Unlock()

    var lenNewEntries = len(newEntries)
	var lenEntries = len(this.entries)
    var notifyChann []chan int = nil
    var fullEntry []logInstance = make([]logInstance, lenNewEntries)

	for i, v := range fullEntry{
        v.entry = newEntries[i]
        v.notifyApplication = make(chan int)
        notifyChann = append(notifyChann, v.notifyApplication)

		if int(this.logSize) < lenEntries {
			this.entries[this.logSize] = v
		} else {
			this.entries = append(this.entries, v)
			lenEntries++
		}
		this.logSize++
	}

    return notifyChann
}

// DeleteFromEntry implements LogEntry.
func (this *log) DeleteFromEntry(entryIndex uint) {
    for i := int(entryIndex); i < len(this.entries); i++ {
        this.entries[i] = logInstance{
            entry: nil,
            notifyApplication: make(chan int),
        }
        this.logSize--
    }
}

func (this *log) GetCommitIndex() int64 {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.commitIndex
}

// IncreaseCommitIndex implements LogEntry.
func (this *log) IncreaseCommitIndex() {
	this.lock.Lock()
	defer this.lock.Unlock()
	if this.commitIndex < int64(this.logSize)-1 {
		this.commitIndex++
		this.newEntryToApply <- 1
	}
}

// MinimumCommitIndex implements LogEntry.
func (this *log) MinimumCommitIndex(val uint) {
    this.lock.Lock()
    defer this.lock.Unlock()

    if val < this.logSize {
        this.commitIndex = int64(val)
        return
    }
    this.commitIndex = int64(this.logSize) - 1
    if this.commitIndex > int64(this.lastApplied) {
        this.newEntryToApply <- 1
    }
}

func (this *log) LastLogIndex() int {
	this.lock.RLock()
	defer this.lock.RUnlock()

    var committedEntr = this.GetCommittedEntries()
    return len(committedEntr)-1
}

// LastLogTerm implements LogEntry.
func (this *log) LastLogTerm() uint {
    var committedEntr = this.GetCommittedEntries()
    var lasLogIdx = this.LastLogIndex()

    if lasLogIdx >= 0{
        return uint(committedEntr[lasLogIdx].Term)
    }
    return 0

}

//clusterConf

// GetNumberNodesInCurrentConf implements LogEntry.
func (this *log) GetNumberNodesInCurrentConf() int {
	return this.cConf.GetNumberNodesInCurrentConf()
}

// IsInConf implements LogEntry.
func (this *log) IsInConf(nodeIp string) bool {
	return this.cConf.IsInConf(nodeIp)
}

// CommitConfig implements LogEntry.
func (this *log) CommitConfig() {
	this.cConf.CommitConfig()
}

// ConfStatus implements LogEntry.
func (this *log) ConfChanged() bool {
	return this.cConf.ConfChanged()
}

// GetConfig implements LogEntry.
func (this *log) GetConfig() []string {
	return this.cConf.GetConfig()
}

// UpdateConfiguration implements LogEntry.
func (this *log) UpdateConfiguration(confOp protobuf.Operation, nodeIps []string) {
	this.cConf.UpdateConfiguration(confOp, nodeIps)
}

// utility

func (this *log) updateLastApplied() error {
	for {
		select {
		case <-this.newEntryToApply:
			this.lastApplied++
			var entry *logInstance = &this.entries[this.lastApplied]

			l.Printf("updating entry: %v", entry)
			switch entry.entry.OpType {
			case p.Operation_JOIN_CONF_ADD, p.Operation_JOIN_CONF_DEL, p.Operation_COMMIT_CONFIG:
				this.applyConf(entry.entry.OpType, entry)
			default:
				(*this).localFs.ApplyLogEntry(entry.entry)
			}
		}

	}
}

func (this *log) applyConf(ope protobuf.Operation, entry *logInstance) {
	var confUnfiltered string = string(entry.entry.Payload)
	var confFiltered []string = strings.Split(confUnfiltered, " ")
	l.Printf("applying the new conf:%v\t%v\n", confUnfiltered, confFiltered)
	this.cConf.UpdateConfiguration(ope, confFiltered)
    //HACK: if you are follower this goroutine remain stuck forever 
    //creating a zombie process
    go func ()  {   
        l.Println("notify change conf: ", entry)
        entry.notifyApplication <- 1 
    }()
}

func (this *log) getEntries(startIndex int) []*p.LogEntry{
	var committedEntries []*p.LogEntry = make([]*p.LogEntry , 0)

    if startIndex == -1{
        startIndex = 0
    }

    for i := int(startIndex); i <= int(this.commitIndex); i++ {
        committedEntries = append(committedEntries, this.entries[i].entry)
    }
	return committedEntries
}
