package clusterconf

import "sync"

type conf struct {
    lock      sync.RWMutex
	oldConf   *[]string
	newConf   *[]string
	committed bool
}

// ConfStatus implements Configuration.
func (this *conf) ConfStatus() bool {
    this.lock.RLock()
    defer this.lock.RUnlock()
	return this.committed
}

func (this *conf) GetConfig() []string {
    this.lock.RLock()
    defer this.lock.RUnlock()
    if !this.ConfStatus() {
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
	this.committed = false
}

func (this *conf) CommitConfig() {
    var newConf = make([]string,0)

    this.lock.Lock()
    defer this.lock.Unlock()

	this.oldConf = this.newConf
	this.newConf = &newConf
	this.committed = true
}
