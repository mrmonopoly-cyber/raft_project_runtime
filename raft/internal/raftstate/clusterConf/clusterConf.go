package clusterconf

type Configuration interface{
    GetConfig() []string
    UpdateConfiguration(nodeIps []string)
    CommitConfig()
    OverwriteConf(conf []string)
    ConfStatus() bool
}

func NewConf(baseConf []string) Configuration{
    var newConf = make([]string,0)
    return &conf{
        oldConf: &baseConf,
        newConf: &newConf,
        committed: true,
    }
}
