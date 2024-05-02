package nodeState

type VolatileNodeState interface{
    SetNextIndex(index int)
    SetMatchIndex(index int)
    GetMatchIndex() int
    GetNextIndex() int
    InitVolatileState(lastLogIndex int)
}

func NewVolatileState() VolatileNodeState{
    return &volatileNodeState{
        nextIndex: 0,
        matchIndex: 0,
    }
}

