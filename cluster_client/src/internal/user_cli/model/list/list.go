package lister

import (
     "fmt"
	   "io"
	   "os"
	   "strings" 
     "github.com/charmbracelet/bubbles/list"
 	   tea "github.com/charmbracelet/bubbletea"
	   "github.com/charmbracelet/lipgloss"
     m "raft/client/src/internal/user_cli/model"
)

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(0)
  listHeight = 10
	itemStyle         = lipgloss.NewStyle().PaddingLeft(2)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(0).PaddingBottom(0)
	quitTextStyle     = lipgloss.NewStyle().Margin(0, 0, 0, 0)
	defaultWidth = 20
)

type List struct {
	list     list.Model
	choice   string
	quitting bool
}

func NewList() m.Model { 
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

  var list *List = new(List)
  list.list = l

	return list
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


func (this List) Init() tea.Cmd {
	return nil
}

func (this List) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			}
			return this, tea.Quit
		}
}
	var cmd tea.Cmd
	this.list, cmd = this.list.Update(msg)
	return this, cmd
}


func (this List) View() string {
	return "\n" + this.list.View()
}


func (this *List) Show() string {
  l, err := tea.NewProgram(this).Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

  list := l.(List)
  return list.choice
}
