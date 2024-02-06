package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kevhlee/checkmark/config"
	"github.com/kevhlee/checkmark/ui"
)

func main() {
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
