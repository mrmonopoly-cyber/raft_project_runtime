package usercli

import (
     "raft/client/src/internal/user_cli/list"
)

type Cli interface {
  Start() 
}

type cli struct {}

func NewCli() *cli {
  return new(cli)
}

func (this *cli) Start() {
  var l *lister.List = new(lister.List)
  l.Show()
}



