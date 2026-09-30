package tui

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/Liphium/neoroute/pkg/neodebug/connector"
	"github.com/Liphium/neoroute/pkg/neodebug/model"
	"github.com/tinylib/msgp/msgp"
)

type sentRequest struct {
	name    string
	payload []byte
}

type inputHistory struct {
	width, height int
	connection    connector.Connection
	handledKey    bool
	requests      []sentRequest
	selected      int

	// Key bindings
	back key.Binding
	up   key.Binding
	down key.Binding
	send key.Binding
}

func newInputHistory() inputHistory {
	return inputHistory{
		back: key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
		up:   key.NewBinding(standardUpKey, key.WithHelp("↑", "up")),
		down: key.NewBinding(standardDownKey, key.WithHelp("↓", "down")),
		send: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "send")),
	}
}

// Children implements keyProvider.
func (m inputHistory) Children() []keyProvider {
	return []keyProvider{}
}

// FooterKeys implements keyProvider.
func (m inputHistory) FooterKeys() []key.Binding {
	return []key.Binding{m.send, m.up, m.down}
}

// FullKeyHelp implements keyProvider.
func (m inputHistory) FullKeyHelp() FullKeyHelp {
	return FullKeyHelp{
		Title: "Request history",
		Keys: [][]key.Binding{
			{m.send, m.up, m.down},
		},
	}
}

func (m *inputHistory) add(name string, payload []byte) {
	m.requests = append(m.requests, sentRequest{name: name, payload: append([]byte(nil), payload...)})
	m.selected = len(m.requests) - 1
}

func (m inputHistory) CanBeOpened() bool {
	return len(m.requests) > 0
}

func (m inputHistory) WantedHeight() int { return 1 + min(len(m.requests), 3) }

func (m inputHistory) View() string {
	title := titleStyle.Render("Request history")
	maxShown := max(0, m.height-1)
	start := max(0, min(m.selected-maxShown/2, len(m.requests)-maxShown))
	end := min(len(m.requests), start+maxShown)
	style := lipgloss.NewStyle().Width(m.width).Padding(0, 1)

	var b strings.Builder
	for i := start; i < end; i++ {

		// Convert the request bytes back to a value we can actually read
		r := m.requests[i]
		v, _, err := msgp.ReadIntfBytes(r.payload)
		detail := ""
		if err == nil {
			detail = fmt.Sprintf(" %v", v)
		}

		// Make sure the line does not exceed the width of the terminal
		line := r.name + detail
		limit := max(0, m.width-4)
		if len([]rune(line)) > limit {
			if limit > 3 {
				line = string([]rune(line)[:limit-3]) + "..."
			} else {
				line = string([]rune(line)[:limit])
			}
		}

		if i == m.selected {
			b.WriteString(style.Render(highlightStyle.Render(SymbolArrowRight) + " " + selectedStyle.Render(line)))
		} else {
			b.WriteString(style.Render(secondaryTextStyle.Render(SymbolArrowRight) + " " + unselectedStyle.Render(line)))
		}
		if i < end-1 {
			b.WriteByte('\n')
		}
	}
	return title + "\n" + b.String()
}

func (m *inputHistory) Update(msg tea.Msg) (exit bool, cmd tea.Cmd) {
	m.handledKey = false

	k, ok := msg.(tea.KeyPressMsg)
	if !ok {
		return false, nil
	}

	switch {
	case key.Matches(k, m.up):
		m.handledKey = true
		m.selected = max(0, m.selected-1)
		return false, nil
	case key.Matches(k, m.down):
		m.handledKey = true
		m.selected = min(len(m.requests)-1, m.selected+1)
		return false, nil
	case key.Matches(k, m.back):
		m.handledKey = true
		return true, nil
	case key.Matches(k, m.send):
		m.handledKey = true
		request := m.requests[m.selected]
		val, _, err := msgp.ReadIntfBytes(request.payload)
		if err != nil {
			return false, model.Plain(model.Error("Request payload is not decodable."))
		}
		return true, m.connection.Send(request.name, val)
	}

	return false, nil
}

// Wrapper for the send function to make sure it's always added to history
func (m *inputHistory) Send(route string, data any) tea.Cmd {
	payload, err := msgp.AppendIntf(nil, data)
	if err != nil {
		return model.Plain(model.Error("Input data is not encodable."))
	}
	m.add(route, payload)
	return m.connection.Send(route, data)
}
