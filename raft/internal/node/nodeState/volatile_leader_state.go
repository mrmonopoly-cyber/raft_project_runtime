package nodeState

type VolatileNodeState interface{
    SetNextIndex(id string, index int)
    SetMatchIndex(id string, index int)
    GetMatchIndex(id string, index int)
    GetNextIndex(id string, index int)
    InitState()
}

type volatileNodeState struct {
  nextIndex int
  matchIndex int
}

func (this *volatileNodeState) SetNextIndex(id string, index int) {
  this.nextIndex = index
}

func (this *volatileNodeState) SetMatchIndex(id string, index int) {
  this.matchIndex = index
}

func (this *volatileNodeState) GetMatchIndex(id string, index int) {
  this.matchIndex = index
}

func (this *volatileNodeState) GetNextIndex(id string, index int) {
  this.nextIndex = index
}

func (this* volatileNodeState) InitState(){
    (*this).nextIndex = 0;
    (*this).matchIndex = 0;
}
