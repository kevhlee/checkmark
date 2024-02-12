package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

type ConfirmModel struct {
	Action Action

	form   *huh.Form
	parent tea.Model
}

func NewConfirm(parent tea.Model, action Action, message string) *ConfirmModel {
	model := &ConfirmModel{
		Action: action,
		parent: parent,
	}

	model.form = huh.NewForm(
		huh.NewGroup(huh.NewConfirm().
			Key("confirm").
			Title(message).
			Affirmative("Confirm").
			Negative("Decline"),
		),
	)

	model.form.WithKeyMap(defaultFormKeys)
	model.form.WithWidth(120)

	return model
}

func (m ConfirmModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m *ConfirmModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	model, cmd := m.form.Update(msg)
	if form, ok := model.(*huh.Form); ok {
		m.form = form
	}

	switch m.form.State {
	case huh.StateAborted:
		return m.parent, nil

	case huh.StateCompleted:
		if m.form.GetBool("confirm") {
			return m.parent.Update(*m)
		} else {
			return m.parent, nil
		}

	default:
		return m, cmd
	}
}

func (m ConfirmModel) View() string {
	return m.form.View()
}
