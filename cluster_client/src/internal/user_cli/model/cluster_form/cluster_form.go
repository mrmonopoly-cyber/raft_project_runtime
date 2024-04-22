package clusterform

import (
	"fmt"
	"io"
	"strings"
  m "raft/client/src/internal/user_cli/model"
	"github.com/charmbracelet/bubbles/list"
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
	titleStyle        = lipgloss.NewStyle().MarginLeft(0)
  listHeight = 10
	itemStyle         = lipgloss.NewStyle().PaddingLeft(2)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(0).PaddingBottom(0)
	quitTextStyle     = lipgloss.NewStyle().Margin(0, 0, 0, 0)
	defaultWidth = 20
)

/* State constants */
const (
  LIST = "list"
  INPUT = "input"
)

type State string

type ClusterForm struct {
  inputs  []textinput.Model
	list     list.Model
	choice   string
	quitting bool
  focused int
	err error
  state State
}

type itemDelegate struct{}

type item string

func (i item) FilterValue() string { return "" }

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}



func NewClusterForm(operation string) m.Model { 

  var l list.Model = newList()

  var inputs []textinput.Model = newInputFields(operation)

  var txtInput *ClusterForm = new(ClusterForm)
	txtInput.inputs = inputs
  txtInput.focused = 0
  txtInput.err = nil
  txtInput.list = l
  txtInput.state = LIST
  txtInput.choice = "Rename"

  return txtInput
}

func (this ClusterForm) Init() tea.Cmd {
	return textinput.Blink 
}

func (this ClusterForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
  switch this.state {
    case LIST:
      return this.updateList(msg)
    case INPUT:
      return this.updateInputFields(msg)
    default:
      return this.updateList(msg)
  }
}


func (this ClusterForm) View() string {

  switch this.state {
  case INPUT:
    return this.viewInputsField()
  case LIST : 
    return this.viewList() 
  default:
    return this.viewList()
}
}


func (this ClusterForm) Show() (map[string]string, error) {
  l, err := tea.NewProgram(this).Run()

  form := l.(*ClusterForm)
  var value map[string]string = map[string]string{}
  value["fileName"] = form.inputs[0].Value()
  value["operation"] = form.choice
  if len(form.inputs) == 2 {
    value["addParam"] = form.inputs[1].Value()
  }
  return value, err

}

func (this *ClusterForm) nextInput() {
	this.focused = (this.focused + 1) % len(this.inputs)
}

func (this *ClusterForm) prevInput() {
	this.focused--
	// Wrap around
	if this.focused < 0 {
		this.focused = len(this.inputs) - 1
	}
}

func (this *ClusterForm) viewInputsField() string {
  var title string = "" 
  var subtitle string = ""
	if this.choice != "" {
    switch this.choice {
    case "Read":  
      title = "Name the file you want to read: "
    case "Write":  
      title = "Name the file you want to write: "
      subtitle = "Insert file path: "
    case "Delete":  
      title = "Name the file you want to delete: "
    case "Create":  
      title = "Name of file: "
      subtitle = "Insert file path: "
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
		  `

  %s  
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

func (this *ClusterForm) viewList() string {
	return "\n" + this.list.View()
}

func newList() list.Model {
  var items []list.Item
  items = []list.Item{
		item("Read"),
		item("Write"),
		item("Delete"),
		item("Create"),
		item("Rename"),
	}
	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "What do you want to do?"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
  l.SetShowHelp(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
  
  return l
}

func newInputFields(operation string) []textinput.Model {
  var inputs []textinput.Model
  switch operation {
  case "Rename", "Create", "Write" :
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

  return inputs
}

func (this *ClusterForm) updateList(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		this.list.SetWidth(msg.Width)
		return this, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			this.quitting = true
			return this, tea.Quit

		case "enter":
			i, ok := this.list.SelectedItem().(item)
			if ok {
				this.choice = string(i)
        this.inputs = newInputFields(this.choice)
        this.state = INPUT
			}
			return this, nil
		}
}
	var cmd tea.Cmd
	this.list, cmd = this.list.Update(msg)
  return this, cmd
}

func (this *ClusterForm) updateInputFields(msg tea.Msg) (tea.Model, tea.Cmd) {
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
    case tea.KeyCtrlZ:
      this.choice = ""
      this.state = LIST
		}
		for i := range this.inputs {
			this.inputs[i].Blur()
		}
		this.inputs[this.focused].Focus()
	}

	for i := range this.inputs {
		this.inputs[i], cmds[i] = this.inputs[i].Update(msg)
	}
	return this, tea.Batch(cmds...)

}
