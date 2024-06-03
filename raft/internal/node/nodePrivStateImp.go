package node

import nodematchidx "raft/internal/raftstate/nodeMatchIdx"

type volPrivState struct {
	statepool nodematchidx.NodeCommonMatch
}

// GetMatchIndex implements Node.
func (this *node) GetMatchIndex() int {
    return this.getData(nodematchidx.MATCH)
}

// GetNextIndex implements Node.
func (this *node) GetNextIndex() int {
    return this.getData(nodematchidx.NEXT)
}

// InitVolatileState implements Node.
func (this *node) InitVolatileState(lastLogIndex int) {
    this.statepool.InitVolatileState(this.GetIp(),lastLogIndex)
}

// NodeUpdated implements Node.
func (this *node) NodeUpdated() {
    this.statepool.DoneUpdating(this.GetIp())
}

// SetMatchIndex implements Node.
func (this *node) SetMatchIndex(index int) error{
    return this.statepool.UpdateNodeState(this.GetIp(),nodematchidx.MATCH,index)
}

// SetNextIndex implements Node.
func (this *node) SetNextIndex(index int) error{
    return this.statepool.UpdateNodeState(this.GetIp(),nodematchidx.NEXT,index)
}

// Updated implements Node.
func (this *node) Updated() bool {
    return this.statepool.Updated(this.GetIp())
}

//utility
func (this *node) getData(dataType nodematchidx.INDEX) int{
    var ip = this.GetIp()
    var res,err = this.statepool.GetNodeIndex(ip,dataType)
    if err != nil {
        return -10
    }
    return res
}
