package singleconf

import (
	"log"
	genericmessage "raft/internal/genericMessage"
	"raft/internal/node"
	"raft/internal/raft_log"
	clustermetadata "raft/internal/raftstate/clusterMetadata"
	nodeIndexPool "raft/internal/raftstate/confPool/NodeIndexPool"
	"raft/internal/rpcs"
	"raft/internal/rpcs/AppendEntryRpc"
	"raft/pkg/raft-rpcProtobuf-messages/rpcEncoding/out/protobuf"
	"sync"
)

type singleConfImp struct {
	nodeList *sync.Map
	conf     sync.Map
	numNodes uint
	raft_log.LogEntry
    nodeIndexPool.NodeIndexPool
    clustermetadata.ClusterMetadata
}

func (s *singleConfImp) AppendEntry(entry *raft_log.LogInstance) {
    s.LogEntry.AppendEntry(entry)
    if  s.GetRole() == clustermetadata.FOLLOWER ||
        s.numNodes <= 1{ 
        //INFO: FOLLOWER or THE ONLY NODE IN THE CONF
        s.LogEntry.IncreaseCommitIndex()
        return
    }

    //INFO:LEADER
    //Propagate to all nodes in this conf
    s.conf.Range(func(key, value any) bool {
        var v,f = s.nodeList.Load(key)
        var fNode node.Node
        var appendRpc rpcs.Rpc
        var enriesToSend []*protobuf.LogEntry
        var rawMex []byte
        var err error

        if !f{
            log.Println("node not yet connected or crashes or it's myself, skipping send")
            return false
        }
        fNode = v.(node.Node)
        
        switch entry.Entry.OpType{
        case protobuf.Operation_JOIN_CONF_DEL, protobuf.Operation_JOIN_CONF_ADD:
            enriesToSend = s.GetEntries()
            appendRpc = AppendEntryRpc.NewAppendEntryRPC(s.ClusterMetadata,s.LogEntry,-1,0,enriesToSend)
        default:
            panic("not implemented")
        }
            
        rawMex,err = genericmessage.Encode(appendRpc)
        if err != nil{
            log.Panicln("error encoding: ", appendRpc, err)
        }

        err = fNode.Send(rawMex)
        if err != nil{
            log.Panicln("error sending rpc to: ",appendRpc,err)
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

func newSingleConfImp(  fsRootDir string, 
                        conf []string, 
                        oldEntries []*protobuf.LogEntry,
                        nodeList *sync.Map,
                        commonStatePool nodeIndexPool.NodeIndexPool,
                        commonMetadata clustermetadata.ClusterMetadata) *singleConfImp{
    var res = &singleConfImp{
        nodeList: nodeList,
        conf: sync.Map{},
        numNodes: 0,
        LogEntry: raft_log.NewLogEntry(fsRootDir,oldEntries),
        NodeIndexPool: commonStatePool,
        ClusterMetadata: commonMetadata,
    }

    for _, v := range conf {
        res.conf.Store(v,v)
        res.numNodes++
    }
    return res
    
}
