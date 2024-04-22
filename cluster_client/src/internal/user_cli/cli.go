package usercli

import (
	m "raft/client/src/internal/user_cli/model"
  clusterform "raft/client/src/internal/user_cli/model/cluster_form"
)

type Cli interface {
  Start() (map[string]string, error)
}

type cli struct {}

func NewCli() *cli {
  return new(cli)
}

func (this *cli) Start() (map[string]string, error) {
  var text_input m.Model = clusterform.NewClusterForm("Read")
  return text_input.Show()
}



