package usercli

import (
	m "raft/client/src/user_cli/model"
  clusterform "raft/client/src/user_cli/model/cluster_form"
)

type Cli interface {
  Start() map[string]string
}

type cli struct {}

func NewCli() *cli {
  return new(cli)
}

func (this *cli) Start() map[string]string {
  var fileNames map[string]string
  var text_input m.Model = clusterform.NewClusterForm("Read")
  fileNames = text_input.Show()
  return fileNames
}



