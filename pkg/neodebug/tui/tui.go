package tui

import (
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

var _ tea.Model = tui{}

// Basic colors that we're gonna use
var (
	textColor          = lipgloss.Color("252") // main text
	secondaryTextColor = lipgloss.Color("246") // secondary text

	separatorColor = lipgloss.Color("240") // subtle separators

	highlightColor = lipgloss.Color("67")  // muted blue
	errorColor     = lipgloss.Color("167") // muted red
	successColor   = lipgloss.Color("71")  // muted green
	warningColor   = lipgloss.Color("214") // yellow / warn
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
	highlightStyle     = lipgloss.NewStyle().Foreground(highlightColor)

	separatorStyle = lipgloss.NewStyle().Foreground(separatorColor)
)

type fullScreenView string

const (
	fullScreenNone    fullScreenView = "none"
	fullScreenHistory fullScreenView = "history"
	fullScreenInput   fullScreenView = "input"
)

func renderDivider(w int) string {
	return separatorStyle.Render(strings.Repeat(SymbolDivider, w))
}

type tui struct {
	transporter neoschema.TransporterSchema

	full    fullScreenView
	input   Input
	history History

	// Viewport / layout
	width    int
	height   int
	tooSmall bool

	// Key bindings + help
	keys appKeyMap
	help help.Model
}

func Run(transporter neoschema.TransporterSchema) *tui {
	h := help.New()
	h.ShowAll = false

	// Customize the short styles
	h.Styles.ShortKey = lipgloss.NewStyle().Foreground(textColor)
	h.Styles.ShortDesc = lipgloss.NewStyle().Foreground(secondaryTextColor)
	h.Styles.ShortSeparator = lipgloss.NewStyle().Foreground(separatorColor)
	h.Styles.Ellipsis = lipgloss.NewStyle().Foreground(separatorColor)

	// Customize the full help styles
	h.Styles.FullKey = lipgloss.NewStyle().Foreground(textColor)
	h.Styles.FullDesc = lipgloss.NewStyle().Foreground(secondaryTextColor)
	h.Styles.FullSeparator = lipgloss.NewStyle().Foreground(separatorColor)

	return &tui{
		transporter: transporter,
		full:        fullScreenNone,
		input:       newInput(transporter),
		history:     NewHistory(0, 0),
		width:       0,
		height:      0,
		keys:        newAppKeyMap(),
		help:        h,
	}
}

func (m tui) Init() tea.Cmd {
	return tea.Batch(connector.Connect(m.transporter), m.input.Init())
}

func (m tui) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd // Possibly a command that we want to return from the children
	var cmds []tea.Cmd

	var handled, differentHandling bool
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		handled = true
		m.width = msg.Width
		m.height = msg.Height
		m.help.SetWidth(msg.Width)
		m.input.SetWidth(msg.Width)
		m.history.SetWidth(msg.Width)
		m.relayoutHeight() // This only handles height

	case tea.MouseWheelMsg:
		differentHandling = true

		historyVisible := msg.Mouse().Y <= m.history.viewport.Height()
		if historyVisible {
			m.history, _ = m.history.Update(msg)
		}

	case tea.KeyPressMsg:
		differentHandling = true

		// Let the children handle keys first, when they handled them, we don't
		if m.inputVisible() {
			cmd = m.inputHandleMsg(msg)
			if m.input.handledKey {
				return m, cmd
			}
		}
		if m.historyVisible() {
			m.history, cmd = m.history.Update(msg)
			if m.history.handledKey {
				return m, cmd
			}
		}

		switch {
		case key.Matches(msg, m.keys.Quit):
			handled = true
			return m, tea.Quit

		case key.Matches(msg, m.keys.Help):
			handled = true
			m.help.ShowAll = !m.help.ShowAll
			m.relayoutHeight()

		case key.Matches(msg, m.keys.ExpandHistory):
			handled = true
			if m.full == fullScreenHistory {
				m.setFull(fullScreenNone)
			} else {
				m.setFull(fullScreenHistory)
			}

		case key.Matches(msg, m.keys.ExpandInput):
			handled = true
			if m.full == fullScreenInput {
				m.setFull(fullScreenNone)
			} else {
				m.setFull(fullScreenInput)
			}

		case key.Matches(msg, m.keys.ExitFullscreen):
			handled = true
			m.setFull(fullScreenNone)
		}
	}

	// If we've already handled the event, we don't want our children to handle it as well (we like already gave out events anyway)
	if !handled && !differentHandling {
		cmds = append(cmds, m.inputHandleMsg(msg))

		m.history, cmd = m.history.Update(msg)
		cmds = append(cmds, cmd)
	}

	if len(cmds) == 0 {
		return m, nil
	}
	return m, tea.Batch(cmds...)
}

func (m *tui) inputHandleMsg(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	// When the height of the input changes, we want to immediately want to relayout ourselves
	prev := m.input.WantedHeight()
	m.input, cmd = m.input.Update(msg)
	if prev != m.input.WantedHeight() {
		m.relayoutHeight()
	}

	return cmd
}

func (m *tui) inputVisible() bool {
	return m.full != fullScreenHistory
}

func (m *tui) historyVisible() bool {
	return m.full != fullScreenInput
}

func (m *tui) relayoutHeight() {

	// When we're in fullscreen, actually adjust differently
	if m.full != fullScreenNone {
		m.setFull(m.full)
		return
	}

	inputHeight := m.input.WantedHeight()
	historyHeight := m.height - inputHeight - 3 /* help bar + dividers */
	if historyHeight <= 0 {
		m.tooSmall = true
		return
	}

	m.history.SetHeight(historyHeight)
	m.input.SetHeight(inputHeight)
}

func (m *tui) setFull(full fullScreenView) {
	m.full = full
	switch m.full {
	case fullScreenNone:
		m.relayoutHeight()
	case fullScreenInput:
		m.input.SetHeight(m.height - 2 /* help view + divider */)
	case fullScreenHistory:
		m.history.SetHeight(m.height - 2 /* help view + divider */)
	}
}

func (m tui) View() tea.View {

	// Configure the main view
	view := tea.NewView("")
	view.MouseMode = tea.MouseModeCellMotion
	view.AltScreen = true

	if m.width == 0 || m.height == 0 || m.tooSmall {
		view.SetContent("loading neodebug...")
		return view
	}

	// All of the hotkeys are shown on a different page
	if m.help.ShowAll {
		view.SetContent("all of your keys")
		return view
	}

	divider := renderDivider(m.width)
	helpView := m.help.View(m.keys)

	if m.full != fullScreenNone {

		switch m.full {
		case fullScreenHistory:
			history := m.history.View()
			view.SetContent(lipgloss.JoinVertical(lipgloss.Left, history, divider, helpView))
			return view

		case fullScreenInput:
			cursor, input := m.input.View()
			view.Cursor = cursor
			view.SetContent(lipgloss.JoinVertical(lipgloss.Left, input, divider, helpView))
			return view
		}
	}

	// Render normally with history and input
	history := m.history.View()
	_, input := m.input.View()

	view.SetContent(lipgloss.JoinVertical(lipgloss.Left, history, divider, input, divider, helpView))
	return view
}
