package singleconf

import (
	"log"
	genericmessage "raft/internal/genericMessage"
	"raft/internal/node"
	"raft/internal/raft_log"
	clustermetadata "raft/internal/raftstate/clusterMetadata"
	nodeIndexPool "raft/internal/raftstate/confPool/NodeIndexPool"
	nodestate "raft/internal/raftstate/confPool/NodeIndexPool/nodeState"
	commonmatch "raft/internal/raftstate/confPool/singleConf/commonMatch"
	"raft/internal/rpcs"
	"raft/internal/rpcs/AppendEntryRpc"
	"raft/pkg/raft-rpcProtobuf-messages/rpcEncoding/out/protobuf"
	"sync"
)

type singleConfImp struct {
	nodeList *sync.Map
	conf     sync.Map
	numNodes uint
	nodeIndexPool.NodeIndexPool
	clustermetadata.ClusterMetadata
	commonmatch.CommonMatch

	raft_log.LogEntry
	commitC chan int
}

func (s *singleConfImp) SendHearthbit(){
    s.conf.Range(func(key, value any) bool {
        var v,f = s.nodeList.Load(key)
        var nNode node.Node
        var nextIndex int
        var hearthbit rpcs.Rpc
        var rawMex []byte

        if !f{
            s.nodeNotFound(key)
            return true
        }
        nNode = v.(node.Node)
        var state,err = s.FetchNodeInfo(nNode.GetIp())
        if err != nil{
            log.Panicln(err)
        }

        nextIndex = state.FetchData(nodestate.NEXTT)

        log.Println("sending hearthbit to: ",nNode.GetIp())

        if nextIndex < len(s.GetEntries()){
            var entries = s.GetEntriesRange(nextIndex)
            prevEntr,err := s.GetEntriAt(int64(nextIndex)-1)
            if err != nil{
                log.Panicln(err)
            }
            hearthbit = AppendEntryRpc.NewAppendEntryRPC(
                s.ClusterMetadata,
                s.LogEntry,
                int64(nextIndex)-1,
                prevEntr.Entry.Term,
                entries)
        }else {
            hearthbit = AppendEntryRpc.GenerateHearthbeat(s.LogEntry,s.ClusterMetadata)
        }

        rawMex,err = genericmessage.Encode(hearthbit)
        if err != nil{
            log.Panicln(err)
        }

        err = nNode.Send(rawMex)
        if err != nil {
            log.Println(err)
        }
        


        return true
    })
}

// CommiEntryC implements SingleConf.
func (s *singleConfImp) CommiEntryC() <-chan int {
	return s.commitC
}

func (s *singleConfImp) AppendEntry(entry *raft_log.LogInstance) {
	s.LogEntry.AppendEntry(entry)
	if s.GetRole() == clustermetadata.FOLLOWER || s.numNodes <= 1 {
		//INFO: FOLLOWER or THE ONLY NODE IN THE CONF
        log.Println("auto commit")
		s.commitC <- 1
        s.IncreaseCommitIndex()
		return
	}

	//INFO:LEADER
	//Propagate to all nodes in this conf
    log.Println("propagate to all follower: ",entry.Entry)
	s.conf.Range(func(key, value any) bool {
		var v, f = s.nodeList.Load(key)
		var fNode node.Node
		var appendRpc rpcs.Rpc
		var enriesToSend []*protobuf.LogEntry
		var rawMex []byte
		var err error

		if !f {
            return s.nodeNotFound(key)
		}
		fNode = v.(node.Node)

		switch entry.Entry.OpType {
		case protobuf.Operation_JOIN_CONF_DEL, protobuf.Operation_JOIN_CONF_ADD:
			enriesToSend = s.GetEntries()
		default:
            enriesToSend = append(enriesToSend, entry.Entry)
		}

        appendRpc = AppendEntryRpc.NewAppendEntryRPC(
            s.ClusterMetadata, s.LogEntry, -1, 0, enriesToSend)

		rawMex, err = genericmessage.Encode(appendRpc)
		if err != nil {
			log.Panicln("error encoding: ", appendRpc, err)
		}

        log.Println("sending rpc to node: ",appendRpc.ToString(),fNode.GetIp())
		err = fNode.Send(rawMex)
		if err != nil {
			log.Panicln("error sending rpc to: ", appendRpc, err)
		}

		return true
	})
}

// GetConfig implements SingleConf.
func (s *singleConfImp) GetConfig() []string {
	var res []string = nil

	s.conf.Range(func(key, value any) bool {
		res = append(res, value.(string))
		return true
	})

	return res
}

//utility

func (s *singleConfImp) nodeNotFound(key any) bool {
    if key == s.ClusterMetadata.GetMyIp(clustermetadata.PRI){
        log.Println("skiping myself from propagation: ", key)
        return true
    }
    log.Println("node not yet connected or crashes, skipping send: ", key)
    return true
}

func (s *singleConfImp) updateEntryCommit() {
	for {
		<-s.CommonMatch.CommitNewEntryC()
        log.Println("new entry to commit")
		s.commitC <- int(s.GetCommitIndex()) + 1
        s.IncreaseCommitIndex()
		//TODO: every time the common match is updated commit an entry
	}
}

func newSingleConfImp(conf []string,
	oldEntries []*protobuf.LogEntry,
	nodeList *sync.Map,
	commonStatePool nodeIndexPool.NodeIndexPool,
	commonMetadata clustermetadata.ClusterMetadata) *singleConfImp {
	var res = &singleConfImp{
		nodeList:        nodeList,
		conf:            sync.Map{},
		numNodes:        0,
		NodeIndexPool:   commonStatePool,
		ClusterMetadata: commonMetadata,
		CommonMatch:     nil,
		LogEntry:        raft_log.NewLogEntry(oldEntries,false),
		commitC:         make(chan int),
	}
	var nodeStates []nodestate.NodeState = nil

	for _, v := range conf {
		res.conf.Store(v, v)
		var st, err = commonStatePool.FetchNodeInfo(v)
		if err != nil {
			log.Panicln("state for node not found: ", v)
		}
		nodeStates = append(nodeStates, st)
		res.numNodes++
	}

    log.Println("node to subs: ",nodeStates)
	res.CommonMatch = commonmatch.NewCommonMatch(nodeStates)

	go res.updateEntryCommit()

	return res

}
