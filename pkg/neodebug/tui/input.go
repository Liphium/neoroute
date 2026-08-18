package tui

import (
	"maps"
	"slices"
	"strings"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"github.com/Liphium/neoroute/neoschema"
	"github.com/Liphium/neoroute/pkg/neodebug/connector"
	"github.com/Liphium/neoroute/pkg/neodebug/model"
)

type inputFieldState int

// Wishlist:
// - Add reconnections

// All of the states for the debug UI's input field thingy
const (
	// When you choose a route from the list of available ones.
	stateRouteSelect inputFieldState = iota

	// When you are viewing the entire request in the view
	stateCreateRequest

	// When the little spinner for connecting is shown
	stateConnecting

	// When the connection has been closed
	stateClosed
)

type inputRequestHeightMsg struct {
	height int
}

var _ keyProvider = Input{}

type Input struct {
	handledKey bool
	width      int
	height     int
	state      inputFieldState
	schema     neoschema.TransporterSchema

	// Connected / closed state
	connection connector.Connection
	spinner    spinner.Model

	// Route select
	routeSelect inputRouteSelect

	// Request creator
	requestCreator inputRequestCreator
}

// Children implements keyProvider.
func (m Input) Children() []keyProvider {
	switch m.state {
	case stateRouteSelect:
		return []keyProvider{m.routeSelect}
	case stateCreateRequest:
		return []keyProvider{m.requestCreator}
	}

	return []keyProvider{}
}

// FooterKeys implements keyProvider.
func (m Input) FooterKeys() []key.Binding {
	return []key.Binding{}
}

// FullKeyHelp implements keyProvider.
func (m Input) FullKeyHelp() FullKeyHelp {
	return FullKeyHelp{}
}

func newInput(schema neoschema.TransporterSchema) Input {
	s := spinner.New(spinner.WithSpinner(spinner.Dot))

	return Input{
		state:       stateConnecting,
		spinner:     s,
		schema:      schema,
		routeSelect: newRouteSelect(slices.Collect(maps.Keys(schema.Routes))),
	}
}

func (m *Input) SetWidth(w int) {
	m.width = w
	m.routeSelect.SetWidth(w)
	m.requestCreator.SetWidth(w)
}

func (m *Input) SetHeight(h int) {
	m.height = h
	m.routeSelect.SetHeight(h)
	m.requestCreator.SetHeight(h)
}

func (m Input) WantedHeight() int {
	switch m.state {
	case stateConnecting, stateClosed:
		return 1
	case stateRouteSelect:
		return m.routeSelect.WantedHeight()
	case stateCreateRequest:
		return m.requestCreator.WantedHeight()
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

	case RouteSelectedMsg:
		differentHandling = true

		// When selected, switch to editing the input
		route, ok := m.schema.Routes[msg.Route]
		if !ok {
			return m, model.Plain(model.Error("Couldn't find route selected."))
		}
		if !route.HasRequest {
			return m, model.Plain(model.Info("Coming soon..."))
		}

		// Switch to new creation state
		m.state = stateCreateRequest
		m.requestCreator = newInputRequestCreator(msg.Route, route.Request, m.width, m.height)

		return m, nil

	case SendMsg:
		differentHandling = true

		// Switch to state and handle cancellation properly
		m.state = stateRouteSelect
		if msg.Cancelled {
			return m, nil
		}

		// When sent, switch to the route selector again
		_, ok := m.schema.Routes[msg.Route]
		if !ok {
			return m, model.Plain(model.Error("Couldn't find route selected."))
		}

		// TODO: Actually send the route

		return m, model.Plain(model.Info("Would send route now, but that's kinda not implemented ig"))
	}

	// Forward any msgs not handled here down
	if !differentHandling {
		switch m.state {
		case stateRouteSelect:
			var cmd tea.Cmd
			m.routeSelect, cmd = m.routeSelect.Update(msg)
			m.handledKey = m.routeSelect.handledKey
			return m, cmd
		case stateCreateRequest:
			var cmd tea.Cmd
			m.requestCreator, cmd = m.requestCreator.Update(msg)
			m.handledKey = true // Always give this back cause other than priority keys everything should be handled here
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

	case stateCreateRequest:
		cursor, view := m.requestCreator.View()
		return cursor, view
	}

	return nil, strings.TrimSuffix(strings.Repeat(textStyle.Render("to do")+"\n", m.height), "\n")
}
