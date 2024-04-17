package lister

import (
     "fmt"
	   "io"
	   "os"
	   "strings" 

     "github.com/charmbracelet/bubbles/list"
 	   tea "github.com/charmbracelet/bubbletea"
	   "github.com/charmbracelet/lipgloss"
)

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(0)
  listHeight = 14
	itemStyle         = lipgloss.NewStyle().PaddingLeft(0)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(0).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(0)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(0).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(0, 0, 0, 0)
	defaultWidth = 20
)

type List struct {}

type model struct {
	list     list.Model
	choice   string
	quitting bool
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


func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = string(i)
			}
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.choice != "" {
    switch m.choice {
    case "Read":  
      return quitTextStyle.Render(fmt.Sprint("Name the file you want to read: "))
    case "Write":  
      return quitTextStyle.Render(fmt.Sprint("Name the file you want to write: "))
    case "Delete":  
      return quitTextStyle.Render(fmt.Sprint("Name the file you want to delete: "))
    case "Create":  
      return quitTextStyle.Render(fmt.Sprint("Name of file: "))
    case "Rename":  
      return quitTextStyle.Render(fmt.Sprint("Name the file you want to rename and the new name: "))
    default:
      return quitTextStyle.Render(fmt.Sprint("Error: unknown operation"))
  }
	}
	if m.quitting {
		return quitTextStyle.Render("Quitting.")
	}
	return "\n" + m.list.View()
}


func (this *List) Show() {

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
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

  var m model

	m = model{list: l}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
