package tui

import (
	"charm.land/bubbles/v2/key"
)

type FullKeyHelp struct {
	Title string
	Keys  [][]key.Binding
}

type keyProvider interface {
	Children() []keyProvider

	FooterKeys() []key.Binding
	FullKeyHelp() FullKeyHelp
}

type appKeyMap struct {
	Quit           key.Binding
	GoToBottom     key.Binding
	ExpandHistory  key.Binding
	ExpandInput    key.Binding
	ExitFullscreen key.Binding
	Help           key.Binding
}

func newAppKeyMap() appKeyMap {
	return appKeyMap{
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		GoToBottom: key.NewBinding(
			key.WithKeys("ctrl+b"),
			key.WithHelp("ctrl+b", "go to bottom"),
		),
		ExpandHistory: key.NewBinding(
			key.WithKeys("ctrl+h"),
			key.WithHelp("ctrl+h", "expand history"),
		),
		ExpandInput: key.NewBinding(
			key.WithKeys("ctrl+e"),
			key.WithHelp("ctrl+e", "expand input"),
		),
		ExitFullscreen: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("escape", "exit fullscreen"),
		),
		Help: key.NewBinding(
			key.WithKeys("ctrl+h"),
			key.WithHelp("ctrl+h", "toggle help"),
		),
	}
}
