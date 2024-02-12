package ui

import (
	"fmt"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kevhlee/checkmark/config"
	"github.com/kevhlee/checkmark/task"
)

type Action int

const (
	AddNewTask Action = iota
	ClearCompletedTasks
	DeleteCurrentTask
	EditCurrentTask
)

type RootKeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Add    key.Binding
	Clear  key.Binding
	Delete key.Binding
	Edit   key.Binding
	Mark   key.Binding
	Help   key.Binding
	Quit   key.Binding
}

type RootModel struct {
	help     help.Model
	index    int
	keys     RootKeyMap
	quitting bool
	tasks    []task.Task
}

func New(cfg config.Config) *RootModel {
	model := &RootModel{
		help:     help.New(),
		index:    0,
		tasks:    cfg.Tasks,
		quitting: false,
	}

	model.keys = RootKeyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑ / k", "move up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓ / j", "move down"),
		),
		Add: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "add task"),
		),
		Clear: key.NewBinding(
			key.WithKeys("c"),
			key.WithHelp("c", "clear completed tasks"),
		),
		Delete: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "delete task"),
		),
		Edit: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("e", "edit task"),
		),
		Mark: key.NewBinding(
			key.WithKeys(" ", "enter"),
			key.WithHelp("Space", "mark task"),
		),
		Help: key.NewBinding(
			key.WithKeys("h"),
			key.WithHelp("h", "toggle help"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c", "esc"),
			key.WithHelp("q", "quit"),
		),
	}

	return model
}

func (m RootModel) ShortHelp() []key.Binding {
	return []key.Binding{m.keys.Help, m.keys.Quit}
}

func (m RootModel) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{m.keys.Up, m.keys.Down},
		{m.keys.Add, m.keys.Clear, m.keys.Delete, m.keys.Edit, m.keys.Mark},
		{m.keys.Help, m.keys.Quit},
	}
}

func (m RootModel) Init() tea.Cmd {
	return nil
}

func (m *RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ConfirmModel:
		switch msg.Action {
		case ClearCompletedTasks:
			m.tasks = slices.DeleteFunc(m.tasks, func(t task.Task) bool {
				return t.Done
			})
			m.normalizeCursor()

		case DeleteCurrentTask:
			m.tasks = append(m.tasks[:m.index], m.tasks[m.index+1:]...)
			m.normalizeCursor()
		}

	case EditorModel:
		switch msg.Action {
		case AddNewTask:
			m.tasks = append(m.tasks, task.Task{Name: msg.Name, Priority: msg.Priority})

		case EditCurrentTask:
			m.tasks[m.index].Name = msg.Name
			m.tasks[m.index].Priority = msg.Priority
		}

		task.SortTasks(m.tasks)

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Up):
			m.index = max(m.index-1, 0)

		case key.Matches(msg, m.keys.Down):
			m.index = min(m.index+1, len(m.tasks)-1)

		case key.Matches(msg, m.keys.Add):
			editor := NewEditor(m, AddNewTask, "", task.LowPriority)
			return editor, editor.Init()

		case key.Matches(msg, m.keys.Edit):
			if len(m.tasks) > 0 {
				editor := NewEditor(m, EditCurrentTask, m.tasks[m.index].Name, m.tasks[m.index].Priority)
				return editor, editor.Init()
			}

		case key.Matches(msg, m.keys.Clear):
			if len(m.tasks) > 0 {
				confirm := NewConfirm(m, ClearCompletedTasks, "Clear completed tasks?")
				return confirm, confirm.Init()
			}

		case key.Matches(msg, m.keys.Delete):
			if len(m.tasks) > 0 {
				confirm := NewConfirm(m, DeleteCurrentTask, fmt.Sprintf("Delete '%s'?", m.tasks[m.index].Name))
				return confirm, confirm.Init()
			}

		case key.Matches(msg, m.keys.Mark):
			if len(m.tasks) > 0 {
				m.tasks[m.index].Done = !m.tasks[m.index].Done
			}

		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll

		case key.Matches(msg, m.keys.Quit):
			m.quitting = true
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m RootModel) View() string {
	if m.quitting {
		return ""
	}

	tasksView := strings.Builder{}

	for i, task := range m.tasks {
		if i == m.index {
			tasksView.WriteString(styleCursor.Render("> "))
		} else {
			tasksView.WriteString(styleCursor.Render("  "))
		}

		if task.Done {
			tasksView.WriteString(styleCheck.Render("✓  "))
		} else {
			tasksView.WriteString(styleCheck.Render("   "))
		}

		tasksView.WriteString(task.Priority.Symbol())
		tasksView.WriteString(" ")
		tasksView.WriteString(task.Name)
		tasksView.WriteString("\n")
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		styleHeader.Render("Tasks:"),
		tasksView.String(),
		m.help.View(m),
	)
}

func (m *RootModel) normalizeCursor() {
	m.index = min(m.index, len(m.tasks)-1)
}
