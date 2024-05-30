package nodeState

type VolatileNodeState interface{
    SetNextIndex(index int) error
    SetMatchIndex(index int) error
    GetMatchIndex() int
    GetNextIndex() int
    InitVolatileState(lastLogIndex int)
    Updated() bool
    NodeUpdated()
}

func NewVolatileState() VolatileNodeState{
    return &volatileNodeState{
        nextIndex: -1,
        matchIndex: -1,
        updated: false,
    }
}

