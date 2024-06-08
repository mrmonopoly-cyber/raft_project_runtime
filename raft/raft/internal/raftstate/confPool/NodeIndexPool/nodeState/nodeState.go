package nodestate

type INFO uint
const (
    MATCH INFO = iota
    NEXTT INFO = iota
)

type NodeState interface{
    UpdateNodeState(info INFO, val int)
    FetchData(info INFO) int
    Subscribe(info INFO) (int,<- chan int)
}

func NewNodeState(nodeIp string) NodeState{
    return newNodeStateImpl(nodeIp)
}
