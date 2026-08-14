package tui

import (
	"github.com/Liphium/neoroute/neoschema"
	"github.com/Liphium/neoroute/pkg/neodebug/connector"
	tea "github.com/charmbracelet/bubbletea"
)

// Idea for UI:
// - Always visible history of all of the things that have received / sent (events, requests, responses)
// 	- v to get the history into full view
// 	- different sections (the individual)
// - Otherwise the input is focused where you can send stuff from (with all of the states below)
// 	- selectHandler: Search a handler in a 3 line high text input
//  - stateCreateRequest: Create a request a JSON like editor that parsed the complete response for the thing we're trying to send
// 	- stateEditRequest: When you click enter on one of the fields in the JSON like thingy, this is where you edit it
// - Bottom bar always constant with some information that we can get from other places (like hotkeys, etc.)

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

type sectionType int

const (
	sectionReceived sectionType = iota
	sectionSent
	sectionSystem
)

type Content struct {
	Error bool
	Line  string
}
type content = Content

var _ tea.Model = model{}

type model struct {
	viewingHistory bool // If we're currently traversing the history
	transporter    neoschema.TransporterSchema

	// Input field
	currentState inputFieldState

	// History
	received []Section
	selected int

	// Viewport things
	width int
}

func Run(transporter neoschema.TransporterSchema) *model {
	return &model{
		viewingHistory: false,
		transporter:    transporter,
		currentState:   stateConnecting,
		received: []Section{
			{
				T:       sectionSystem,
				Loading: true,
				Label:   "Initialization",
				Data: []content{
					{Line: "Starting neodebug..."},
				},
			},
		},
		width: 20,
	}
}

func (m model) Init() tea.Cmd {
	return connector.Connect(m.transporter)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return ""
}
