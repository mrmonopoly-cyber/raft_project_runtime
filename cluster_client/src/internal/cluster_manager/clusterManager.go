package cluster_manager

type ClusterManager interface {
  Init()  
  Terminate() bool
  GetConfig()  
  SetConfig()
}

func NewClusterManager(initialConfig []string) ClusterManager {
  return &clusterManagerImpl{
    nodes: initialConfig,
    sources: *newSources(), 
  }
}
