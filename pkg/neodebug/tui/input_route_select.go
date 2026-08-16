package tui

import tea "charm.land/bubbletea/v2"

type inputRouteSelect struct {
	width  int
	height int
}

func (m inputRouteSelect) SetWidth(w int) {
	m.width = w
}

func (m inputRouteSelect) SetHeight(h int) {
	m.height = h
}

func (m inputRouteSelect) WantedHeight() int {
	return 5 /* input field + divider + 3 recommendations */
}

func (m inputRouteSelect) Update(msg tea.Msg) (inputRouteSelect, tea.Cmd) {
	return m, nil
}

func (m inputRouteSelect) View() string {
	return ""
}
