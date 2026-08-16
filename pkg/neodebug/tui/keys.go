package tui

import (
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/viewport"
)

type appKeyMap struct {
	ScrollingActive bool // If the viewport has been scrolled (will stay at that point, this will show the scroll to bottom key)

	Quit           key.Binding
	GoToBottom     key.Binding
	ExpandHistory  key.Binding
	ExpandInput    key.Binding
	ExitFullscreen key.Binding
	ViewportKeys   viewport.KeyMap
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
		ViewportKeys: viewport.DefaultKeyMap(),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),
	}
}

func (k appKeyMap) ShortHelp() []key.Binding {
	keys := []key.Binding{k.ViewportKeys.Up, k.ViewportKeys.Down, k.Help, k.Quit}

	// When we're scrolling, the snap to bottom thingy should be at the complete left (so it's easily visible).
	if k.ScrollingActive {
		keys = append([]key.Binding{k.GoToBottom}, keys...)
	}

	return keys
}

func (k appKeyMap) FullHelp() [][]key.Binding {

	// TODO: Add the viewport and other keys in here anywhere
	return [][]key.Binding{
		{k.GoToBottom},
		{k.Help, k.Quit},
	}
}
