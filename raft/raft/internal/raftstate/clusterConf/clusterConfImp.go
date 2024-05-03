package clusterconf

type conf struct {
	oldConf   *[]string
	newConf   *[]string
	committed bool
}

// ConfStatus implements Configuration.
func (this conf) ConfStatus() bool {
	return this.committed
}

// OverwriteConf implements Configuration.
func (this conf) OverwriteConf(conf []string) {
	this.oldConf = &conf
    this.newConf = &conf
	this.committed = false
}

func (this conf) GetConfig() []string {
	return append(*this.oldConf, *this.newConf...)
}

func (this conf) UpdateConfiguration(nodeIps []string) {
	*this.newConf = append(*this.newConf, nodeIps...)
	this.committed = false
}

func (this conf) CommitConfig() {
    var newConf = make([]string,0)

	this.oldConf = this.newConf
	this.newConf = &newConf
	this.committed = true
}
