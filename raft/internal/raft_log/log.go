package raft_log

import p "raft/pkg/rpcEncoding/out/protobuf"

type LogEntry interface{
    GetEntries() []*p.LogEntry
    GetCommitIndex() int64
    More_recent_log(last_log_index int64, last_log_term uint64) bool 
    SetCommitIndex(val int64)
    AppendEntries(newEntries []*p.LogEntry, index int)
    LastLogIndex() int
    UpdateLastApplied() int 
    InitState()
    /* testing */
    AppendDummyEntry(term uint64)
}

type log struct {
    entries     []p.LogEntry
    lastApplied  int
    commitIndex  int64
}

func NewLogEntry() LogEntry {
  var l = new(log)
  l.commitIndex = 0
  l.lastApplied = 0
  l.entries = make([]p.LogEntry, 0)

  return l
}

func (this *log) GetEntries() []*p.LogEntry{
  var e []*p.LogEntry = make([]*p.LogEntry, 0)
  for i, en := range this.entries {
    e[i] = &en
  }
  return e
}

func (this *log) LastLogIndex() int {
  return len(this.entries)-1
}

func (this *log) More_recent_log(last_log_index int64, last_log_term uint64) bool {
    if last_log_index >= this.commitIndex{
        var entries []*p.LogEntry = this.GetEntries()
        if len(entries) <= int(last_log_index) {
            return true
        }
        if last_log_term >= *(entries[last_log_index]).Term {
            return true
        }
    }
    return false

}

func (this *log) AppendEntries(newEntries []*p.LogEntry, index int) {
  for i, en := range newEntries {
    this.entries[index + i + 1] = *en
  }
}

func (this *log) UpdateLastApplied() int {
  if int(this.commitIndex) > this.lastApplied {
    this.lastApplied++
    return this.lastApplied
  }
  return -1
}

func (this *log) GetCommitIndex() int64 {
  return this.commitIndex
}

func (this *log) SetCommitIndex(val int64) {
  this.commitIndex = val
}

func (this *log) InitState() {
  this.commitIndex = 0
  this.lastApplied = 0
}

/* testing */
func (this *log) AppendDummyEntry(term uint64) {
  desc := "ciao"
  op := p.Operation_WRITE

  newEntry := &p.LogEntry{
     Description: &desc,
    Term: &term,
    OpType: &op, 
  }
  
  this.entries = append(this.entries, *newEntry)
}
