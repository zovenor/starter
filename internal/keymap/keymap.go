package keymap

import (
	"github.com/charmbracelet/bubbles/key"
)

type KeyMap struct {
	Up      key.Binding
	Down    key.Binding
	Run     key.Binding
	RunAll  key.Binding
	Stop    key.Binding
	StopAll key.Binding
	Detach  key.Binding
	Quit    key.Binding
}

var Keys = KeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Run: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "run executable app"),
	),
	RunAll: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "run all executable apps"),
	),
	Stop: key.NewBinding(
		key.WithKeys("ctrl+r"),
		key.WithHelp("ctrl+r", "stop executable app"),
	),
	StopAll: key.NewBinding(
		key.WithKeys("ctrl+a"),
		key.WithHelp("ctrl+a", "stop all executable apps"),
	),
	Detach: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "detach/attach app"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q/esc/ctrl+c", "quit"),
	),
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Run, k.RunAll, k.Stop, k.StopAll, k.Detach, k.Quit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down},
		{k.Run, k.RunAll, k.Stop, k.StopAll, k.Detach, k.Quit},
	}
}
