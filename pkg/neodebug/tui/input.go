package tui

import (
	"strings"

	tea "charm.land/bubbletea/v2"
)

type inputFieldState int

// All of the states for the debug UI's input field thingy
const (
	// When you choose a handler from the list of available ones.
	stateSelectHandler inputFieldState = iota

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
}

func newInput() Input {
	return Input{}
}

func (m *Input) SetWidth(w int) {
	m.width = w
}

func (m *Input) SetHeight(h int) {
	m.height = h
}

func (m Input) WantedHeight(windowHeight int) int {
	return 5
}

func (m Input) Update(msg tea.Msg) (Input, tea.Cmd) {
	m.handledKey = false

	return m, nil
}

// Returns the cursor as well, when an input field is focused.
func (m Input) View() (*tea.Cursor, string) {
	if m.height < 5 /* divider + 4 for input + any content */ {
		return nil, "too small to fit"
	}

	return nil, strings.TrimSuffix(strings.Repeat(symbolStyle.Render("i")+"\n", m.height), "\n")
}
