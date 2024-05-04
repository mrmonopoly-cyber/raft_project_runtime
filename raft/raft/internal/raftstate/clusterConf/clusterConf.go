package clusterconf

import "sync"

type Configuration interface{
    GetConfig() []string
    UpdateConfiguration(nodeIps []string)
    CommitConfig()
    ConfStatus() bool
    IsInConf(nodeIp string) bool
}

func NewConf(baseConf []string) Configuration{
    var newConf = make([]string,0)
    return &conf{
        lock: sync.RWMutex{},
        oldConf: &baseConf,
        newConf: &newConf,
        changed: false,
        joinConf: false,
    }
}
