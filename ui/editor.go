package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/kevhlee/checkmark/task"
)

type EditorModel struct {
	Action   Action
	Name     string
	Priority task.Priority

	form   *huh.Form
	parent tea.Model
}

func NewEditor(parent tea.Model, action Action, name string, priority task.Priority) *EditorModel {
	model := &EditorModel{
		Action:   action,
		Name:     name,
		Priority: priority,
		parent:   parent,
	}

	model.form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("name").
				Value(&model.Name).
				Title("Name: ").
				CharLimit(120).
				Validate(func(s string) error {
					if len(strings.TrimSpace(s)) == 0 {
						return fmt.Errorf("Empty input")
					}
					return nil
				}),
			huh.NewSelect[task.Priority]().
				Key("priority").
				Value(&model.Priority).
				Title("Priority:").
				Options(huh.NewOptions(task.Priorities...)...),
		),
	)

	model.form = model.form.WithKeyMap(defaultFormKeys)
	model.form = model.form.WithWidth(120)

	return model
}

func (m EditorModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m *EditorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	model, cmd := m.form.Update(msg)
	if form, ok := model.(*huh.Form); ok {
		m.form = form
	}

	switch m.form.State {
	case huh.StateAborted:
		return m.parent, nil

	case huh.StateCompleted:
		return m.parent.Update(*m)

	default:
		return m, cmd
	}
}

func (m EditorModel) View() string {
	return m.form.View()
}
