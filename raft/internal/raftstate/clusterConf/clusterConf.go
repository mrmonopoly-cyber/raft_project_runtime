package clusterconf

type Configuration interface{
    GetConfig() []string
    UpdateConfiguration(nodeIps []string)
    CommitConfig()
}

func NewConf(baseConf []string) Configuration{
    return conf{
        oldConf: baseConf,
        newCong: make([]string, 0),
        committed: true,
    }
}
