package nodeState

type volatileNodeState struct {
	nextIndex  int
	matchIndex int
	updated    bool
}

// NodeUpdated implements VolatileNodeState.
func (this *volatileNodeState) NodeUpdated() {
    this.updated = true
}

// Updated implements VolatileNodeState.
func (this *volatileNodeState) Updated() bool {
	return this.updated
}

func (this *volatileNodeState) SetNextIndex(index int) error{
	this.nextIndex = index
    return nil
}

func (this *volatileNodeState) SetMatchIndex(index int) error{
	this.matchIndex = index
    return nil
}

func (this *volatileNodeState) GetMatchIndex() int {
	return this.matchIndex
}

func (this *volatileNodeState) GetNextIndex() int {
	return this.nextIndex
}

func (this *volatileNodeState) InitVolatileState(lastLogIndex int) {
	(*this).nextIndex = lastLogIndex
	(*this).matchIndex = 0
}

func (this *volatileNodeState) NextIndexStep() {
	(*this).nextIndex++
	(*this).matchIndex++
}
