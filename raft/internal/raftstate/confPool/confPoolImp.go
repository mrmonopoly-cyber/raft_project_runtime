package confpool

import (
	"errors"
    "github.com/fatih/color"
	"log"
	localfs "raft/internal/localFs"
	"raft/internal/node"
	"raft/internal/raft_log"
	clustermetadata "raft/internal/raftstate/clusterMetadata"
	nodeIndexPool "raft/internal/raftstate/confPool/NodeIndexPool"
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
    lock         sync.Locker

	fsRootDir    string
	mainConf     singleconf.SingleConf
	newConf      singleconf.SingleConf
	confQueue    queue.Queue[tuple]
	emptyNewConf chan int
	nodeList     sync.Map
	numNodes     uint

	nodeIndexPool.NodeIndexPool
	commonMetadata clustermetadata.ClusterMetadata

    entryToCommiC chan int
    raft_log.LogEntry
    localfs.LocalFs
}


// SendHearthBit implements ConfPool.
func (c *confPool) SendHearthBit() {
    c.mainConf.SendHearthbit()
    if c.newConf != nil{
        c.newConf.SendHearthbit()
    }
}

// GetNode implements ConfPool.
func (c *confPool) GetNode(ip string) (node.Node, error) {
	var v, f = c.nodeList.Load(ip)
	if !f {
		return nil, errors.New("node not found: " + ip)
	}
	return v.(node.Node), nil
}

func (c *confPool) GetNodeList() *sync.Map {
	return &c.nodeList
}

// UpdateNodeList implements ConfPool.
func (c *confPool) UpdateNodeList(op OP, node node.Node) {
	switch op {
	case ADD:
        log.Println("storing a new Node")
		c.nodeList.Store(node.GetIp(), node)
        c.UpdateStatusList(nodeIndexPool.ADD,node.GetIp())
		c.numNodes++
	case REM:
		c.nodeList.Delete(node.GetIp())
		c.numNodes--
	}
}

func (c *confPool) AppendEntry(entry []*raft_log.LogInstance, prevLogIndex int) uint {
    c.lock.Lock()
    defer c.lock.Unlock()

    var appended uint = 0

    color.Yellow("appending entry, general pool: %v %v\n", entry, prevLogIndex)
    for i, v := range entry {
        appended = c.LogEntry.AppendEntry([]*raft_log.LogInstance{v},prevLogIndex+i)
        if appended > 0{
            color.Cyan("appending entry, general pool done\n")
            go func ()  {
                c.entryToCommiC <- 1
            }()
        }
    }

    return appended
}

func (c *confPool) IncreaseCommitIndex(){
    <- c.entryToCommiC
}

func (c *confPool) appendEntryToConf(){
    for{
        <- c.entryToCommiC
        for c.GetCommitIndex() < int64(c.GetLogSize()){
            log.Println("appendEntryToConf: waiting signal")
            var entry = c.GetEntriAt(c.GetCommitIndex()+1)

            switch entry.Entry.OpType {
            case protobuf.Operation_JOIN_CONF_ADD:
                <-c.emptyNewConf
                c.newConf = c.appendJoinConfADD(entry)
            case protobuf.Operation_JOIN_CONF_DEL:
                <-c.emptyNewConf
                c.newConf = c.appendJoinConfDEL(entry)
            }

            log.Println("notifying main conf to commit a new entry")
            c.mainConf.NotifyAppendEntryC() <- 1
            if c.newConf != nil{
                log.Println("notifying new conf to commit a new entry")
                c.newConf.NotifyAppendEntryC() <- 1
            }
            c.increaseCommitIndex()
        }
    }
}

func (c *confPool) appendJoinConfDEL(entry *raft_log.LogInstance) singleconf.SingleConf {
	panic("not implemented")
}

func (c *confPool) appendJoinConfADD(entry *raft_log.LogInstance) singleconf.SingleConf {
	var confUnfiltered string = string(entry.Entry.Payload)
	var confFiltered []string = strings.Split(confUnfiltered, raft_log.SEPARATOR)
    var mainConf = c.mainConf.GetConfig()

	confFiltered = confFiltered[0 : len(confFiltered)-1]
	for i := range confFiltered {
		confFiltered[i], _ = strings.CutSuffix(confFiltered[i], " ")
		var ip *string = &confFiltered[i]
		if *ip == "" || *ip == " " {
			continue
		}
		c.NodeIndexPool.UpdateStatusList(nodeIndexPool.ADD, *ip)
	}

    for _,v := range mainConf {
        confFiltered = append(confFiltered, v)
    }

	var newConf = singleconf.NewSingleConf(
		confFiltered,
        c.LogEntry,
		&c.nodeList,
		c.NodeIndexPool,
		c.commonMetadata)
	return newConf
}

// daemon
func (c *confPool) increaseCommitIndex() {
    color.Cyan("commit Index: waiting commit of main conf on ch: %v\n",c.mainConf.CommiEntryC())
    var activeC = <-c.mainConf.CommiEntryC()

    if activeC == 0{
        return
    }

    color.Cyan("main conf committed")
    if c.newConf != nil {
        color.Cyan("commit Index: waiting commit of new conf on ch: %v\n",c.newConf.CommiEntryC())
        <-c.newConf.CommiEntryC()
        color.Cyan("new conf committed")
    }
    color.Cyan("increasing commitIndex")
    c.LogEntry.IncreaseCommitIndex()
}

func (c *confPool) updateLastApplied() {
	for {
        color.Cyan("waiting to apply new entry")
        var toApplyIdx = <- c.ApplyEntryC()
        color.Red("ready to apply new entry")
		var entr = c.GetEntriAt(int64(toApplyIdx))

        color.Cyan("applying new entry: %v",entr)
        switch entr.Entry.OpType {
        case protobuf.Operation_COMMIT_CONFIG_ADD:
            color.Yellow("start applying commitADD:")
            c.lock.Lock()

            c.mainConf.CloseCommitEntryC()

            c.mainConf = c.newConf
            c.newConf = nil
            c.emptyNewConf <- 1
            color.Green("commit config applied [main,new]: ", c.mainConf, c.newConf)

            go c.increaseCommitIndex()

            c.lock.Unlock()
            color.Yellow("done applying commitADD:")
        case protobuf.Operation_COMMIT_CONFIG_REM:
            panic("Not implemented")
        case protobuf.Operation_JOIN_CONF_ADD, protobuf.Operation_JOIN_CONF_DEL:
            if c.commonMetadata.GetRole() == clustermetadata.LEADER{
                var commit = protobuf.LogEntry{
                    Term:   c.commonMetadata.GetTerm(),
                    OpType: protobuf.Operation_COMMIT_CONFIG_ADD,
                }
                c.AppendEntry([]*raft_log.LogInstance{c.NewLogInstance(&commit, nil)},-2)
            }
            color.Green("join conf applied")
        case protobuf.Operation_READ,protobuf.Operation_WRITE,protobuf.Operation_DELETE,
             protobuf.Operation_CREATE, protobuf.Operation_RENAME:

            c.LocalFs.ApplyLogEntry(entr.Entry)
        default:
            log.Panicln("unrecognized opration: ",entr.Entry.OpType)
        }

        color.Red("finish switch")

        if entr.AtCompletion != nil {
            entr.AtCompletion()
        }
	}
}

func confPoolImpl(rootDir string, commonMetadata clustermetadata.ClusterMetadata) *confPool {
	var res = &confPool{
        lock: &sync.Mutex{},

		mainConf:         nil,
		newConf:          nil,
		confQueue:        queue.NewQueue[tuple](),
		emptyNewConf:     make(chan int),
		nodeList:         sync.Map{},
		numNodes:         0,
		fsRootDir:        rootDir,
		NodeIndexPool:    nodeIndexPool.NewLeaederCommonIdx(),
		commonMetadata:   commonMetadata,
        entryToCommiC: make(chan int),
        LogEntry: raft_log.NewLogEntry(nil,true),
        LocalFs: localfs.NewFs(rootDir),
	}
	var mainConf = singleconf.NewSingleConf(
		nil,
		res.LogEntry,
		&res.nodeList,
		res.NodeIndexPool,
		res.commonMetadata)

	res.mainConf = mainConf

    go res.appendEntryToConf()
	go res.updateLastApplied()

    go func(){
        res.emptyNewConf <- 1
    }()

	return res
}
