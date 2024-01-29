package main

import (
	"fmt"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

//
// Styles
//

var (
	styleCheck  = lipgloss.NewStyle().Foreground(lipgloss.Color("#afffd0"))
	styleCursor = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff80ee"))

	styleHeader = lipgloss.NewStyle().Bold(true).
			Foreground(lipgloss.Color("#acffee"))

	styleSelected = lipgloss.NewStyle().
			Margin(1, 1).
			Padding(0, 3).
			Background(lipgloss.Color("212")).
			Foreground(lipgloss.Color("230"))

	styleUnselected = lipgloss.NewStyle().
			Margin(1, 1).
			Padding(0, 3).
			Background(lipgloss.Color("235")).
			Foreground(lipgloss.Color("254"))

	cursorSelected   = styleCursor.Render("> ")
	cursorUnselected = styleCursor.Render("  ")
	statusDone       = styleCheck.Render("✓  ")
	statusPending    = styleCheck.Render("   ")
)

//
// Keys
//

type KeyMap struct {
	Add    key.Binding
	Check  key.Binding
	Clear  key.Binding
	Delete key.Binding
	Edit   key.Binding
	Enter  key.Binding
	Up     key.Binding
	Down   key.Binding
	Left   key.Binding
	Right  key.Binding
	Help   key.Binding
	Quit   key.Binding
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down},
		{k.Add, k.Check, k.Delete, k.Clear, k.Edit},
		{k.Help, k.Quit},
	}
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

var keys = KeyMap{
	Add: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "Add new task"),
	),
	Check: key.NewBinding(
		key.WithKeys(tea.KeyEnter.String(), tea.KeySpace.String()),
		key.WithHelp("Enter / Space", "Mark task"),
	),
	Clear: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "Clear completed tasks"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "Delete task"),
	),
	Edit: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "Edit task"),
	),
	Enter: key.NewBinding(
		key.WithKeys(tea.KeyEnter.String()),
	),
	Up: key.NewBinding(
		key.WithKeys(tea.KeyUp.String(), "k"),
		key.WithHelp("↑ / k", "Go up"),
	),
	Down: key.NewBinding(
		key.WithKeys(tea.KeyDown.String(), "j"),
		key.WithHelp("↓ / j", "Go down"),
	),
	Left: key.NewBinding(
		key.WithKeys(tea.KeyLeft.String(), "h"),
		key.WithHelp("← / h", "Go left"),
	),
	Right: key.NewBinding(
		key.WithKeys(tea.KeyRight.String(), "l"),
		key.WithHelp("→ / l", "Go right"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "Toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys(tea.KeyEsc.String(), tea.KeyCtrlC.String()),
		key.WithHelp("Esc / Ctrl+C", "Quit"),
	),
}

//
// TUI
//

type TUI struct {
	Tasks []Task

	cursor   int
	help     help.Model
	quitting bool
}

func StartTUI(config *Config) error {
	tui := &TUI{
		Tasks:  config.Tasks,
		cursor: 0,
		help:   help.New(),
	}

	model, err := tea.NewProgram(tui).Run()
	if err != nil {
		return err
	}
	config.Tasks = model.(*TUI).Tasks
	return nil
}

func (t TUI) Init() tea.Cmd {
	return nil
}

func (t *TUI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Add):
			form := NewForm(t, false, "")
			return form, form.Init()

		case key.Matches(msg, keys.Edit):
			if len(t.Tasks) > 0 {
				form := NewForm(t, true, t.Tasks[t.cursor].Name)
				return form, form.Init()
			}

		case key.Matches(msg, keys.Check):
			if len(t.Tasks) > 0 {
				t.Tasks[t.cursor].Done = !t.Tasks[t.cursor].Done
			}

		case key.Matches(msg, keys.Clear):
			if len(t.Tasks) > 0 {
				confirm := NewConfirm(t, ClearCompletedTasks, "Clear completed tasks?")
				return confirm, confirm.Init()
			}

		case key.Matches(msg, keys.Delete):
			if len(t.Tasks) > 0 {
				confirm := NewConfirm(t, DeleteTask, fmt.Sprintf("Delete '%s'?", t.Tasks[t.cursor].Name))
				return confirm, confirm.Init()
			}

		case key.Matches(msg, keys.Up):
			t.PrevTask()

		case key.Matches(msg, keys.Down):
			t.NextTask()

		case key.Matches(msg, keys.Help):
			t.help.ShowAll = !t.help.ShowAll

		case key.Matches(msg, keys.Quit):
			t.quitting = true
			cmd = tea.Quit
		}

	case Confirm:
		if msg.confirmation {
			switch msg.action {
			case DeleteTask:
				t.Tasks = append(t.Tasks[:t.cursor], t.Tasks[t.cursor+1:]...)
				t.NormalizeCursor()

			case ClearCompletedTasks:
				t.Tasks = slices.DeleteFunc(t.Tasks, func(task Task) bool {
					return task.Done
				})
				t.NormalizeCursor()
			}
		}

	case Form:
		if msg.edit {
			t.Tasks[t.cursor].Name = msg.title.Value()
		} else {
			t.Tasks = append(t.Tasks, Task{Name: msg.title.Value()})
		}

	}

	return t, cmd
}

func (t TUI) View() string {
	if t.quitting {
		return ""
	}

	tasksView := strings.Builder{}

	for i, task := range t.Tasks {
		if i == t.cursor {
			tasksView.WriteString(cursorSelected)
		} else {
			tasksView.WriteString(cursorUnselected)
		}

		if task.Done {
			tasksView.WriteString(statusDone)
		} else {
			tasksView.WriteString(statusPending)
		}

		tasksView.WriteString(task.Name)
		tasksView.WriteString("\n")
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		styleHeader.Render("Tasks: "),
		tasksView.String(),
		t.help.View(keys),
	)
}

func (t *TUI) PrevTask() {
	t.cursor = max(t.cursor-1, 0)
}

func (t *TUI) NextTask() {
	t.cursor = min(t.cursor+1, len(t.Tasks)-1)
}

func (t *TUI) NormalizeCursor() {
	if len(t.Tasks) == 0 {
		t.cursor = 0
	} else {
		t.cursor = min(t.cursor, len(t.Tasks)-1)
	}
}

//
// Form
//

type Form struct {
	parent tea.Model

	edit  bool
	title textinput.Model
}

func NewForm(parent tea.Model, edit bool, title string) Form {
	form := Form{
		edit:   edit,
		parent: parent,
		title:  textinput.New(),
	}

	form.title.SetValue(title)
	form.title.Focus()

	return form
}

func (f Form) Init() tea.Cmd {
	return textinput.Blink
}

func (f Form) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			return f.parent, nil

		case key.Matches(msg, keys.Enter):
			return f.parent.Update(f)
		}
	}

	var cmd tea.Cmd
	f.title, cmd = f.title.Update(msg)
	return f, cmd
}

func (f Form) View() string {
	var headerView string
	if f.edit {
		headerView = styleHeader.Render("Edit task:")
	} else {
		headerView = styleHeader.Render("Create new task:")
	}
	return lipgloss.JoinVertical(lipgloss.Left, headerView, f.title.View())
}

//
// Confirm
//

type ConfirmAction int

const (
	DeleteTask ConfirmAction = iota
	ClearCompletedTasks
)

type Confirm struct {
	parent tea.Model

	action       ConfirmAction
	message      string
	confirmation bool
}

func NewConfirm(parent tea.Model, action ConfirmAction, message string) *Confirm {
	return &Confirm{
		action:       action,
		confirmation: true,
		message:      message,
		parent:       parent,
	}
}

func (c Confirm) Init() tea.Cmd {
	return nil
}

func (c *Confirm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Left) || msg.String() == "y":
			c.confirmation = true

		case key.Matches(msg, keys.Right) || msg.String() == "n":
			c.confirmation = false

		case key.Matches(msg, keys.Quit):
			c.confirmation = false
			fallthrough

		case key.Matches(msg, keys.Enter):
			return c.parent.Update(*c)
		}
	}

	return c, nil
}

func (c Confirm) View() string {
	msgYes := "Yes"
	msgNo := "No"

	if c.confirmation {
		msgYes = styleSelected.Render(msgYes)
		msgNo = styleUnselected.Render(msgNo)
	} else {
		msgYes = styleUnselected.Render(msgYes)
		msgNo = styleSelected.Render(msgNo)
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		styleHeader.Render(c.message),
		lipgloss.JoinHorizontal(lipgloss.Left, msgYes, msgNo),
	)
}
