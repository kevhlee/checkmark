package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kevhlee/checkmark/config"
	"github.com/kevhlee/checkmark/ui"
)

func main() {
	// var (
	// 	name     string
	// 	priority task.Priority
	// )
	//
	// form := huh.NewForm(
	// 	huh.NewGroup(
	// 		huh.NewInput().
	// 			Title("Title:").
	// 			Value(&name).
	// 			Validate(func(s string) error {
	// 				if len(strings.TrimSpace(s)) == 0 {
	// 					return fmt.Errorf("Empty input")
	// 				}
	// 				return nil
	// 			}),
	// 		huh.NewSelect[task.Priority]().
	// 			Title("Priority").
	// 			Options(
	// 				huh.NewOption(task.LowPriority.FullString(), task.LowPriority),
	// 				huh.NewOption(task.HighPriority.FullString(), task.HighPriority),
	// 				huh.NewOption(task.FirePriority.FullString(), task.FirePriority),
	// 			).
	// 			Value(&priority),
	// 	),
	// )
	//
	// if err := form.Run(); err != nil {
	// 	panic(err)
	// }
	// fmt.Println(name, priority)

	cfg, err := config.Init()
	if err != nil {
		panic(err)
	}

	program := tea.NewProgram(ui.New(cfg))
	if _, err := program.Run(); err != nil {
		panic(err)
	}

	if err := config.Save(cfg); err != nil {
		panic(err)
	}
}
