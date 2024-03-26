package raftstate

type VolatileLeaderState struct {
  nextIndex map[string]uint64
  matchIndex map[string]uint64
}

func (this *VolatileLeaderState) InitNextIndex(serverList []string, lastLogIndex uint64) {
  var m map[string]uint64 = make(map[string]uint64)
  for _, i := range serverList {
    m[i] = lastLogIndex
  }
  this.nextIndex = m
}

func (this *VolatileLeaderState) InitMatchIndex(serverList []string) {
  var m map[string]uint64 = make(map[string]uint64)
  for _, i := range serverList {
    m[i] = 0
  }
  this.matchIndex = m
}

func (this *VolatileLeaderState) DecrementNextIndex(id string, index uint64) {
  this.matchIndex[id] = index
}
