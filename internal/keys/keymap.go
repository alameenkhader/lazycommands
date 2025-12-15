package keys

import "github.com/charmbracelet/bubbles/key"

// KeyMap defines the keyboard shortcuts for the application
type KeyMap struct {
	Up       key.Binding
	Down     key.Binding
	Continue key.Binding
	Stop     key.Binding
	Quit     key.Binding
}

// DefaultKeyMap returns the default key bindings
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "move up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "move down"),
		),
		Continue: key.NewBinding(
			key.WithKeys("y"),
			key.WithHelp("y", "continue"),
		),
		Stop: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "stop"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q/ctrl+c", "quit"),
		),
	}
}
