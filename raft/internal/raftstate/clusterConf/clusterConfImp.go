package clusterconf

type conf struct{
    oldConf []string
    newCong []string
    committed bool
}

func (this conf) GetConfig() []string{
    return append(this.oldConf,this.newCong...)
}

func (this conf) UpdateConfiguration(nodeIps []string){
    this.newCong = append(this.newCong, nodeIps...)
    this.committed = false
}

func (this conf) CommitConfig(){
    this.oldConf = this.newCong
    this.newCong = make([]string, 0)
    this.committed = true
}
