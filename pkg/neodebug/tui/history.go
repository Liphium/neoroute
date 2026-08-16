package tui

import (
	"strings"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/Liphium/neoroute/pkg/neodebug/model"
)

type ToggleSnapMsg struct{}

var (
	infoStyle   = lipgloss.NewStyle().Foreground(textColor).Bold(true)
	errorStyle  = lipgloss.NewStyle().Foreground(errorColor).Bold(true)
	symbolStyle = lipgloss.NewStyle().Foreground(successColor).Bold(true)
	eventStyle  = lipgloss.NewStyle().Foreground(highlightColor).Bold(true)
)

type History struct {
	handledKey   bool
	snapToBottom bool
	content      []model.Content
	viewport     viewport.Model
}

func NewHistory(w, h int) History {
	v := viewport.New(viewport.WithWidth(w), viewport.WithHeight(h))
	v.SoftWrap = false
	v.MouseWheelEnabled = true
	v.FillHeight = true

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
	// TODO: Determine if this enough (maybe we need to scroll down more?)
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
		b.WriteString(secondaryTextStyle.Render(c.Creation().Format("03:04 PM") + " "))

		switch c := c.(type) {
		case model.InfoContent:
			b.WriteString(infoStyle.Render("INF"))
			b.WriteRune(' ')
			b.WriteString(textStyle.Render(c.Message))

		case model.ErrorContent:
			b.WriteString(errorStyle.Render("ERR"))
			b.WriteRune(' ')
			b.WriteString(textStyle.Render(c.Message))

		case model.SendingContent:
			b.WriteString(symbolStyle.Render(SymbolArrowRight))
			b.WriteRune(' ')
			b.WriteString(textStyle.Render(c.Route))

		case model.ResponseContent:
			b.WriteString(symbolStyle.Render(SymbolArrowLeft))
			b.WriteRune(' ')
			b.WriteString(textStyle.Render(c.Route))

		case model.EventContent:
			b.WriteString(eventStyle.Render(SymbolArrowLeft))
			b.WriteRune(' ')
			b.WriteString(textStyle.Render(c.Name))
		}

		b.WriteString("\n")
	}
	return strings.TrimSuffix(b.String(), "\n")
}

func (m History) View() string { return m.viewport.View() }

func (m History) ViewportKeyMap() viewport.KeyMap { return m.viewport.KeyMap }
