package tui

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Quit key.Binding
	Cancel key.Binding
}

var DefaultKeyMap = KeyMap {
	Quit: key.NewBinding(
		key.WithKeys(),
	),
}