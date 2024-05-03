package clusterconf

import "sync"

type conf struct {
    lock      sync.RWMutex
	oldConf   *[]string
	newConf   *[]string
	changed bool
    joinConf bool
}

// ConfStatus implements Configuration.
func (this *conf) ConfStatus() bool {
    this.lock.RLock()
    defer this.lock.RUnlock()
    if this.changed {
        this.changed = false
        return false 
    }
	return true
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
    var newConf []string

    this.lock.Lock()
    defer this.lock.Unlock()

    newConf = append(*this.newConf, nodeIps...)
	this.newConf = &newConf
	this.changed = true
    this.joinConf = true

}

func (this *conf) CommitConfig() {
    var newConf = make([]string,0)

    this.lock.Lock()
    defer this.lock.Unlock()

	this.oldConf = this.newConf
	this.newConf = &newConf
	this.changed = false
    this.joinConf = false
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
