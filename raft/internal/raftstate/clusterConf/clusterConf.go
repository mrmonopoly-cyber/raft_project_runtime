package clusterconf

import "sync"

type CONF_OPE uint8

const (
    ADD CONF_OPE = iota
    DEL CONF_OPE = iota
)

type Configuration interface{
    GetConfig() []string
    UpdateConfiguration(op CONF_OPE,nodeIps []string)
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
