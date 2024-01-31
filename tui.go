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
	cursor   int
	help     help.Model
	quitting bool
	tasks    []Task
}

func StartTUI(config *TaskConfig) error {
	tui := &TUI{
		tasks: config.Tasks,
		help:  help.New(),
	}

	if _, err := tea.NewProgram(tui).Run(); err != nil {
		return err
	}

	config.Tasks = tui.tasks
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
			if len(t.tasks) > 0 {
				form := NewForm(t, true, t.tasks[t.cursor].Name)
				return form, form.Init()
			}

		case key.Matches(msg, keys.Check):
			if len(t.tasks) > 0 {
				t.tasks[t.cursor].Done = !t.tasks[t.cursor].Done
			}

		case key.Matches(msg, keys.Clear):
			if len(t.tasks) > 0 {
				confirm := NewConfirmPrompt(t, ClearCompletedTasks, "Clear completed tasks?")
				return confirm, confirm.Init()
			}

		case key.Matches(msg, keys.Delete):
			if len(t.tasks) > 0 {
				confirm := NewConfirmPrompt(t, DeleteTask, fmt.Sprintf("Delete '%s'?", t.tasks[t.cursor].Name))
				return confirm, confirm.Init()
			}

		case key.Matches(msg, keys.Up):
			t.PrevTask()

		case key.Matches(msg, keys.Down):
			t.NextTask()

		case key.Matches(msg, keys.Help):
			t.help.ShowAll = !t.help.ShowAll

		case key.Matches(msg, keys.Quit) || msg.String() == "q":
			t.quitting = true
			cmd = tea.Quit
		}

	case ConfirmPrompt:
		if msg.Confirmation {
			switch msg.Action {
			case DeleteTask:
				t.tasks = append(t.tasks[:t.cursor], t.tasks[t.cursor+1:]...)
				t.NormalizeCursor()

			case ClearCompletedTasks:
				t.tasks = slices.DeleteFunc(t.tasks, func(task Task) bool {
					return task.Done
				})
				t.NormalizeCursor()
			}
		}

	case Form:
		if msg.edit {
			t.tasks[t.cursor].Name = msg.titlePrompt.Value()
		} else {
			t.tasks = append(t.tasks, Task{Name: msg.titlePrompt.Value()})
		}

	}

	return t, cmd
}

func (t TUI) View() string {
	if t.quitting {
		return ""
	}

	view := strings.Builder{}

	for i, task := range t.tasks {
		if i == t.cursor {
			view.WriteString(cursorSelected)
		} else {
			view.WriteString(cursorUnselected)
		}

		if task.Done {
			view.WriteString(statusDone)
		} else {
			view.WriteString(statusPending)
		}

		view.WriteString(task.Name)
		view.WriteString("\n")
	}

	return lipgloss.JoinVertical(lipgloss.Left, styleHeader.Render("Tasks: "), view.String(), t.help.View(keys))
}

func (t *TUI) PrevTask() {
	t.cursor = max(t.cursor-1, 0)
}

func (t *TUI) NextTask() {
	t.cursor = min(t.cursor+1, len(t.tasks)-1)
}

func (t *TUI) NormalizeCursor() {
	if len(t.tasks) == 0 {
		t.cursor = 0
	} else {
		t.cursor = min(t.cursor, len(t.tasks)-1)
	}
}

//
// Form
//

type Form struct {
	edit        bool
	parent      tea.Model
	titlePrompt textinput.Model
}

func NewForm(parent tea.Model, edit bool, title string) Form {
	form := Form{
		edit:        edit,
		parent:      parent,
		titlePrompt: textinput.New(),
	}

	form.titlePrompt.SetValue(title)
	form.titlePrompt.Focus()

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
	f.titlePrompt, cmd = f.titlePrompt.Update(msg)
	return f, cmd
}

func (f Form) View() string {
	var headerView string
	if f.edit {
		headerView = styleHeader.Render("Edit task:")
	} else {
		headerView = styleHeader.Render("Create new task:")
	}
	return lipgloss.JoinVertical(lipgloss.Left, headerView, f.titlePrompt.View())
}

//
// Confirm
//

type ConfirmAction int

const (
	DeleteTask ConfirmAction = iota
	ClearCompletedTasks
)

type ConfirmPrompt struct {
	Action       ConfirmAction
	Confirmation bool

	message string
	parent  tea.Model
}

func NewConfirmPrompt(parent tea.Model, action ConfirmAction, message string) *ConfirmPrompt {
	return &ConfirmPrompt{
		Action:       action,
		Confirmation: true,

		message: message,
		parent:  parent,
	}
}

func (p ConfirmPrompt) Init() tea.Cmd {
	return nil
}

func (p *ConfirmPrompt) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Left) || msg.String() == "y":
			p.Confirmation = true

		case key.Matches(msg, keys.Right) || msg.String() == "n":
			p.Confirmation = false

		case key.Matches(msg, keys.Quit):
			p.Confirmation = false
			fallthrough

		case key.Matches(msg, keys.Enter):
			if p.Confirmation {
				return p.parent.Update(*p)
			} else {
				return p.parent, nil
			}
		}
	}

	return p, nil
}

func (p ConfirmPrompt) View() string {
	msgYes := "Yes"
	msgNo := "No"

	if p.Confirmation {
		msgYes = styleSelected.Render(msgYes)
		msgNo = styleUnselected.Render(msgNo)
	} else {
		msgYes = styleUnselected.Render(msgYes)
		msgNo = styleSelected.Render(msgNo)
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		styleHeader.Render(p.message),
		lipgloss.JoinHorizontal(lipgloss.Left, msgYes, msgNo),
	)
}
