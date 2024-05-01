package clusterconf

type conf struct{
    oldConf []string
    newConf []string
    committed bool
}

func (this conf) GetConfig() []string{
    return append(this.oldConf,this.newConf...)
}

func (this conf) UpdateConfiguration(nodeIps []string){
    this.newConf = append(this.newConf, nodeIps...)
    this.committed = false
}

func (this conf) CommitConfig(){
    this.oldConf = this.newConf
    this.newConf = make([]string, 0)
    this.committed = true
}
