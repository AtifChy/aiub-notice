package list

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/evertras/bubble-table/table"
)

type KeyMap struct {
	table.KeyMap
	RowOpen key.Binding
	Help    key.Binding
	Quit    key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.RowDown, k.RowUp, k.RowSelectToggle, k.RowOpen, k.Help, k.Quit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.RowDown, k.RowUp},
		{k.RowSelectToggle, k.RowOpen},
		{k.Filter, k.FilterBlur, k.FilterClear},
		{k.Help, k.Quit},
	}
}

func DefaultKeyMap() KeyMap {
	km := table.DefaultKeyMap()
	km.RowSelectToggle = key.NewBinding(
		key.WithKeys(" "),
		key.WithHelp("space", "select"),
	)
	return KeyMap{
		KeyMap: km,
		RowOpen: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "open"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "more"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q/ctrl+c", "quit"),
		),
	}
}
