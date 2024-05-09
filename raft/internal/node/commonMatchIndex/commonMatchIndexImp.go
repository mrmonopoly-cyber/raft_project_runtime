package commonmatchindex

import "sync"

type commonMatchIndex struct {
	nodeNum     uint
	matchIndex  int
	updatedNode uint
	updateIndex chan int
	ipMap       sync.Map
}

// ResetCommonsMatchIndex implements CommonMatchIndex.
func (c *commonMatchIndex) ResetCommonsMatchIndex(matchIndex int) {
    c.matchIndex = matchIndex
    c.ipMap = sync.Map{}
    c.updatedNode = 0
    c.updateIndex = make(chan int)
}

// ChanUpdateNode implements CommonMatchIndex.
func (c *commonMatchIndex) ChanUpdateNode() chan int {
	return c.updateIndex
}

// ChangeNodeNum implements CommonMatchIndex.
func (c *commonMatchIndex) ChangeNodeNum(nodeNum uint) {
	c.nodeNum = nodeNum
}

// GetCommonMatchIndex implements CommonMatchIndex.
func (c *commonMatchIndex) GetCommonMatchIndex() int {
	return c.matchIndex
}

// GetNumUpdatedNode implements CommonMatchIndex.
func (c *commonMatchIndex) GetNumUpdatedNode() uint {
	return c.updatedNode
}

// IncreaseMatchIndex implements CommonMatchIndex.
func (c *commonMatchIndex) IncreaseUpdatedNode(ipNode string) {
	var _, found = c.ipMap.Load(ipNode)
	if found {
		return
	}
	c.updatedNode++
	c.ipMap.Store(ipNode, ipNode)
	if c.updatedNode > c.nodeNum/2 {
		c.updateIndex <- c.matchIndex
		c.updatedNode = 0
		c.ipMap = sync.Map{}
	}
}
