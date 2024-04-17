package model

import (
 	   tea "github.com/charmbracelet/bubbletea"
) 

type Model interface {
  Show() string
  View() string
  Update(msg tea.Msg) (tea.Model, tea.Cmd)
}


