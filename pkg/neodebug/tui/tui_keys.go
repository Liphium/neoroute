package tui

import (
	"slices"
	"strings"

	"charm.land/bubbles/v2/key"
	"charm.land/lipgloss/v2"
)

var _ keyProvider = tui{}

// Children implements keyProvider.
func (m tui) Children() []keyProvider {
	switch m.full {
	case fullScreenHistory:
		return []keyProvider{m.history}
	case fullScreenInput:
		return []keyProvider{m.input}
	default:
		return []keyProvider{m.history, m.input} // Like this input will have priority over the history
	}
}

// FooterKeys implements keyProvider.
func (m tui) FooterKeys() []key.Binding {
	return []key.Binding{m.helpKey, m.quit}
}

// FullKeyHelp implements keyProvider.
func (m tui) FullKeyHelp() FullKeyHelp {
	return FullKeyHelp{
		Title: "App control",
		Keys: [][]key.Binding{
			[]key.Binding{m.expandInput, m.expandHistory, m.exitFullscreen},
			[]key.Binding{m.goToBottom, m.helpKey, m.quit},
		},
	}
}

func collectFooters(provider keyProvider, current []key.Binding) []key.Binding {

	// Collect all of the bindings having things further down the line overwrite the parent's bindings
	for _, child := range provider.Children() {
		keys := child.FooterKeys()
		toAdd := []key.Binding{}
		for _, k := range keys {
			current = slices.DeleteFunc(current, func(c key.Binding) bool {
				for _, ck := range c.Keys() {
					if slices.Contains(k.Keys(), ck) {
						return true
					}
				}
				return false
			})
			toAdd = append(toAdd, k)
		}

		current = append(toAdd, current...)
		current = collectFooters(child, current)
	}

	return current
}

func (m tui) HelpBar() string {
	binds := collectFooters(m, m.FooterKeys())
	return m.help.ShortHelpView(binds)
}

func (m tui) FullHelpView(current keyProvider, view string) string {

	// Add own key help to the thing
	view += m.renderKeyHelp(current.FullKeyHelp())

	// Add all of the children help views
	for _, child := range current.Children() {
		view = m.FullHelpView(child, view)
	}

	return view
}

func (m tui) renderKeyHelp(help FullKeyHelp) string {
	if help.Keys == nil {
		return ""
	}

	var b strings.Builder
	b.WriteString(titleStyle.Render(help.Title))
	b.WriteString("\n")
	b.WriteString(m.help.FullHelpView(help.Keys))
	return b.String() + "\n\n"
}

func (m tui) fullHelpHint() string {
	hintLabel := lipgloss.NewStyle().Bold(true).Foreground(textColor).Render("Hint:")
	sec := secondaryTextStyle.Render
	normal := textStyle.Render
	return hintLabel + " " + sec("Instead of ") + normal(SymbolArrowUp) + sec(" and ") + normal(SymbolArrowDown) + sec(", you can use ") + normal("shift+tab") + sec(" and ") + normal("tab") + sec(", respectively. ") + normal("Vim keybinds") + sec(" also work as you would expect.")
}
