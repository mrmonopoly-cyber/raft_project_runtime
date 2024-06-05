package commonmatch

import nodestate "raft/internal/raftstate/confPool/NodeIndexPool/nodeState"

type CommonMatch interface{
    CommitNewEntryC() <- chan int
}

func NewCommonMatch(nodeStates []nodestate.NodeState) CommonMatch  {
    return NewCommonMatchImp(nodeStates)
}
