package confpool

import (
	"errors"
	"log"
	"raft/internal/node"
	"raft/internal/raft_log"
	clustermetadata "raft/internal/raftstate/clusterMetadata"
	nodeIndexPool "raft/internal/raftstate/confPool/NodeIndexPool"
	"raft/internal/raftstate/confPool/queue"
	singleconf "raft/internal/raftstate/confPool/singleConf"
	"raft/pkg/raft-rpcProtobuf-messages/rpcEncoding/out/protobuf"
	"reflect"
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
	confQueue    queue.Queue[tuple]
	emptyNewConf chan int
	nodeList     sync.Map
	numNodes     uint
	nodeIndexPool.NodeIndexPool
    commonMetadata clustermetadata.ClusterMetadata
}

// GetNode implements ConfPool.
func (c *confPool) GetNode(ip string) (node.Node, error) {
	var v, f = c.nodeList.Load(ip)
	if !f {
		return nil, errors.New("node not found: " + ip)
	}
	return v.(node.Node), nil
}

// NewLogInstanceBatch implements ConfPool.
func (c *confPool) NewLogInstanceBatch(entry []*protobuf.LogEntry, post []func()) []*raft_log.LogInstance {
	var res []*raft_log.LogInstance = nil
	var postLen = len(post)

	for i, v := range entry {
		var inst = raft_log.LogInstance{
			Entry: v,
		}
		if i < postLen {
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
		Entry:        entry,
		Committed:    make(chan int),
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

	log.Println("appending entry, general pool: ", entry)
	switch entry.Entry.OpType {
	case protobuf.Operation_JOIN_CONF_ADD:
		var confUnfiltered string = string(entry.Entry.Payload)
        log.Printf("debug unfiltered conf: %vEND\n",confUnfiltered)
		var confFiltered []string = strings.Split(confUnfiltered, raft_log.SEPARATOR)
        for i := range confFiltered {
            var ip = strings.TrimSuffix(confFiltered[i]," ")
            if ip == "" || ip == " "{
                continue
            }
	        log.Println("adding node in new conf: ", ip)
			c.NodeIndexPool.UpdateStatusList(nodeIndexPool.ADD, ip)
		}
        if c.mainConf.GetConfig() != nil {
            confFiltered = append(confFiltered, c.mainConf.GetConfig()...)
        }
        log.Println("new conf: ", confFiltered,len(confFiltered))
		var newConf = singleconf.NewSingleConf(
			c.fsRootDir,
			confFiltered,
            c.mainConf.GetEntries(),
			&c.nodeList,
			c.NodeIndexPool,
            c.commonMetadata)
        if c.newConf != nil {
            log.Println("checking conf is the same: ", newConf.GetConfig(), c.newConf)
        }
		// log.Println("checking conf is the same: ", newConf.GetConfig(), c.newConf.GetConfig())
		//WARN: DANGEROUS
		if c.newConf == nil || !reflect.DeepEqual(c.newConf.GetConfig(), newConf.GetConfig()) {
			c.confQueue.Push(tuple{SingleConf: newConf, LogInstance: entry})
			return
		}
	case protobuf.Operation_JOIN_CONF_DEL:
		panic("not implemented")
	}

	log.Println("append entry main conf: ", entry)
	c.mainConf.AppendEntry(entry)
	if c.newConf != nil {
		log.Println("append entry new conf: ", entry)
		c.newConf.AppendEntry(entry)
		joinConf = true
	}

	go func() {
		log.Println("waiting main conf commit of entry: ", entry)
		<-entry.Committed
		if joinConf {
			log.Println("waiting new conf commit of entry: ", entry)
			<-entry.Committed
		}
		log.Println("entry committed: ", entry)
		switch {
		case entry.Entry.OpType == protobuf.Operation_COMMIT_CONFIG_ADD:
			c.mainConf = c.newConf
			c.newConf = nil
			c.emptyNewConf <- 1
			log.Println("commit config applied [main,new]: ", c.mainConf, c.newConf)
		}
		if entry.AtCompletion != nil {
			entry.AtCompletion()
		}
	}()

}

// daemon
func (c *confPool) joinNextConf() {
	go func() {
		c.emptyNewConf <- 1
	}()
	for {
		log.Println("waiting commiting of new conf")
		<-c.emptyNewConf
		log.Println("waiting new conf to join")
		<-c.confQueue.WaitEl()
		var co = c.confQueue.Pop()
		c.newConf = co.SingleConf
		c.AppendEntry(co.LogInstance)
	}
}

func confPoolImpl(rootDir string, commonMetadata clustermetadata.ClusterMetadata) *confPool {
	var res = &confPool{
		mainConf:      nil,
		newConf:       nil,
		confQueue:     queue.NewQueue[tuple](),
		emptyNewConf:  make(chan int),
		nodeList:      sync.Map{},
		numNodes:      0,
		fsRootDir:     rootDir,
		NodeIndexPool: nodeIndexPool.NewLeaederCommonIdx(),
        commonMetadata: commonMetadata,
	}
	res.mainConf = singleconf.NewSingleConf(
		rootDir,
		nil,
        nil,
		&res.nodeList,
		res.NodeIndexPool,
        res.commonMetadata)

	go res.joinNextConf()

	return res
}
