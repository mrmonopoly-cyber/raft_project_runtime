package clusterconf

import (
	"log"
	"sync"
)

type conf struct {
	lock     sync.RWMutex
	oldConf  *[]string
	newConf  *[]string
	changed  bool
	joinConf bool
}

// GetNumberNodesInCurrentConf implements Configuration.
func (this *conf) GetNumberNodesInCurrentConf() int {
	var conf []string = this.GetConfig()
    log.Printf("num of nodes in conf :%v, conf: %v\n", len(conf), conf)
    return len(conf)
}

// ConfChanged implements Configuration.
func (this *conf) ConfChanged() bool {
	this.lock.RLock()
	defer this.lock.RUnlock()

	if this.changed {
		this.changed = false
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
	log.Printf("Updating conf with new conf: %v\n", nodeIps)
	this.lock.Lock()
	defer this.lock.Unlock()

	var newConf []string

	newConf = append(*this.newConf, nodeIps...)
	this.newConf = &newConf
	this.changed = true
	this.joinConf = true
}

func (this *conf) CommitConfig() {
	var newConf = make([]string, 0)

	this.lock.Lock()
	defer this.lock.Unlock()

	this.oldConf = this.newConf
	this.newConf = &newConf
	this.joinConf = false
}

func (this *conf) IsInConf(nodeIp string) bool {
	var currConf []string = this.GetConfig()
	for _, v := range currConf {
		if v == nodeIp {
			return true
		}
	}
	return false
}
