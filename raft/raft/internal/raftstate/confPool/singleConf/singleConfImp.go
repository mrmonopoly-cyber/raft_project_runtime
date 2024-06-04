package singleconf

import (
	"raft/internal/raft_log"
	"sync"
)

type singleConfImp struct {
	nodeList *sync.Map
	conf     sync.Map
	numNodes uint
	raft_log.LogEntry
}


// GetConfig implements SingleConf.
func (s *singleConfImp) GetConfig() []string {
	var res []string = nil

    s.conf.Range(func(key, value any) bool {
        res = append(res, value.(string))
        return true
    })

    return res
}
