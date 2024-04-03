package nodeState

type VolatileNodeState interface{
    SetNextIndex(index int)
    SetMatchIndex(index int)
    GetMatchIndex() int
    GetNextIndex() int
    InitVolatileState()
}

type volatileNodeState struct {
  nextIndex int
  matchIndex int
}
func NewVolatileState() VolatileNodeState{
    return &volatileNodeState{
        nextIndex: 0,
        matchIndex: 0,
    }
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

func (this* volatileNodeState) InitVolatileState(){
    (*this).nextIndex = 0;
    (*this).matchIndex = 0;
}

