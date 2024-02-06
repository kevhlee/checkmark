package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kevhlee/checkmark/task"
)

type EditorKeyMap struct {
	Left  key.Binding
	Right key.Binding
	Enter key.Binding
	Quit  key.Binding
}

type EditorModel struct {
	Action Action

	index      int
	keys       EditorKeyMap
	nameInput  textinput.Model
	parent     tea.Model
	priorities []task.Priority
	quitting   bool
}

func NewEditor(parent tea.Model, action Action, name string, priority task.Priority) *EditorModel {
	model := &EditorModel{
		Action:     action,
		index:      int(priority),
		nameInput:  textinput.New(),
		parent:     parent,
		priorities: []task.Priority{task.LowPriority, task.HighPriority, task.FirePriority},
	}

	model.keys = EditorKeyMap{
		Left:  key.NewBinding(key.WithKeys("left")),
		Right: key.NewBinding(key.WithKeys("right")),
		Enter: key.NewBinding(key.WithKeys("enter")),
		Quit:  key.NewBinding(key.WithKeys("ctrl+c", "esc")),
	}

	model.nameInput.SetValue(name)
	model.nameInput.Focus()

	return model
}

func (m EditorModel) Name() string {
	return strings.TrimSpace(m.nameInput.Value())
}

func (m EditorModel) Priority() task.Priority {
	return m.priorities[m.index]
}

func (m EditorModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *EditorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Left):
			if !m.nameInput.Focused() {
				m.index = max(m.index-1, 0)
			}

		case key.Matches(msg, m.keys.Right):
			if !m.nameInput.Focused() {
				m.index = min(m.index+1, len(m.priorities)-1)
			}

		case key.Matches(msg, m.keys.Quit):
			m.quitting = true
			return m.parent, nil

		case key.Matches(msg, m.keys.Enter):
			if m.nameInput.Focused() {
				m.nameInput.Blur()
			} else {
				m.quitting = true
				return m.parent.Update(*m)
			}
		}
	}

	var cmd tea.Cmd

	if m.nameInput.Focused() {
		m.nameInput, cmd = m.nameInput.Update(msg)
	}

	return m, cmd
}

func (m EditorModel) View() string {
	if m.quitting {
		return ""
	}

	nameBuilder := strings.Builder{}
	nameBuilder.WriteString("Name:")
	nameBuilder.WriteString("\n")
	nameBuilder.WriteString(m.nameInput.View())
	nameBuilder.WriteString("\n")

	priorityBuilder := strings.Builder{}
	priorityBuilder.WriteString("Priority:")
	priorityBuilder.WriteString("\n")

	for i, priority := range m.priorities {
		cursor := "  "
		if i == m.index {
			cursor = "> "
		}

		if m.nameInput.Focused() {
			priorityBuilder.WriteString(cursor)
		} else {
			priorityBuilder.WriteString(styleCursor.Render(cursor))
		}

		priorityBuilder.WriteString(priority.Symbol())
		priorityBuilder.WriteString(" ")
		priorityBuilder.WriteString(priority.String())
		priorityBuilder.WriteString("    ")
	}
	priorityBuilder.WriteString("\n")

	nameView := nameBuilder.String()
	priorityView := priorityBuilder.String()

	if m.nameInput.Focused() {
		nameView = styleEditorSelected.Render(nameView)
		priorityView = styleEditorUnselected.Render(priorityView)
	} else {
		nameView = styleEditorUnselected.Render(nameView)
		priorityView = styleEditorSelected.Render(priorityView)
	}

	return lipgloss.JoinVertical(lipgloss.Left, nameView, priorityView)
}
