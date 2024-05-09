package commonmatchindex

import "sync"

type CommonMatchIndex interface{
    GetCommonMatchIndex() int
    GetNumUpdatedNode() uint
    IncreaseUpdatedNode(ipNode string)
    IncreaseNodeNum()
    DecreaseNodeNum()
    ChanUpdateNode() chan int
    ResetCommonsMatchIndex(matchIndex int)
}

func NewCommonMatchIndex(nodeNum uint) CommonMatchIndex{
    return &commonMatchIndex{
        nodeNum: nodeNum,
        updatedNode: 0,
        matchIndex: -1,
        updateIndex: make(chan int),
        ipMap: sync.Map{},
    }
}
