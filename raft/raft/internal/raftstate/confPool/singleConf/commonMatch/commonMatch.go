package commonmatch

import nodestate "raft/internal/raftstate/confPool/NodeIndexPool/nodeState"

type CommonMatch interface{
    CommitNewEntryC() <- chan int
    StopNotify()
}

func NewCommonMatch(initialCommonCommitIdx int, nodeStates []nodestate.NodeState) CommonMatch  {
    return NewCommonMatchImp(initialCommonCommitIdx, nodeStates)
}
