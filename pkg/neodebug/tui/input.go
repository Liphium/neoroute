package tui

import (
	"strings"

	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
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
)

type inputRequestHeightMsg struct {
	height int
}

type Input struct {
	handledKey bool
	width      int
	height     int
	state      inputFieldState
	spinner    spinner.Model
}

func newInput() Input {
	s := spinner.New(spinner.WithSpinner(spinner.Dot))

	return Input{
		state:   stateConnecting,
		spinner: s,
	}
}

func (m *Input) SetWidth(w int) {
	m.width = w
}

func (m *Input) SetHeight(h int) {
	m.height = h
}

// TODO: For dynamic height changing, diff this in TUI and when it changes change the height of the input.
//
// Also TODO: Make this return what the selected dialog wants.
func (m Input) WantedHeight() int {
	switch m.state {
	case stateConnecting:
		return 1
	case stateRouteSelect:
		return 5
	}
	return 5
}

func (m Input) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m Input) Update(msg tea.Msg) (Input, tea.Cmd) {
	m.handledKey = false

	switch msg := msg.(type) {
	case spinner.TickMsg:
		// Only update for connecting state (the only state we render in here)
		if m.state == stateConnecting {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
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
	}

	return nil, strings.TrimSuffix(strings.Repeat(symbolStyle.Render("i")+"\n", m.height), "\n")
}
