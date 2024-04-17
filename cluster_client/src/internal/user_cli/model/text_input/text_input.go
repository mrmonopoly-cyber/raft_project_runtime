package textinput

import (
	"fmt"
	"os"
	m "raft/client/src/internal/user_cli/model"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)


const (
	hotPink  = lipgloss.Color("#FF06B7")
	darkGray = lipgloss.Color("#767676")
)

var (
	inputStyle    = lipgloss.NewStyle().Foreground(hotPink)
	continueStyle = lipgloss.NewStyle().Foreground(darkGray)
)

type (
	errMsg error
)

type TextInput struct {
  inputs  []textinput.Model
	operation string
  fileName  string
  fileNewName string
  focused int
	err error
}

func NewTextInput(operation string) m.Model { 
  var inputs []textinput.Model

  switch operation {
  case "Rename" :
      inputs = make([]textinput.Model, 2)
      inputs[0] = textinput.New()
      inputs[0].Placeholder = "Insert file name here..."
	    inputs[0].Focus()
	    inputs[0].CharLimit = 30
	    inputs[0].Width = 40
	    inputs[0].Prompt = ""

	    inputs[1] = textinput.New()
	    inputs[1].Placeholder = "Insert new name file here... "
 	    inputs[1].CharLimit = 30
	    inputs[1].Width = 40
	    inputs[1].Prompt = ""
  default :
      inputs = make([]textinput.Model, 1)
      inputs[0] = textinput.New()
      inputs[0].Placeholder = "Insert file name here..."
	    inputs[0].Focus()
	    inputs[0].CharLimit = 30
	    inputs[0].Width = 40
	    inputs[0].Prompt = ""
}


  var txtInput *TextInput = new(TextInput)
  txtInput.operation = operation
	txtInput.inputs = inputs
  txtInput.focused = 0
  txtInput.err = nil

  return txtInput
}

func (this TextInput) Init() tea.Cmd {
	return textinput.Blink 
}

func (this TextInput) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, len(this.inputs))
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if this.focused == len(this.inputs)-1 {
				return this, tea.Quit
			}
			this.nextInput()
		case tea.KeyCtrlC, tea.KeyEsc:
			return this, tea.Quit
		case tea.KeyShiftTab, tea.KeyCtrlP:
			this.prevInput()
		case tea.KeyTab, tea.KeyCtrlN:
			this.nextInput()
		}
		for i := range this.inputs {
			this.inputs[i].Blur()
		}
		this.inputs[this.focused].Focus()
    // We handle errors just like any other message
	case errMsg:
		this.err = msg
		return this, nil
	}

	for i := range this.inputs {
		this.inputs[i], cmds[i] = this.inputs[i].Update(msg)
	}
	return this, tea.Batch(cmds...)
}


func (this TextInput) View() string {

  var title string = "" 
  var subtitle string = ""
	if this.operation != "" {
    switch this.operation {
    case "Read":  
      title = "Name the file you want to read: "
    case "Write":  
      title = "Name the file you want to write: "
    case "Delete":  
      title = "Name the file you want to delete: "
    case "Create":  
      title = "Name of file: "
    case "Rename":  
      title = "Name the file you want to rename: "
      subtitle = "New name: "
  }
	}

  if subtitle == "" {
    return fmt.Sprintf(		  
      `  %s 
  %s
      `,
		  inputStyle.Width(50).Render(title),
		  this.inputs[0].View(),) + "\n"
  } else {
    return fmt.Sprintf(
		  `  %s  
  %s 

  %s 
  %s     
      `,
		  inputStyle.Width(50).Render(title),
		  this.inputs[0].View(),
      inputStyle.Width(50).Render(subtitle),
      this.inputs[1].View(), ) + "\n"
  }
}


func (this TextInput) Show() string {
  l, err := tea.NewProgram(this).Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

  txtInput := l.(TextInput)
  return txtInput.fileName

}

func (this *TextInput) nextInput() {
	this.focused = (this.focused + 1) % len(this.inputs)
}

func (this *TextInput) prevInput() {
	this.focused--
	// Wrap around
	if this.focused < 0 {
		this.focused = len(this.inputs) - 1
	}
}
