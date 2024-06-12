package cluster_manager

type ClusterManager interface {
  Init()  
  Terminate() bool
  GetConfig()  
  SetConfig()
}

func NewClusterManager() ClusterManager {
  return &clusterManagerImpl{
    sources: *newSources(), 
  }
}
