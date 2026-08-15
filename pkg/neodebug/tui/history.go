package tui

import (
	"strings"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type Content interface{}

type ErrorContent struct {
	Message string
}

type SendingContent struct {
	Route string
}

type EventContent struct {
	Name  string
	Event any
}

type ResponseContent struct {
	Route string
	Data  any
}

type AddContentMsg struct{ Content }
type ToggleSnapMsg struct{}

var (
	errorStyle  = lipgloss.NewStyle().Foreground(errorColor).Bold(true)
	symbolStyle = lipgloss.NewStyle().Foreground(successColor).Bold(true)
	eventStyle  = lipgloss.NewStyle().Foreground(highlightColor).Bold(true)
)

type History struct {
	snapToBottom bool
	content      []Content
	viewport     viewport.Model
}

func NewHistory(w, h int) History {
	return History{
		viewport:     viewport.New(viewport.WithWidth(w), viewport.WithHeight(h)),
		snapToBottom: true,
	}
}

func (h History) Init() tea.Cmd { return nil }

func (h History) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		h.viewport.SetWidth(msg.Width)
		h.viewport.SetHeight(msg.Height)
		h.viewport.SetContent(h.renderContent())

	case AddContentMsg:
		c := msg.Content
		h.content = append(h.content, c)
		h.refreshViewport()
		// TODO: Make sure we don't scroll the viewport to the bottom when the new thingy is added

	case ToggleSnapMsg:
		h.toggleSnap()
	}

	oldOffset := h.viewport.YOffset()
	var cmd tea.Cmd
	h.viewport, cmd = h.viewport.Update(msg)
	cmds = append(cmds, cmd)
	if h.viewport.YOffset() != oldOffset {
		h.snapToBottom = h.viewport.AtBottom()
	}

	return h, tea.Batch(cmds...)
}

func (h *History) toggleSnap() {
	h.snapToBottom = !h.snapToBottom
	h.refreshViewport()
}

func (h *History) refreshViewport() {
	y := h.viewport.YOffset()
	h.viewport.SetContent(h.renderContent())
	if h.snapToBottom {
		h.viewport.GotoBottom()
	} else {
		h.viewport.SetYOffset(y)
	}
}

func (h History) renderContent() string {
	var b strings.Builder
	for _, c := range h.content {

		switch c := c.(type) {
		case ErrorContent:
			b.WriteString(errorStyle.Render("error"))
			b.WriteRune(' ')
			b.WriteString(textStyle.Render(c.Message))

		case SendingContent:
			b.WriteString(symbolStyle.Render(SymbolArrowRight))
			b.WriteRune(' ')
			b.WriteString(textStyle.Render(c.Route))

		case ResponseContent:
			b.WriteString(symbolStyle.Render(SymbolArrowLeft))
			b.WriteRune(' ')
			b.WriteString(textStyle.Render(c.Route))

		case EventContent:
			b.WriteString(eventStyle.Render(SymbolArrowLeft))
			b.WriteRune(' ')
			b.WriteString(textStyle.Render(c.Name))
		}

		b.WriteString("\n")
	}
	return strings.TrimSuffix(b.String(), "\n")
}

func (h History) View() tea.View     { return tea.NewView(h.viewport.View()) }
func (h History) ViewString() string { return h.viewport.View() }
