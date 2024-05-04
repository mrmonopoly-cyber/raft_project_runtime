package clusterconf

import "sync"

type conf struct {
    lock      sync.RWMutex
	oldConf   *[]string
	newConf   *[]string
	chenged bool
    joinConf bool
}

// ConfStatus implements Configuration.
func (this *conf) ConfStatus() bool {
    this.lock.RLock()
    defer this.lock.RUnlock()
    
    if this.chenged{
        this.chenged = false
        return true
    }
    return false
}

func (this *conf) GetConfig() []string {
    this.lock.RLock()
    defer this.lock.RUnlock()
    if !this.joinConf {
        return append(*this.oldConf, *this.newConf...)
    }
    return *this.oldConf
}

func (this *conf) UpdateConfiguration(nodeIps []string) {
    this.lock.Lock()
    defer this.lock.Unlock()

    var newConf []string

    newConf = append(*this.newConf, nodeIps...)
	this.newConf = &newConf
	this.chenged = true
    this.joinConf = true
}

func (this *conf) CommitConfig() {
    var newConf = make([]string,0)

    this.lock.Lock()
    defer this.lock.Unlock()

	this.oldConf = this.newConf
	this.newConf = &newConf
	this.joinConf= false
}

func (this *conf) IsInConf(nodeIp string) bool{
    var currConf []string = this.GetConfig()
    for _, v := range currConf {
        if v == nodeIp {
            return true
        }
    }
    return false
}
