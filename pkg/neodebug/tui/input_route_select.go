package tui

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/Liphium/neoroute/pkg/neodebug/v2/config"
	"github.com/Liphium/neoroute/pkg/neodebug/v2/model"
	"github.com/Liphium/neoroute/v2/neoschema"
	"github.com/charmbracelet/x/ansi"
	"github.com/sahilm/fuzzy"
)

var (
	selectedStyle   = textStyle.Bold(true)
	unselectedStyle = secondaryTextStyle
)

type RouteSelectedMsg struct {
	Route  string
	Preset *config.Preset
	Edit   bool
}

var _ keyProvider = inputRouteSelect{}

type inputRouteSelect struct {
	handledKey    bool
	width         int
	height        int
	input         textinput.Model
	routes        []routeItem
	results       []routeItem
	selectedRoute int

	// keys
	clearSearch key.Binding
	filter      key.Binding
	up          key.Binding
	down        key.Binding
	enter       key.Binding
	edit        key.Binding
}

type routeItem struct {
	Name    string
	Warning string // We only want to warn about invalid presets since having wrong ones is also a good way to test neoroute
	Preset  *config.Preset
}

func (r routeItem) Route() string {
	if r.Preset != nil {
		return r.Preset.Route
	}
	return r.Name
}

// Children implements keyProvider.
func (m inputRouteSelect) Children() []keyProvider {
	return []keyProvider{}
}

// FooterKeys implements keyProvider.
func (m inputRouteSelect) selectionKeys() []key.Binding {
	enter, edit := m.enter, m.edit
	enter.SetEnabled(len(m.results) > 0)
	edit.SetEnabled(false)
	if m.selectedRoute >= 0 && m.selectedRoute < len(m.results) {
		item := m.results[m.selectedRoute]
		if item.Preset != nil {
			enter.SetHelp("enter", "send")
			edit.SetEnabled(item.Warning == "")
		}
	}
	return []key.Binding{enter, edit}
}

func (m inputRouteSelect) FooterKeys() []key.Binding {
	keys := m.selectionKeys()
	if m.input.Focused() {
		return append(keys, m.clearSearch, m.up, m.down)
	}
	return append(keys, m.filter, m.up, m.down)
}

func (m inputRouteSelect) FullKeyHelp() FullKeyHelp {
	return FullKeyHelp{
		Title: "Route selection",
		Keys: [][]key.Binding{
			append(m.selectionKeys(), m.clearSearch, m.filter),
			{m.up, m.down},
		},
	}
}

func newRouteSelect(routes map[string]neoschema.RouteSchema) inputRouteSelect {
	i := textinput.New()
	i.Placeholder = ""
	i.Prompt = "Filter: "
	styles := textinput.DefaultDarkStyles()
	styles.Blurred.Prompt, styles.Focused.Prompt = secondaryTextStyle, secondaryTextStyle
	styles.Blurred.Text, styles.Focused.Text = textStyle, textStyle
	i.SetStyles(styles)
	i.SetVirtualCursor(false)

	// Add all presets to the routes
	items := []routeItem{}
	for name, _ := range routes {
		items = append(items, routeItem{Name: name})
	}
	for _, preset := range config.Config.Presets {
		items = append(items, routeItem{
			Name:    fmt.Sprintf("Preset: %s - %s", preset.Route, preset.Name),
			Preset:  &preset,
			Warning: preset.Check(routes),
		})
	}

	return inputRouteSelect{
		input:         i,
		routes:        items,
		results:       items,
		selectedRoute: 0,

		// Define default keys
		clearSearch: key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "clear search")),
		filter:      key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "search")),
		up:          key.NewBinding(standardUpKey, key.WithHelp("↑", "up")),
		down:        key.NewBinding(standardDownKey, key.WithHelp("↓", "down")),
		edit:        key.NewBinding(key.WithKeys("alt+e"), key.WithHelp("alt+e", "edit")),
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
		case key.Matches(msg, m.enter), key.Matches(msg, m.edit):
			handled = true
			m.handledKey = true

			if len(m.results) == 0 {
				break
			}

			item := m.results[m.selectedRoute]
			edit := key.Matches(msg, m.edit)
			if edit && (item.Preset == nil || item.Warning != "") {
				break
			}
			return m, model.Plain(RouteSelectedMsg{Route: item.Route(), Preset: item.Preset, Edit: edit})
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
	var selected *routeItem
	if m.selectedRoute >= 0 && m.selectedRoute < len(m.results) {
		item := m.results[m.selectedRoute]
		selected = &item
	}
	m.input.Blur()
	m.input.SetValue("")

	m.results = m.routes
	m.selectedRoute = 0
	if selected != nil {
		for i, item := range m.results {
			if item == *selected {
				m.selectedRoute = i
				break
			}
		}
	}
}

func (m *inputRouteSelect) updateSearch(query string) {
	m.selectedRoute = 0
	if query == "" {
		m.results = m.routes
		return
	}

	names := make([]string, len(m.routes))
	for i, route := range m.routes {
		names[i] = route.Name
	}

	matches := fuzzy.Find(query, names)
	m.results = make([]routeItem, 0, len(matches))
	for _, match := range matches {
		m.results = append(m.results, m.routes[match.Index])
	}
}

func (m inputRouteSelect) View() (*tea.Cursor, string) {
	var cursor *tea.Cursor

	// Render the title (when searching just the input)
	titleLine := titleStyle.Render("Select a route or preset") + " " + secondaryTextStyle.Render("(/ to filter)")
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

		warning := ""
		if result.Warning != "" {
			// Reserve the row padding, arrow, space, and warning parentheses.
			available := m.width - 2 - ansi.StringWidth(SymbolArrowRight+" "+result.Name)
			if available >= 6 {
				warning = secondaryTextStyle.Render(" (" + ansi.Truncate(result.Warning, available-3, "...") + ")")
			}
		}

		if i == m.selectedRoute {
			b.WriteString(resultsStyle.Render(highlightStyle.Render(SymbolArrowRight) + " " + selectedStyle.Render(result.Name) + warning))
		} else {
			b.WriteString(resultsStyle.Render(secondaryTextStyle.Render(SymbolArrowRight) + " " + unselectedStyle.Render(result.Name) + warning))
		}

		b.WriteString("\n")
	}

	fill := strings.TrimSuffix(strings.Repeat("\n", m.height-1-(end-start)), "\n")
	return cursor, strings.TrimSuffix(titleLine+"\n"+b.String()+fill, "\n")
}
