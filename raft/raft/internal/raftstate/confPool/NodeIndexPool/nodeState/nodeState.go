package nodestate

type INFO uint
const (
    MATCH INFO = iota
    NEXTT INFO = iota
)

type NodeState interface{
    UpdateNodeState(info INFO, val int)
    FetchData(info INFO) int
}

func NewNodeState() NodeState{
    return newNodeStateImpl()
}
