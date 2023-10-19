package todoFinder

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"os"
)

// TODO: Change render function of list so long todos are not turned to elipses, but are printed for multiple lines (this may not be possible)
func CreateTodoList(todos []Todo, dirToSearch string) {
	// Convert Todos to items type. Need to create arrays of both types of items so can access all info for items in the model
	todoItems := []item{}
	listItems := []list.Item{}
	for _, todo := range todos {
		desc := todo.FilePath() + " Line: " + todo.LineNumber()
		newItem := item{todo.Text(), desc}
		todoItems = append(todoItems, newItem)
		listItems = append(listItems, newItem)
	}

	m := model{list: list.New(listItems, list.NewDefaultDelegate(), 0, 0), items: todoItems, selected: make(map[int]bool)}
	m.list.Title = "Todos for " + dirToSearch

	// Add key bindings for enter key
	m.list.AdditionalShortHelpKeys = addShortHelpKeys
	m.list.AdditionalFullHelpKeys = addFullHelpKeys

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title, desc string
}

type itemArray []item

// Getters
func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.desc + i.title }

type model struct {
	list     list.Model
	items    itemArray
	selected map[int]bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			index := m.list.Index()
			m.ToggleCheckMark(index)
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) ToggleCheckMark(index int) {
	_, ok := m.selected[index]
	originalItem := m.items[index]
	// If the key exists remove from map, otherwise add it
	if ok {
		m.list.SetItem(index, originalItem)
		delete(m.selected, index)
	} else {
		checkedItem := item{"✅ " + originalItem.Title(), originalItem.Description()}
		m.list.SetItem(index, checkedItem)
		m.selected[index] = true
	}
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

func addShortHelpKeys() []key.Binding {
	newBindings := []key.Binding{}
	newBindings = append(newBindings, key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "mark complete"),
	))
	return newBindings
}

func addFullHelpKeys() []key.Binding {
	newBindings := []key.Binding{}
	newBindings = append(newBindings, key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "mark complete"),
	))
	return newBindings
}
