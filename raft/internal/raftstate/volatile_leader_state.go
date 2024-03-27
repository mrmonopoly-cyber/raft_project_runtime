package raftstate

type VolatileLeaderState struct {
  nextIndex map[string]int
  matchIndex map[string]int
}

func (this *VolatileLeaderState) InitNextIndex(serverList []string, lastLogIndex int) {
  var m map[string]int = make(map[string]int)
  for _, i := range serverList {
    m[i] = lastLogIndex
  }
  this.nextIndex = m
}

func (this *VolatileLeaderState) InitMatchIndex(serverList []string) {
  var m map[string]int = make(map[string]int)
  for _, i := range serverList {
    m[i] = -1
  }
  this.matchIndex = m
}

func (this *VolatileLeaderState) SetNextIndex(id string, index int) {
  this.nextIndex[id] = index
}

func (this *VolatileLeaderState) SetMatchIndex(id string, index int) {
  this.matchIndex[id] = index
}
