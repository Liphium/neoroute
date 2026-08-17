package tui

import (
	"maps"
	"slices"
	"strings"

	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"github.com/Liphium/neoroute/neoschema"
	"github.com/Liphium/neoroute/pkg/neodebug/connector"
)

type inputFieldState int

// All of the states for the debug UI's input field thingy
const (
	// When you choose a route from the list of available ones.
	stateRouteSelect inputFieldState = iota

	// When you are viewing the entire request in the view
	stateCreateRequest

	// When you edit an individual field of the request
	stateEditRequest

	// When the little spinner for connecting is shown
	stateConnecting

	// When the connection has been closed
	stateClosed
)

type inputRequestHeightMsg struct {
	height int
}

type Input struct {
	handledKey bool
	width      int
	height     int
	state      inputFieldState

	// Connected / closed state
	connection connector.Connection
	spinner    spinner.Model

	// Route select
	routeSelect inputRouteSelect
}

func newInput(schema neoschema.TransporterSchema) Input {
	s := spinner.New(spinner.WithSpinner(spinner.Dot))

	return Input{
		state:       stateConnecting,
		spinner:     s,
		routeSelect: newRouteSelect(slices.Collect(maps.Keys(schema.Routes))),
	}
}

func (m *Input) SetWidth(w int) {
	m.width = w
	m.routeSelect.SetWidth(w)
}

func (m *Input) SetHeight(h int) {
	m.height = h
	m.routeSelect.SetHeight(h)
}

func (m Input) WantedHeight() int {
	switch m.state {
	case stateConnecting, stateClosed:
		return 1
	case stateRouteSelect:
		return m.routeSelect.WantedHeight()
	}
	return 5
}

func (m Input) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m Input) Update(msg tea.Msg) (Input, tea.Cmd) {
	m.handledKey = false

	differentHandling := false
	switch msg := msg.(type) {
	case spinner.TickMsg:
		differentHandling = true

		// Only update for connecting state (the only state we render in here)
		if m.state == stateConnecting {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}

	case connector.ConnectedMsg:
		differentHandling = true
		m.state = stateRouteSelect
		m.connection = msg.Connection
		return m, m.connection.WaitForEvent()

	case connector.ClosedMsg:
		differentHandling = true
		m.state = stateClosed
		return m, nil

	case connector.DoWaitMsg:
		differentHandling = true
		return m, m.connection.WaitForEvent()
	}

	// Forward any msgs not handled here down
	if !differentHandling {
		switch m.state {
		case stateRouteSelect:
			var cmd tea.Cmd
			m.routeSelect, cmd = m.routeSelect.Update(msg)
			m.handledKey = m.routeSelect.handledKey
			return m, cmd
		}
	}

	return m, nil
}

// Returns the cursor as well, when an input field is focused.
func (m Input) View() (*tea.Cursor, string) {
	if m.height < m.WantedHeight() {
		return nil, "too small to fit"
	}

	switch m.state {
	case stateConnecting:
		connectText := highlightStyle.Render(m.spinner.View()) + textStyle.Render("Connecting to transporter...")
		fill := strings.TrimSuffix(strings.Repeat("\n", m.height), "\n")
		return nil, connectText + fill

	case stateClosed:
		closeText := textStyle.Render("Connection closed.")
		fill := strings.TrimSuffix(strings.Repeat("\n", m.height), "\n")
		return nil, closeText + fill

	case stateRouteSelect:
		cursor, view := m.routeSelect.View()
		return cursor, view
	}

	return nil, strings.TrimSuffix(strings.Repeat(textStyle.Render("to do")+"\n", m.height), "\n")
}
