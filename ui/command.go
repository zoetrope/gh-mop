package ui

import (
	"github.com/charmbracelet/bubbles/textinput"
)

func newCommandModel() textinput.Model {
	ti := textinput.New()
	ti.Placeholder = "command"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20
	return ti
}
