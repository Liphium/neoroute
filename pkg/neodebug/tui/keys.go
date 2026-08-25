package tui

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

// Key standards
var (
	standardUpKey    = key.WithKeys("up", "shift+tab", "k")
	standardDownKey  = key.WithKeys("down", "tab", "j")
	standardRightKey = key.WithKeys("right", "l")
	standardLeftKey  = key.WithKeys("left", "h")
)

// matchesFocusedNavigation checks navigation keybind, excluding Vim keys when input focused.
func matchesFocusedNavigation(msg tea.KeyPressMsg, focused bool, binding key.Binding) bool {
	if focused {
		return matchesNoVim(msg, binding)
	}
	return key.Matches(msg, binding)
}

// matchesNoVim checks if a key is pressed without any 1 character keys
func matchesNoVim(msg tea.KeyPressMsg, binding key.Binding) bool {
	keys := binding.Keys()
	noVim := []string{}
	for _, k := range keys {
		if len(k) != 1 {
			noVim = append(noVim, k)
		}
	}

	temp := key.NewBinding(key.WithKeys(noVim...))
	return key.Matches(msg, temp)
}

type FullKeyHelp struct {
	Title string
	Keys  [][]key.Binding
}

type keyProvider interface {
	Children() []keyProvider

	FooterKeys() []key.Binding
	FullKeyHelp() FullKeyHelp
}
