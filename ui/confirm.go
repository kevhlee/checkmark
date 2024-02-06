package ui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ConfirmKeyMap struct {
	Left  key.Binding
	Right key.Binding
	Enter key.Binding
	Quit  key.Binding
}

type ConfirmModel struct {
	Action Action

	confirm  bool
	keys     ConfirmKeyMap
	message  string
	parent   tea.Model
	quitting bool
}

func NewConfirm(parent tea.Model, action Action, message string) *ConfirmModel {
	model := &ConfirmModel{
		Action:  action,
		confirm: true,
		message: message,
		parent:  parent,
	}

	model.keys = ConfirmKeyMap{
		Left:  key.NewBinding(key.WithKeys("left")),
		Right: key.NewBinding(key.WithKeys("right")),
		Enter: key.NewBinding(key.WithKeys("enter")),
		Quit:  key.NewBinding(key.WithKeys("ctrl+c", "esc")),
	}

	return model
}

func (m ConfirmModel) Init() tea.Cmd {
	return nil
}

func (m *ConfirmModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Left):
			m.confirm = true

		case key.Matches(msg, m.keys.Right):
			m.confirm = false

		case key.Matches(msg, m.keys.Quit):
			m.confirm = false
			fallthrough

		case key.Matches(msg, m.keys.Enter):
			m.quitting = true
			if m.confirm {
				return m.parent.Update(*m)
			} else {
				return m.parent, nil
			}
		}
	}

	return m, nil
}

func (m ConfirmModel) View() string {
	if m.quitting {
		return ""
	}

	viewConfirm := "Confirm"
	viewDecline := "Decline"

	if m.confirm {
		viewConfirm = styleConfirmSelected.Render(viewConfirm)
		viewDecline = styleConfirmUnselected.Render(viewDecline)
	} else {
		viewConfirm = styleConfirmUnselected.Render(viewConfirm)
		viewDecline = styleConfirmSelected.Render(viewDecline)
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.message,
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			viewConfirm,
			viewDecline,
		),
	)
}
