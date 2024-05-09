package commonmatchindex

import (
	"log"
	"sync"
)

type commonMatchIndex struct {
	nodeNum     uint
	matchIndex  int
	updatedNode uint
	updateIndex chan int
	ipMap       sync.Map
}

// DecreaseNodeNum implements CommonMatchIndex.
func (c *commonMatchIndex) DecreaseNodeNum() {
    c.nodeNum--
}

// IncreaseNodeNum implements CommonMatchIndex.
func (c *commonMatchIndex) IncreaseNodeNum() {
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
    log.Println("debug: ok loading: ", ipNode)
	c.updatedNode++
	c.ipMap.Store(ipNode, ipNode)
    log.Println("debug: ok storing: ", ipNode)
    if c.nodeNum>1 && c.updatedNode > c.nodeNum/2 { //HACK: to test and check
		c.updateIndex <- c.matchIndex
		c.updatedNode = 0
		c.ipMap = sync.Map{}
	}
}
