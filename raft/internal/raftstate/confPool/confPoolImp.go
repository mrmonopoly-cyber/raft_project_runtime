package confpool

import (
	"errors"
	"log"
	"maps"
	genericmessage "raft/internal/genericMessage"
	localfs "raft/internal/localFs"
	"raft/internal/node"
	"raft/internal/raft_log"
	clustermetadata "raft/internal/raftstate/clusterMetadata"
	nodeIndexPool "raft/internal/raftstate/confPool/NodeIndexPool"
	nodestate "raft/internal/raftstate/confPool/NodeIndexPool/nodeState"
	singleconf "raft/internal/raftstate/confPool/singleConf"
	"raft/internal/rpcs/UpdateNode"
	"raft/pkg/raft-rpcProtobuf-messages/rpcEncoding/out/protobuf"
	"strings"
	"sync"

	"github.com/fatih/color"
)

type tuple struct {
	singleconf.SingleConf
	*raft_log.LogInstance
}

type confPool struct {
    lock         sync.Locker

	mainConf     singleconf.SingleConf
	newConf      singleconf.SingleConf
	emptyNewConf chan int
	nodeList     sync.Map
	numNodes     uint

    toUpdateNode map[string]string
    signalNewNode chan string


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

func (c* confPool) GetConfig() map[string]string{
    var mainConf = c.mainConf.GetConfig()
    if c.newConf != nil{
        maps.Copy(mainConf,c.newConf.GetConfig())
    }
    return mainConf
}

// GetNumNodesInConf implements SingleConf.
func (s *confPool) GetNumNodesInConf() uint{
    return s.numNodes
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
        go func(){
            c.signalNewNode <- node.GetIp()
        }()
	case REM:
		c.nodeList.Delete(node.GetIp())
		c.numNodes--
	}
}

func (c *confPool) AppendEntry(entry []*raft_log.LogInstance, prevLogIndex int) uint {

    c.lock.Lock()
    color.Yellow("appending entry, general pool: %v %v\n", entry, prevLogIndex)
    var appended = c.LogEntry.AppendEntry(entry,prevLogIndex)
    c.lock.Unlock()

    for i := 0; i < int(appended); i++ {
        color.Cyan("appending entry, general pool done\n")
        go func ()  {
            c.entryToCommiC <- 1
        }()
        
    }

    return appended
}

func (c *confPool) MinimumCommitIndex(val uint){
    var oldCommitIndex = c.GetCommitIndex()
    c.LogEntry.MinimumCommitIndex(val)
    var newCommitIndex = c.GetCommitIndex()
    var applyC = c.ApplyEntryC()
    for i := oldCommitIndex; i <= newCommitIndex; i++ {
        applyC <- 1
    }
}


func (c *confPool) IncreaseCommitIndex(){
    color.Red("you should not use this function, doing nothing\n")
}

func (c *confPool) appendEntryToConf(){
    for{
        <- c.entryToCommiC
        for c.GetCommitIndex() < int64(c.GetLogSize()-1){
            color.Cyan("appendEntryToConf: waiting signal: CI %state, LS %V\n",
                c.GetCommitIndex(),c.GetLogSize())
            var entry = c.GetEntriAt(c.GetCommitIndex()+1)


            if entry.Entry.OpType == protobuf.Operation_JOIN_CONF_FULL{
                <-c.emptyNewConf
                //INFO: critical case 1 of change in configuration.
                //updated all the new nodes until they are all pared
                var newConf = c.extractConfPayloadConf(entry.Entry)
                c.updateNewerNode(newConf)
                c.newConf = c.appendJoinConf(&newConf)
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

func (c *confPool) appendJoinConf(newConf *map[string]string) singleconf.SingleConf {
	return singleconf.NewSingleConf( *newConf,
        c.LogEntry,
		&c.nodeList,
		c.NodeIndexPool,
		c.commonMetadata)
}

func (c *confPool) extractConfPayloadConf(entry *protobuf.LogEntry) map[string]string{
	var confUnfiltered string = string(entry.Payload)
	var confFiltered []string = strings.Split(confUnfiltered, raft_log.SEPARATOR)
    var mainConf map[string]string = map[string]string{}

	confFiltered = confFiltered[0 : len(confFiltered)-1]
	for i := range confFiltered {
        var ip *string = &confFiltered[i]
		*ip, _ = strings.CutSuffix(*ip, " ")
		if *ip == "" || *ip == " " {
			continue
		}
        mainConf[*ip] = *ip
        c.UpdateStatusList(nodeIndexPool.ADD,*ip)
	}
    return mainConf
}

//INFO: send updated message to new node, it may be not present at the moment
//be aware
func (c *confPool) updateNewerNode(newConf map[string]string)  {
    var currConf = c.GetConfig()
    var changeVoteRight = UpdateNode.ChangeVoteRightNode(false)
    var rawMex,err = genericmessage.Encode(changeVoteRight)
    if err != nil{
        log.Panicln(err)
    }

    for _, v := range newConf {
        if currConf[v] == ""{
            c.NodeIndexPool.UpdateStatusList(nodeIndexPool.ADD, v)
            
            var val,f = c.nodeList.Load(v)
            var nNode node.Node
            if !f{
                color.Yellow("Node not yet present in the network: ",val)
                c.toUpdateNode[v] = v
                continue
            }
            nNode = val.(node.Node)
            nNode.Send(rawMex)
            go c.checkIfNodeIsUpdated(nNode)
        }
    }
}

func (c *confPool) checkIfNodeIsUpdated(nNode node.Node){
    var changeVoteRight = UpdateNode.ChangeVoteRightNode(true)
    var ip = nNode.GetIp()
    rawMex,err := genericmessage.Encode(changeVoteRight)
    state,err := c.FetchNodeInfo(ip)
    if err != nil{
        log.Panicln(err)
    }

    var _,subC = state.Subscribe(nodestate.MATCH)
    for{
        var match = <- subC
        if match >= int(c.GetCommitIndex()){
            break
        }
    }
    color.Green("node %v updated\n",ip)
    //TODO: unsubscribe
    nNode.Send(rawMex)
}

// daemon

//INFO: manage case of the node not yet present in the network
func (c *confPool) checkNodeToUpdate(){
    var changeVoteRight = UpdateNode.ChangeVoteRightNode(false)
    var rawMex,err = genericmessage.Encode(changeVoteRight)
    if err != nil{
        log.Panicln(err)
    }

    for {
        var ip = <- c.signalNewNode
        if c.toUpdateNode[ip] == ""{
            continue
        }
        var val,_ = c.nodeList.Load(ip)
        var nNode = val.(node.Node)
        nNode.Send(rawMex)
        delete(c.toUpdateNode,ip)
        go c.checkIfNodeIsUpdated(nNode)
    }
}

func (c *confPool) increaseCommitIndex() {
    color.Cyan("commit Index: waiting commit of main conf on ch: %v\n",c.mainConf.CommiEntryC())
    var activeC = <-c.mainConf.CommiEntryC()

    //TODO: to explain
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
		var entr = c.GetEntriAt(int64(toApplyIdx))

        color.Cyan("applying new entry: %v",entr)
        switch entr.Entry.OpType {
        case protobuf.Operation_COMMIT_CONFIG_ADD:
            color.Yellow("start applying commitADD:")
            var newConf = c.extractConfPayloadConf(entr.Entry)

            c.mainConf.CloseCommitEntryC()
            c.mainConf = c.appendJoinConf(&newConf)

            c.newConf = nil
            go func(){
                c.emptyNewConf <- 1
            }()

            if c.newConf != nil{
                color.Green("commit config applied [main,new]: ", 
                c.mainConf.GetConfig(), c.newConf.GetConfig())
            }else{
                color.Green("commit config applied [main,new]: ", 
                c.mainConf.GetConfig(), c.newConf)
            }
            color.Yellow("done applying commit:")
            
            //INFO: critical case 2 of change in configuration
            var newCommittedConf = c.mainConf.GetConfig()
            var myIp = c.commonMetadata.GetMyIp(clustermetadata.PRI)
            var found = false
            
            for _, v := range newCommittedConf {
                if v == myIp{
                    found = true
                    break
                }
            }

            if !found{
                color.Red("i'm no more in the conf, becomming follower")
                c.commonMetadata.SetRole(clustermetadata.FOLLOWER)
            }


        case protobuf.Operation_JOIN_CONF_FULL:
            if c.commonMetadata.GetRole() == clustermetadata.LEADER{
                var commit = protobuf.LogEntry{
                    Term:   c.commonMetadata.GetTerm(),
                    OpType: protobuf.Operation_COMMIT_CONFIG_ADD,
                    Payload: entr.Entry.Payload,
                }
                c.AppendEntry([]*raft_log.LogInstance{c.NewLogInstance(&commit)},-2)
            }
            color.Green("join conf applied")
        case protobuf.Operation_READ,protobuf.Operation_WRITE,protobuf.Operation_DELETE,
             protobuf.Operation_CREATE, protobuf.Operation_RENAME:

            var value = c.LocalFs.ApplyLogEntry(entr.Entry)
            entr.ReturnValue <- value
        default:
            log.Panicln("unrecognized opration: ",entr.Entry.OpType)
        }
	}
}

func confPoolImpl(rootDir string, commonMetadata clustermetadata.ClusterMetadata) *confPool {
	var res = &confPool{
        lock: &sync.Mutex{},

		mainConf:         nil,
		newConf:          nil,
		emptyNewConf:     make(chan int),
		nodeList:         sync.Map{},
        toUpdateNode:     map[string]string{},
        signalNewNode:    make(chan string),
		numNodes:         0,
		NodeIndexPool:    nodeIndexPool.NewLeaederCommonIdx(),
		commonMetadata:   commonMetadata,
        entryToCommiC:    make(chan int),
        LogEntry: raft_log.NewLogEntry(),
        LocalFs: localfs.NewFs(rootDir),
	}
	res.mainConf = singleconf.NewSingleConf(
		nil,
		res.LogEntry,
		&res.nodeList,
		res.NodeIndexPool,
		res.commonMetadata)

    go res.appendEntryToConf()
	go res.updateLastApplied()
    go res.checkNodeToUpdate()

    go func(){
        res.emptyNewConf <- 1
    }()

	return res
}
