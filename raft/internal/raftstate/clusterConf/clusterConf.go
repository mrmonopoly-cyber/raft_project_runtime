package clusterconf

type Configuration interface{
    GetConfig() []string
    UpdateConfiguration(nodeIps []string)
    CommitConfig()
}

func NewConf(baseConf []string) Configuration{
    return conf{
        oldConf: baseConf,
        newConf: make([]string, 0),
        committed: true,
    }
}
