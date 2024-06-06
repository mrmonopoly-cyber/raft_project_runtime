package confpool

import (
	"errors"
	"log"
	localfs "raft/internal/localFs"
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
		c.nodeList.Store(node.GetIp(), node)
		c.numNodes++
	case REM:
		c.nodeList.Delete(node.GetIp())
		c.numNodes--
	}
}

func (c *confPool) AppendEntry(entry *raft_log.LogInstance) {
	log.Println("appending entry, general pool: ", entry)
    c.LogEntry.AppendEntry(entry)
    c.entryToCommiC <- 1
}

func (c *confPool) appendEntryToConf(){
    for {
        <- c.entryToCommiC
        var newConf singleconf.SingleConf
        var entry,err = c.GetEntriAt(c.GetCommitIndex()+1)
        if err != nil{
            log.Panicln(err)
        }

        switch entry.Entry.OpType {
        case protobuf.Operation_JOIN_CONF_ADD:
            newConf = c.appendJoinConfADD(entry)
            if c.pushJoinConf(entry,newConf){
                return
            }
        case protobuf.Operation_JOIN_CONF_DEL:
            newConf = c.appendJoinConfDEL(entry)
            if c.pushJoinConf(entry,newConf){
                return
            }
        }

        log.Println("append entry main conf: ", entry)
        c.mainConf.AppendEntry(entry)
        if c.newConf != nil {
            log.Println("append entry new conf: ", entry)
            var entryCopy raft_log.LogInstance = raft_log.LogInstance{
                Entry:        entry.Entry,
                AtCompletion: entry.AtCompletion,
            }
            c.newConf.AppendEntry(&entryCopy)
        }
    }
}

func (c *confPool) pushJoinConf(entry *raft_log.LogInstance, newConf singleconf.SingleConf) bool{
	//WARN: DANGEROUS

	if c.newConf != nil {
		log.Println("checking conf is the same: ", newConf.GetConfig(), c.newConf.GetConfig())
	}
	if c.newConf == nil || !reflect.DeepEqual(c.newConf.GetConfig(), newConf.GetConfig()) {
		c.confQueue.Push(tuple{SingleConf: newConf, LogInstance: entry})
		return true
	}
    return false
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
		c.mainConf.GetEntries(),
		&c.nodeList,
		c.NodeIndexPool,
		c.commonMetadata)
	return newConf
}

// daemon
func (c *confPool) increaseCommitIndex() {
	for {
        log.Println("commit Index: waiting commit of main conf")
		<-c.mainConf.CommiEntryC()
		if c.newConf != nil {
            log.Println("commit Index: waiting commit of new conf")
			<-c.newConf.CommiEntryC()
		}
        c.IncreaseCommitIndex()
	}
}

func (c *confPool) updateLastApplied() {
	for {
        var toApplyIdx = <- c.ApplyEntryC()
		var entr, err = c.GetEntriAt(int64(toApplyIdx))
		if err != nil {
			log.Panicln(err)
		}

        switch entr.Entry.OpType {
        case protobuf.Operation_COMMIT_CONFIG_ADD:
            c.mainConf = c.newConf
            c.newConf = nil
            c.emptyNewConf <- 1
            log.Println("commit config applied [main,new]: ", c.mainConf, c.newConf)
        case protobuf.Operation_COMMIT_CONFIG_REM:
            panic("Not implemented")
        case protobuf.Operation_JOIN_CONF_ADD, protobuf.Operation_JOIN_CONF_DEL:
            if c.commonMetadata.GetRole() == clustermetadata.LEADER{
                var commit = protobuf.LogEntry{
                    Term:   c.commonMetadata.GetTerm(),
                    OpType: protobuf.Operation_COMMIT_CONFIG_ADD,
                }
                c.AppendEntry(c.NewLogInstance(&commit, nil))
            }
        case protobuf.Operation_READ,protobuf.Operation_WRITE,protobuf.Operation_DELETE,
             protobuf.Operation_CREATE, protobuf.Operation_RENAME:

            c.LocalFs.ApplyLogEntry(entr.Entry)
        default:
            log.Panicln("unrecognized opration: ",entr.Entry.OpType)
        }

        if entr.AtCompletion != nil {
            entr.AtCompletion()
        }

	}
}

func (c *confPool) joinNextConf() {
	for {
		log.Println("waiting commiting of new conf")
		<-c.emptyNewConf
		log.Println("waiting new conf to join")
		<-c.confQueue.WaitEl()
		var co = c.confQueue.Pop()
        log.Println("new conf to join: ",co.SingleConf.GetConfig())
		c.newConf = co.SingleConf
		c.AppendEntry(co.LogInstance)
	}
}

func confPoolImpl(rootDir string, commonMetadata clustermetadata.ClusterMetadata) *confPool {
	var res = &confPool{
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
		nil,
		&res.nodeList,
		res.NodeIndexPool,
		res.commonMetadata)

	res.mainConf = mainConf

	go res.joinNextConf()
	go res.increaseCommitIndex()
	go res.updateLastApplied()
    go res.appendEntryToConf()

    res.emptyNewConf <- 1

	return res
}
