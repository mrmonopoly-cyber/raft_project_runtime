package clusterconf

import (
	"log"
	"raft/pkg/raft-rpcProtobuf-messages/rpcEncoding/out/protobuf"
	"sync"
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
    var res []string = make([]string,0)

    for _,v  := range this.oldConf {
        resMap[v] = v
    }
	if this.joinConf {
        for _,v  := range this.newConf{
            resMap[v] = v
        }
	}
    
    for _, v := range resMap {
        res = append(res, v)
    }

	return res
}

func (this *conf) UpdateConfiguration(op protobuf.Operation, nodeIps []string) {
	log.Printf("Updating conf with new nodes: %v\n", nodeIps)
	this.lock.Lock()
	defer this.lock.Unlock()

    switch op{
    case protobuf.Operation_JOIN_CONF_ADD:
        for _, v := range nodeIps {
            this.newConf[v] = v
        }
        this.joinConf = true
    case protobuf.Operation_JOIN_CONF_DEL:
        for _, v := range nodeIps {
            delete(this.newConf,v)
            delete(this.oldConf,v)
        }
    case protobuf.Operation_COMMIT_CONFIG:
        for _, v := range nodeIps {
            delete(this.newConf,v)
            this.oldConf[v] = v
        }
        if len(this.newConf) == 0 {
            this.joinConf = false
        }
    default:
        log.Println("invalid configuration operation, doing nothing, given: ", op)
        return
    }
    
	this.changed = true
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



