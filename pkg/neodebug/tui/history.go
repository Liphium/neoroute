package tui

import (
	"strings"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/Liphium/neoroute/pkg/neodebug/model"
)

type ToggleSnapMsg struct{}

var (
	infoStyle          = lipgloss.NewStyle().Foreground(textColor).Bold(true)
	errorStyle         = lipgloss.NewStyle().Foreground(errorColor).Bold(true)
	warnStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Bold(true) // yellow / warn
	symbolStyle        = lipgloss.NewStyle().Foreground(successColor).Bold(true)
	eventStyle         = lipgloss.NewStyle().Foreground(highlightColor).Bold(true)
	newLineRenderStyle = lipgloss.NewStyle().Padding(0, 0, 0, 1)
)

var _ keyProvider = History{}

type History struct {
	handledKey   bool
	snapToBottom bool
	content      []model.Content
	viewport     viewport.Model
}

// Children implements keyProvider.
func (m History) Children() []keyProvider {
	return []keyProvider{}
}

// FooterKeys implements keyProvider.
func (m History) FooterKeys() []key.Binding {
	return []key.Binding{m.viewport.KeyMap.Up, m.viewport.KeyMap.Down}
}

// FullKeyHelp implements keyProvider.
func (m History) FullKeyHelp() FullKeyHelp {
	km := m.viewport.KeyMap
	return FullKeyHelp{
		Title: "History",
		Keys: [][]key.Binding{
			[]key.Binding{km.Up, km.Down},
			[]key.Binding{km.Right, km.Left},
		},
	}
}

func NewHistory(w, h int) History {
	v := viewport.New(viewport.WithWidth(w), viewport.WithHeight(h))
	v.SoftWrap = false
	v.MouseWheelEnabled = true
	v.FillHeight = true
	v.KeyMap.Up = key.NewBinding(standardUpKey, key.WithHelp("↑", "up"))
	v.KeyMap.Down = key.NewBinding(standardDownKey, key.WithHelp("↓", "down"))

	his := History{
		viewport: v,
		content: []model.Content{
			model.InfoContent{
				Message:      "Starting up...",
				BasicContent: model.Now(),
			},
		},
		snapToBottom: true,
	}
	his.viewport.SetContent(his.renderContent())

	return his
}

func (m *History) SetWidth(width int) {
	m.viewport.SetWidth(width)
}

func (m *History) SetHeight(height int) {
	m.viewport.SetHeight(height)
}

func (m *History) GotoBottom() {
	m.snapToBottom = true
	m.viewport.GotoBottom()
}

func (m History) Update(msg tea.Msg) (History, tea.Cmd) {
	m.handledKey = false
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case model.AddContentMsg:
		m.content = append(m.content, msg.Content)

		// Set new content and scroll it to the bottom when we snap
		m.viewport.SetContent(m.renderContent())
		if m.snapToBottom {
			m.viewport.GotoBottom()
		}
	}

	// Update the viewport, when it scrolls we want to determine if we actually should snap to the bottom
	oldOffset := m.viewport.YOffset()
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)
	if m.viewport.YOffset() != oldOffset {
		m.handledKey = true
		m.snapToBottom = m.viewport.AtBottom()
	}

	return m, tea.Batch(cmds...)
}

func (m History) renderContent() string {
	var b strings.Builder
	for _, c := range m.content {

		// Add the creation timestamp
		timestamp := secondaryTextStyle.Render(c.Creation().Format("03:04 PM") + " ")
		b.WriteString(timestamp)

		switch c := c.(type) {
		case model.InfoContent:
			b.WriteString(infoStyle.Render("INF"))
			b.WriteRune(' ')
			b.WriteString(textStyle.Render(c.Message))

		case model.ErrorContent:
			b.WriteString(errorStyle.Render("ERR"))
			b.WriteRune(' ')
			b.WriteString(textStyle.Render(c.Message))

		case model.WarnContent:
			b.WriteString(warnStyle.Render("WRN"))
			b.WriteRune(' ')
			b.WriteString(textStyle.Render(c.Message))

		case model.ResponseContent:

			// Write the first line for the request
			b.WriteString(symbolStyle.Render("REQ"))
			b.WriteRune(' ')
			msg := "Sent request to "
			b.WriteString(lessTextStyle.Render(msg) + textStyle.Bold(true).Render(c.Route))
			b.WriteRune(' ')
			b.WriteString(m.renderValue(13 /* timestamp and stuff */ +len(msg)+len(c.Route), c.Request, symbolStyle))
			b.WriteString("\n")

			// Write the second line for the response to the request
			b.WriteString(timestamp)
			b.WriteString(symbolStyle.Render("RES"))
			b.WriteRune(' ')
			msg = "Got response for "
			b.WriteString(lessTextStyle.Render(msg) + textStyle.Bold(true).Render(c.Route))
			b.WriteRune(' ')
			b.WriteString(m.renderValue(13+len(msg)+len(c.Route), c.Response, symbolStyle))

		case model.EventContent:
			b.WriteString(eventStyle.Render("EVT"))
			b.WriteRune(' ')
			msg := "Received event "
			b.WriteString(lessTextStyle.Render(msg) + textStyle.Bold(true).Render(c.Name))
			b.WriteRune(' ')
			b.WriteString(m.renderValue(13 /* timestamp and stuff */ +len(msg)+len(c.Name), c.Event, eventStyle))
		}

		b.WriteString("\n")
	}
	return strings.TrimSuffix(b.String(), "\n")
}

// Render the value, if oneline add behind, otherwise do line breaking
func (m History) renderValue(lenBefore int, value any, divider lipgloss.Style) string {
	req := render(value, m.viewport.Width()-lenBefore)
	if lipgloss.Height(req) == 1 {
		return textStyle.Render(req)
	} else {
		req = render(value, m.viewport.Width()-2 /* padding + color */)
		renderedDivider := divider.Render(SymbolHorizontalDivider)
		req = newLineRenderStyle.Render(textStyle.Render(req))
		return "\n" + renderedDivider + strings.ReplaceAll(req, "\n", "\n"+renderedDivider)
	}
}

func (m History) View() string { return m.viewport.View() }

func (m History) ViewportKeyMap() viewport.KeyMap { return m.viewport.KeyMap }
