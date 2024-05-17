package clusterconf

import (
	"raft/pkg/raft-rpcProtobuf-messages/rpcEncoding/out/protobuf"
	"sync"
)

type Configuration interface{
    GetConfig() []string
    UpdateConfiguration(op protobuf.Operation,nodeIps []string)
    CommitConfig()
    ConfChanged() bool
    IsInConf(nodeIp string) bool
    GetNumberNodesInCurrentConf() int
}

func NewConf(baseConf []string) Configuration{
    var baseConfMap map[string]string = map[string]string{}
    for _, v := range baseConf {
        baseConfMap[v]=v
    }

    return &conf{
        lock: sync.RWMutex{},
        oldConf: baseConfMap,
        newConf: map[string]string{},
        changed: false,
        joinConf: false,
    }
}
