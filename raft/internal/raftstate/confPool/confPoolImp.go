package confpool

import (
	"errors"
	"log"
	"raft/internal/node"
	"raft/internal/raft_log"
	"raft/internal/raftstate/confPool/queue"
	singleconf "raft/internal/raftstate/confPool/singleConf"
	"raft/pkg/raft-rpcProtobuf-messages/rpcEncoding/out/protobuf"
	"strings"
	"sync"
)

type tuple struct {
	singleconf.SingleConf
	*raft_log.LogInstance
}

type confPool struct {
	fsRootDir    string
	mainConf     singleconf.SingleConf
	newConf      singleconf.SingleConf
	confQueue    *queue.QueueImp[tuple]
	emptyNewConf chan int
	nodeList     sync.Map
	numNodes     uint
}

// GetNode implements ConfPool.
func (c *confPool) GetNode(ip string) (node.Node,error) {
    var v,f = c.nodeList.Load(ip)
    if !f{
        return nil,errors.New("node not found: " + ip)
    }
    return v.(node.Node),nil
}

// NewLogInstanceBatch implements ConfPool.
func (c *confPool) NewLogInstanceBatch(entry []*protobuf.LogEntry, post []func()) []*raft_log.LogInstance {
    var res []*raft_log.LogInstance = nil
    var postLen = len(post)

    for i, v := range entry {
        var inst = raft_log.LogInstance{
            Entry: v,
        }
        if i < postLen{
            inst.AtCompletion = post[i]
        }
        
        res = append(res, &inst)
    }

    return res
}

// DeleteFromEntry implements ConfPool.
func (c *confPool) DeleteFromEntry(entryIndex uint) {
	panic("unimplemented")
}

// GetCommitIndex implements ConfPool.
func (c *confPool) GetCommitIndex() int64 {
	panic("unimplemented")
}

// GetCommittedEntries implements ConfPool.
func (c *confPool) GetCommittedEntries() []raft_log.LogInstance {
	panic("unimplemented")
}

// GetCommittedEntriesRange implements ConfPool.
func (c *confPool) GetCommittedEntriesRange(startIndex int) []raft_log.LogInstance {
	panic("unimplemented")
}

// GetEntriAt implements ConfPool.
func (c *confPool) GetEntriAt(index int64) (*raft_log.LogInstance, error) {
	panic("unimplemented")
}

// GetEntries implements ConfPool.
func (c *confPool) GetEntries() []raft_log.LogInstance {
	panic("unimplemented")
}

// IncreaseCommitIndex implements ConfPool.
func (c *confPool) IncreaseCommitIndex() {
	panic("unimplemented")
}

// LastLogIndex implements ConfPool.
func (c *confPool) LastLogIndex() int {
	panic("unimplemented")
}

// LastLogTerm implements ConfPool.
func (c *confPool) LastLogTerm() uint {
	panic("unimplemented")
}

// MinimumCommitIndex implements ConfPool.
func (c *confPool) MinimumCommitIndex(val uint) {
	panic("unimplemented")
}

// NewLogInstance implements ConfPool.
func (c *confPool) NewLogInstance(entry *protobuf.LogEntry, post func()) *raft_log.LogInstance {
    var res = &raft_log.LogInstance{
        Entry: entry,
        Committed: make(chan int),
        AtCompletion: post,
    }
    return res
}

func (c *confPool) GetNodeList() *sync.Map {
	return &c.nodeList
}

// GetConf implements ConfPool.
func (c *confPool) GetConf() []string {
	if c.mainConf.GetConfig() == nil {
		return nil
	}
	var mainConf = c.mainConf.GetConfig()
	if c.newConf != nil {
		mainConf = append(mainConf, c.newConf.GetConfig()...)
	}
	return mainConf
}

// UpdateNodeList implements ConfPool.
func (c *confPool) UpdateNodeList(op OP, node node.Node) {
	switch op {
	case ADD:
		c.nodeList.Store(node.GetIp(), node)
		c.numNodes++
	case REM:
		c.nodeList.Delete(node.GetIp())
		c.numNodes--
	}
}

func (c *confPool) AppendEntry(entry *raft_log.LogInstance) {
	var joinConf bool = false

    log.Println("appending entry, general pool: ",entry)
	switch entry.Entry.OpType {
	case protobuf.Operation_JOIN_CONF_ADD:
		var confUnfiltered string = string(entry.Entry.Payload)
		var confFiltered []string = strings.Split(confUnfiltered, raft_log.SEPARATOR)
		for i := range confFiltered {
			confFiltered[i] = strings.Trim(confFiltered[i], " ")
		}
		confFiltered = append(confFiltered, c.mainConf.GetConfig()...)
		var newConf = singleconf.NewSingleConf(c.fsRootDir, confFiltered, &c.nodeList)
		c.confQueue.Push(tuple{SingleConf: newConf, LogInstance: entry})
        log.Println("waiting conf pushed: ",newConf)
		for c.newConf == nil {
		} //HACK: POLLING WAIT
        log.Println("checking conf is the same: ",newConf, c.newConf)
		if c.newConf != newConf {
			return
		}
	case protobuf.Operation_JOIN_CONF_DEL:
		panic("not implemented")
	}

	c.mainConf.AppendEntry(entry)
	if c.newConf != nil {
		c.newConf.AppendEntry(entry)
		joinConf = true
	}

	if joinConf {
		<-entry.Committed
	}
	<-entry.Committed
	switch {
	case entry.Entry.OpType == protobuf.Operation_COMMIT_CONFIG_ADD:
		c.mainConf = c.newConf
		c.newConf = nil
		c.emptyNewConf <- 1
	case entry.AtCompletion != nil:
		entry.AtCompletion()
	}

}

func (c *confPool) joinNextConf() {
    c.emptyNewConf <- 1
	for {
		<-c.emptyNewConf
        log.Println("waiting on cc")
		<-c.confQueue.C
		var co = c.confQueue.Pop()
		c.newConf = co.SingleConf
		c.AppendEntry(co.LogInstance)
	}
}

func confPoolImpl(rootDir string) *confPool {
	var res = &confPool{
		mainConf:     nil,
		newConf:      nil,
		confQueue:    queue.NewQueueImp[tuple](),
		emptyNewConf: make(chan int),
		nodeList:     sync.Map{},
		numNodes:     0,
		fsRootDir:    rootDir,
	}
	res.mainConf = singleconf.NewSingleConf(rootDir, nil, &res.nodeList)

    go res.joinNextConf()

	return res
}
