package clusterconf

import (
	"raft/pkg/raft-rpcProtobuf-messages/rpcEncoding/out/protobuf"
	"sync"
)

type Configuration interface{
    GetConfig() []string
    UpdateConfiguration(op protobuf.Operation,nodeIps []string)
    IsInConf(nodeIp string) bool
    GetNumberNodesInCurrentConf() int
    NotifyChangeInConf() <- chan int
}

func NewConf() Configuration{
    return &conf{
        lock: sync.RWMutex{},
        oldConf: map[string]string{},
        newConf: map[string]string{},
        notifyRemove: make(chan int),
        joinConf: false,
    }
}
