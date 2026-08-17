package tui

import (
	"strings"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var (
	selectedStyle   = textStyle.Bold(true)
	unselectedStyle = secondaryTextStyle
)

type routeSelectKeyMap struct {
	ClearSearch key.Binding
	Filter      key.Binding
	Up          key.Binding
	Down        key.Binding
}

func routeSelectDefaultKeyMap() routeSelectKeyMap {
	return routeSelectKeyMap{
		ClearSearch: key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "clear search")),
		Filter:      key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "filter routes")),
		Up:          key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "select previous route")),
		Down:        key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "select next route")),
	}
}

type inputRouteSelect struct {
	routeSelectKeyMap

	handledKey    bool
	width         int
	height        int
	input         textinput.Model
	routes        []string
	results       []string
	selectedRoute int
}

func newRouteSelect(routes []string) inputRouteSelect {
	i := textinput.New()
	i.Placeholder = ""
	i.Prompt = secondaryTextStyle.Render("Filter: ")
	i.SetStyles(textinput.DefaultDarkStyles())
	i.SetVirtualCursor(false)

	return inputRouteSelect{
		routeSelectKeyMap: routeSelectDefaultKeyMap(),
		input:             i,
		routes:            routes,
		results:           routes,
		selectedRoute:     0,
	}
}

func (m *inputRouteSelect) SetWidth(w int) {
	m.width = w
}

func (m *inputRouteSelect) SetHeight(h int) {
	m.height = h
}

func (m inputRouteSelect) WantedHeight() int {
	return 1 + min(len(m.results), 3)
}

func (m inputRouteSelect) Update(msg tea.Msg) (inputRouteSelect, tea.Cmd) {
	m.handledKey = false
	var cmd tea.Cmd

	differentHandling := false
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		differentHandling = true

		// Handle the keys that are specific to the input
		handled := false
		switch {
		case key.Matches(msg, m.Up):
			handled = true
			m.selectedRoute = max(0, m.selectedRoute-1)
		case key.Matches(msg, m.Down):
			handled = true
			m.selectedRoute = min(len(m.results)-1, m.selectedRoute+1)
		case key.Matches(msg, m.ClearSearch):
			if !m.input.Focused() {
				break
			}
			handled = true
			m.clearSearch()
		}
		m.handledKey = handled

		// Make sure the input does not get any keys we handle here
		if handled {
			break
		}

		// When the input is focused, let it handle all of the things
		if m.input.Focused() {
			m.input, cmd = m.input.Update(msg)
			m.updateSearch(m.input.Value())
			m.handledKey = true
			break
		}

		// Handle the filter key (we don't want / to be up there because we want to support it being typed)
		switch {
		case key.Matches(msg, m.Filter):
			m.handledKey = true
			return m, m.input.Focus()
		}
	}

	// Forward other events to the textinput anyway
	if !differentHandling {
		m.input, cmd = m.input.Update(msg)
	}

	return m, cmd
}

func (m *inputRouteSelect) clearSearch() {
	m.input.Blur()
	m.input.SetValue("")

	// TODO: Determine if this is a good idea or not
	/*
		// Restore the selection to the correct route
		if len(m.results) != 0 {
			currentlySelected := m.results[m.selectedRoute]
			idx := slices.Index(m.routes, currentlySelected)
			if idx != -1 {
				m.selectedRoute = idx
			}
		}
	*/

	m.results = m.routes
	m.selectedRoute = 0
}

func (m *inputRouteSelect) updateSearch(query string) {
	m.selectedRoute = 0
	if query == "" {
		m.results = m.routes
		return
	}

	m.results = []string{}
	for _, route := range m.routes {
		if strings.Contains(route, query) {
			m.results = append(m.results, route)
		}
	}
}

func (m inputRouteSelect) View() (*tea.Cursor, string) {
	var cursor *tea.Cursor

	// Render the title (when searching just the input)
	titleLine := titleStyle.Render("Select a route") + " " + secondaryTextStyle.Render("(/ to filter)")
	if m.input.Focused() {
		titleLine = m.input.View()
		cursor = m.input.Cursor()
	}

	// Render all of the results
	resultsStyle := lipgloss.NewStyle().Width(m.width).Padding(0, 1)
	maxShown := m.height - 1

	n := len(m.results)

	var start, end int
	if n <= maxShown {
		start, end = 0, n-1 // all visible
	} else {
		half := maxShown / 2
		start = m.selectedRoute - half
		start = max(0, min(start, n-maxShown)) // clamp → pin start/end
		end = start + maxShown - 1             // inclusive, window size = maxShown
	}

	var b strings.Builder
	for i := start; i <= end; i++ {
		result := m.results[i]
		if i == m.selectedRoute {
			b.WriteString(resultsStyle.Render(highlightStyle.Render(SymbolArrowRight) + " " + selectedStyle.Render(result)))
		} else {
			b.WriteString(resultsStyle.Render(secondaryTextStyle.Render(SymbolArrowRight) + " " + unselectedStyle.Render(result)))
		}

		b.WriteString("\n")
	}

	fill := strings.TrimSuffix(strings.Repeat("\n", m.height-1-(end-start)), "\n")
	return cursor, strings.TrimSuffix(titleLine+"\n"+b.String()+fill, "\n")
}
