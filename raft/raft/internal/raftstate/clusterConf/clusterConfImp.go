package clusterconf

import (
	"log"
	"raft/pkg/raft-rpcProtobuf-messages/rpcEncoding/out/protobuf"
	"reflect"
	"sync"
)

type conf struct {
	lock         sync.RWMutex
	oldConf      map[string]string
	newConf      map[string]string
	joinConf     bool
	notifyChange chan int
}

// NotifyChangeInConf implements Configuration.
func (this *conf) NotifyChangeInConf() <-chan int {
	return this.notifyChange
}

// GetNumberNodesInCurrentConf implements Configuration.
func (this *conf) GetNumberNodesInCurrentConf() int {
	var conf []string = this.GetConfig()
	log.Printf("num of nodes in conf :%v, conf: %v\n", len(conf), conf)
	return len(conf)
}

func (this *conf) GetConfig() []string {
	this.lock.RLock()
	defer this.lock.RUnlock()

	var resMap map[string]string = map[string]string{}
	var res []string = nil

	for _, v := range this.oldConf {
		resMap[v] = v
	}
	if this.joinConf {
		for _, v := range this.newConf {
			resMap[v] = v
		}
	}

	for _, v := range resMap {
		if v != "" {
			res = append(res, v)
		}
	}

	return res
}

func (this *conf) UpdateConfiguration(op protobuf.Operation, nodeIps []string) {
	log.Printf("Updating conf with new nodes: %v\n", nodeIps)
	this.lock.Lock()
	defer this.lock.Unlock()

	switch op {
	case protobuf.Operation_JOIN_CONF_ADD:
		for _, v := range nodeIps {
			this.newConf[v] = v
		}
		this.joinConf = true
	case protobuf.Operation_JOIN_CONF_DEL:
		for _, v := range nodeIps {
			delete(this.newConf, v)
		}
	case protobuf.Operation_COMMIT_CONFIG_ADD:
		for _, v := range nodeIps {
			this.oldConf[v] = v
		}
		if reflect.DeepEqual(this.oldConf, this.newConf) {
			this.joinConf = false
		}

	case protobuf.Operation_COMMIT_CONFIG_REM:
		for _, v := range nodeIps {
			delete(this.oldConf, v)
		}
		if reflect.DeepEqual(this.oldConf, this.newConf) {
			this.joinConf = false
		}
	default:
		log.Println("invalid configuration operation, doing nothing, given: ", op)
		return
	}
    log.Println("new conf updated: ", this.GetConfig())

	this.notifyChange <- 1
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
