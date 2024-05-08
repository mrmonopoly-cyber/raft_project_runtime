package nodeState

type VolatileNodeState interface{
    SetNextIndex(index int)
    SetMatchIndex(index int)
    GetMatchIndex() int
    GetNextIndex() int
    InitVolatileState(lastLogIndex int)
    NextIndexStep()
}

func NewVolatileState() VolatileNodeState{
    return &volatileNodeState{
        nextIndex: -1,
        matchIndex: -1,
    }
}

