package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/Liphium/neoroute/pkg/neodebug/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(2*time.Second, func(t time.Time) tea.Msg { return tickMsg(t) })
}

type model struct {
	h      tui.History
	width  int
	height int
}

func (m model) Init() tea.Cmd { return tea.Batch(m.h.Init(), tick()) }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		}
	case tickMsg:
		id := fmt.Sprintf("sec-%d", time.Now().UnixNano())
		nm, cmd := m.h.Update(tui.AddSectionMsg{Section: tui.Section{
			ID: id, Label: "Event " + id[:12], Loading: rand.Float32() < 0.3,
			Data: []tui.Content{{Line: fmt.Sprintf("random %d", rand.Intn(1000))}, {Line: time.Now().Format(time.RFC3339)}},
		}})
		m.h = nm.(tui.History)
		return m, tea.Batch(cmd, tick())
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// reserve 1 line for help bar
		h := msg.Height - 1
		if h < 1 {
			h = 1
		}
		nm, cmd := m.h.Update(tea.WindowSizeMsg{Width: msg.Width, Height: h})
		m.h = nm.(tui.History)
		return m, cmd
	}

	nm, cmd := m.h.Update(msg)
	m.h = nm.(tui.History)
	return m, cmd
}

func (m model) View() string {
	helpStyle := lipgloss.NewStyle().Faint(true)
	help := helpStyle.Render("q/esc quit • v snap to bottom • ↑/k ↓/j select • ←/h collapse →/l expand")
	// reserve 1 line for help bar
	h := m.h
	// give history height-1 if we know size
	if m.height > 1 {
		// push WindowSize down to history via View sizing; history already tracks via msg,
		// but ensure viewport height reflects reserved line
		_ = h
	}
	return lipgloss.JoinVertical(lipgloss.Left, h.View(), help)
}

func main() {
	h := tui.NewHistory(80, 20)
	p := tea.NewProgram(model{h: h}, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
