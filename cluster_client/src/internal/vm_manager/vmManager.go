package vmmanager

import "raft/client/src/internal/cluster"

type VMManager interface {
  Init() cluster.Cluster  
  Terminate() bool
  AddNode() cluster.Cluster
  RemoveNode()
}

func NewClusterManager() VMManager {
  return &vMManagerImpl{
    sources: *newSources(), 
  }
}
