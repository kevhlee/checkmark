package ui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/huh"
)

var defaultFormKeys = huh.NewDefaultKeyMap()

func init() {
	defaultFormKeys.Quit = key.NewBinding(
		key.WithKeys("ctrl+c", "esc"),
	)
}
