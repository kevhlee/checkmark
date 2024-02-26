package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kevhlee/checkmark/task"
	"github.com/kevhlee/checkmark/ui"
	"github.com/spf13/viper"
)

func init() {
	viper.AddConfigPath("$HOME/.config/checkmark")
	viper.SetConfigName("config")
	viper.SetConfigType("json")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if err := viper.SafeWriteConfig(); err != nil {
				panic(err)
			}
			return
		}
		panic(err)
	}
}

func main() {
	var tasks []task.Task
	if err := viper.UnmarshalKey("tasks", &tasks); err != nil {
		panic(err)
	}

	model := ui.New(tasks)
	if _, err := tea.NewProgram(model).Run(); err != nil {
		panic(err)
	}

	viper.Set("tasks", model.Tasks)

	if err := viper.WriteConfig(); err != nil {
		panic(err)
	}
}
