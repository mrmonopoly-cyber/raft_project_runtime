package cluster

import (
	"raft/client/src/internal/rpcs"
	"raft/client/src/internal/utility"
)

type Cluster interface {
  GetConfig() rpcs.Rpc 
}

func NewCluster(IPs []utility.Pair[string,string]) *clusterImpl {
  return &clusterImpl{
    IPs: IPs,
    leaderIP: "",
    conn: nil,
  }
} 
