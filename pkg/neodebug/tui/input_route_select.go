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

type RouteSelectedMsg struct {
	Route string
}

var _ keyProvider = inputRouteSelect{}

type inputRouteSelect struct {
	handledKey    bool
	width         int
	height        int
	input         textinput.Model
	routes        []string
	results       []string
	selectedRoute int

	// keys
	clearSearch key.Binding
	filter      key.Binding
	up          key.Binding
	down        key.Binding
	enter       key.Binding
}

// Children implements keyProvider.
func (m inputRouteSelect) Children() []keyProvider {
	return []keyProvider{}
}

// FooterKeys implements keyProvider.
func (m inputRouteSelect) FooterKeys() []key.Binding {
	if m.input.Focused() {
		return []key.Binding{m.enter, m.clearSearch, m.up, m.down}
	}
	return []key.Binding{m.enter, m.filter, m.up, m.down}
}

// FullKeyHelp implements keyProvider.
func (m inputRouteSelect) FullKeyHelp() FullKeyHelp {
	return FullKeyHelp{
		Title: "Route selection",
		Keys: [][]key.Binding{
			[]key.Binding{m.enter, m.clearSearch, m.filter},
			[]key.Binding{m.up, m.down},
		},
	}
}

func newRouteSelect(routes []string) inputRouteSelect {
	i := textinput.New()
	i.Placeholder = ""
	i.Prompt = "Filter: "
	styles := textinput.DefaultDarkStyles()
	styles.Blurred.Prompt, styles.Focused.Prompt = secondaryTextStyle, secondaryTextStyle
	styles.Blurred.Text, styles.Focused.Text = textStyle, textStyle
	i.SetStyles(styles)
	i.SetVirtualCursor(false)

	return inputRouteSelect{
		input:         i,
		routes:        routes,
		results:       routes,
		selectedRoute: 0,

		// Define default keys
		clearSearch: key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "clear search")),
		filter:      key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "search")),
		up:          key.NewBinding(standardUpKey, key.WithHelp("↑", "up")),
		down:        key.NewBinding(standardDownKey, key.WithHelp("↓", "down")),
		enter:       key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select")),
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
		case matchesFocusedNavigation(msg, m.input.Focused(), m.up):
			handled = true
			m.selectedRoute = max(0, m.selectedRoute-1)
		case matchesFocusedNavigation(msg, m.input.Focused(), m.down):
			handled = true
			m.selectedRoute = min(len(m.results)-1, m.selectedRoute+1)
		case key.Matches(msg, m.clearSearch):
			if !m.input.Focused() {
				break
			}
			handled = true
			m.clearSearchInput()
		case key.Matches(msg, m.enter):
			handled = true
			m.handledKey = true

			if len(m.results) == 0 {
				break
			}

			return m, func() tea.Msg {
				return RouteSelectedMsg{Route: m.results[m.selectedRoute]}
			}
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
		case key.Matches(msg, m.filter):
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

func (m *inputRouteSelect) clearSearchInput() {
	m.input.Blur()
	m.input.SetValue("")

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
