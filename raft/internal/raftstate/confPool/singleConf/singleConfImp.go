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
		entry.Committed <- 1
		return
	}

	//INFO:LEADER
	//Propagate to all nodes in this conf
	s.conf.Range(func(key, value any) bool {
		var v, f = s.nodeList.Load(key)
		var fNode node.Node
		var appendRpc rpcs.Rpc
		var enriesToSend []*protobuf.LogEntry
		var rawMex []byte
		var err error

		if !f {
			log.Println("node not yet connected or crashes or it's myself, skipping send: ", key)
			return true
		}
		fNode = v.(node.Node)

		switch entry.Entry.OpType {
		case protobuf.Operation_JOIN_CONF_DEL, protobuf.Operation_JOIN_CONF_ADD:
			enriesToSend = s.GetEntries()
			appendRpc = AppendEntryRpc.NewAppendEntryRPC(
				s.ClusterMetadata, s.LogEntry, -1, 0, enriesToSend)
		default:
			panic("not implemented")
		}

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

func (s *singleConfImp) updateEntryCommit() {
	for {
		<-s.CommonMatch.CommitNewEntryC()
		var committedEntry, err = s.LogEntry.GetEntriAt(s.GetCommitIndex() + 1)
		if err != nil {
			log.Panicln(err)
		}
		committedEntry.Committed <- 1
		s.commitC <- 1
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
		LogEntry:        raft_log.NewLogEntry(oldEntries),
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

	res.CommonMatch = commonmatch.NewCommonMatch(nodeStates)

	go res.updateEntryCommit()

	return res

}
