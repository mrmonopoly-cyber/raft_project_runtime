package usercli

import (
	"raft/client/src/internal/user_cli/model"
	"raft/client/src/internal/user_cli/model/list"
	"raft/client/src/internal/user_cli/model/text_input"
)

type Cli interface {
  Start() 
}

type cli struct {}

func NewCli() *cli {
  return new(cli)
}

func (this *cli) Start() (string, string) {
  var list model.Model = lister.NewList()
  var res = list.Show()
  var text_input model.Model = textinput.NewTextInput(res)
  var res1 = text_input.Show()
  return res, res1
}



