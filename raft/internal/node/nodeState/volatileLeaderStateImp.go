package nodeState

type volatileNodeState struct {
  nextIndex int
  matchIndex int
}

func (this *volatileNodeState) SetNextIndex(index int) {
  this.nextIndex = index
}

func (this *volatileNodeState) SetMatchIndex(index int) {
  this.matchIndex = index
}

func (this *volatileNodeState) GetMatchIndex() int {
  return this.matchIndex
}

func (this *volatileNodeState) GetNextIndex() int {
  return this.nextIndex
}

func (this* volatileNodeState) InitVolatileState(lastLogIndex int){
    (*this).nextIndex = lastLogIndex;
    (*this).matchIndex = 0;
}

func (this* volatileNodeState) NextIndexStep(){
    (*this).nextIndex++
    (*this).matchIndex++
}
