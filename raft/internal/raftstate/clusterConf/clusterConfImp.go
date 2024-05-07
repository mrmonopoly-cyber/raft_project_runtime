package clusterconf

import (
	"log"
	"sync"
    "maps"
)

type conf struct {
	lock     sync.RWMutex
	oldConf  map[string]string
	newConf  map[string]string
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

    var resMap map[string]string = map[string]string{}
    var res []string = make([]string, len(this.oldConf))

	if this.joinConf {
        maps.Copy(resMap,this.oldConf)
	}
    maps.Copy(resMap,this.newConf)
    
    for _, v := range resMap {
        res = append(res, v)
    }

	return res
}

func (this *conf) UpdateConfiguration(op CONF_OPE, nodeIps []string) {
	log.Printf("Updating conf with new conf: %v\n", nodeIps)
	this.lock.Lock()
	defer this.lock.Unlock()

    var newConfBase = maps.Clone(this.oldConf)

    switch op{
    case ADD:
        for _, v := range nodeIps {
            newConfBase[v] = v
        }
        this.newConf = newConfBase
    case DEL:
        for _, v := range nodeIps {
            delete(newConfBase,v)
        }
    default:
        log.Println("invalid configuration operation, doing nothing, given: ", op)
        return
    }
    
	this.changed = true
	this.joinConf = true
}

func (this *conf) CommitConfig() {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.oldConf = this.newConf
	this.newConf = map[string]string{}
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

