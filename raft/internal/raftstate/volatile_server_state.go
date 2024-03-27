package raftstate

type VolatileServerState struct {
  lastApplied  int
  commitIndex  uint64
}

func (this *VolatileServerState) UpdateLastApplied() int {
  if int(this.commitIndex) > this.lastApplied {
    this.lastApplied++
    return this.lastApplied
  }
  return -1
}

func (this *VolatileServerState) GetCommitIndex() uint64 {
  return this.commitIndex
}

func (this *VolatileServerState) SetCommitIndex(val uint64) {
  this.commitIndex = val
}

func (this *VolatileServerState) InitState() {
  this.commitIndex = 0
  this.lastApplied = 0
}

