package tui

import (
	"math"
	"strings"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/Liphium/neoroute/neoschema"
	"github.com/Liphium/neoroute/pkg/neodebug/connector"
)

// Idea for UI:
// - Always visible history of all of the things that have received / sent (events, requests, responses)
//   - v to get the history into full view
//   - different sections (the individual)
// - Otherwise the input is focused where you can send stuff from (with all of the states below)
//   - selectHandler: Search a handler in a 3 line high text input
//   - stateCreateRequest: Create a request a JSON like editor that parsed the complete response for the thing we're trying to send
//   - stateEditRequest: When you click enter on one of the fields in the JSON like thingy, this is where you edit it
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

var _ tea.Model = model{}

// Basic colors that we're gonna use
var (
	textColor          = lipgloss.Color("252") // main text (more white)
	secondaryTextColor = lipgloss.Color("248") // secondary text (gray)

	separatorColor = lipgloss.Color("245") // gray for the separators (more gray than the text stuff)

	highlightColor = lipgloss.Color("38")  // blue
	errorColor     = lipgloss.Color("196") // red
	successColor   = lipgloss.Color("46")  // green
)

// Characters that are handy
const (
	SymbolDivider = "─"

	SymbolArrowRight = "→"
	SymbolArrowLeft  = "←"
	SymbolArrowUp    = "↑"
	SymbolArrowDown  = "↓"
)

// Styles for all of the elements
var (
	topPanelStyle    = lipgloss.NewStyle().Padding(0, 1)
	bottomPanelStyle = lipgloss.NewStyle().Padding(0, 1)

	titleStyle         = lipgloss.NewStyle().Bold(true).Foreground(highlightColor)
	textStyle          = lipgloss.NewStyle().Foreground(textColor)
	secondaryTextStyle = lipgloss.NewStyle().Foreground(secondaryTextColor)

	separatorStyle = lipgloss.NewStyle().Foreground(separatorColor)
)

type model struct {
	transporter neoschema.TransporterSchema

	// History
	history History

	// Viewport / layout
	width  int
	height int

	// Key bindings + help
	keys appKeyMap
	help help.Model
}

func Run(transporter neoschema.TransporterSchema) *model {
	h := help.New()
	h.ShowAll = false

	// Customize the short styles
	h.Styles.ShortKey = lipgloss.NewStyle().Foreground(textColor).Bold(true)
	h.Styles.ShortDesc = lipgloss.NewStyle().Foreground(secondaryTextColor)
	h.Styles.ShortSeparator = lipgloss.NewStyle().Foreground(separatorColor)
	h.Styles.Ellipsis = lipgloss.NewStyle().Foreground(separatorColor)

	// Customize the full help styles
	h.Styles.FullKey = lipgloss.NewStyle().Foreground(textColor).Bold(true)
	h.Styles.FullDesc = lipgloss.NewStyle().Foreground(secondaryTextColor)
	h.Styles.FullSeparator = lipgloss.NewStyle().Foreground(separatorColor)

	return &model{
		transporter: transporter,
		history:     NewHistory(0, 0),
		width:       0,
		height:      0,
		keys:        newAppKeyMap(),
		help:        h,
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
		m.height = msg.Height
		m.help.SetWidth(msg.Width)

	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		}
	}

	// Update the history
	h, cmd := m.history.Update(msg)
	cmds = append(cmds, cmd)
	m.history = h.(History)

	// TODO: Calculate new sizes and stuff + update history

	if len(cmds) == 0 {
		return m, nil
	}
	return m, tea.Batch(cmds...)
}

func (m model) View() tea.View {
	if m.width == 0 || m.height == 0 {
		return tea.NewView("loading neodebug...")
	}

	helpView := m.help.View(m.keys)
	remaining := m.height - lipgloss.Height(helpView) - 2 /* dividers */

	topH := int(math.Floor(float64(remaining) / 2))
	bottomH := int(math.Ceil(float64(remaining) / 2))

	top := topPanelStyle.Width(m.width).Height(topH).Render(m.renderTopPanel())
	bottom := bottomPanelStyle.Width(m.width).Height(bottomH).Render(m.renderBottomPanel())
	divider := separatorStyle.Render(strings.Repeat(SymbolDivider, m.width))

	v := tea.NewView(lipgloss.JoinVertical(lipgloss.Left, top, divider, bottom, divider, helpView))
	v.AltScreen = true
	return v
}

func (m model) renderTopPanel() string {
	title := titleStyle.Render("History") + lipgloss.NewStyle().Faint(true).Render(" (h to focus)")
	return lipgloss.JoinVertical(lipgloss.Left, title)
}

func (m model) renderBottomPanel() string {
	title := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("240")).Render("Bottom panel")
	return lipgloss.JoinVertical(lipgloss.Left, title)
}
