package tui

import (
	"strings"
	"testing"

	"github.com/Liphium/neoroute/pkg/neodebug/v2/config"
	"github.com/charmbracelet/x/ansi"
	"github.com/stretchr/testify/require"
)

func TestClearSearchPreservesSelectedItem(t *testing.T) {
	first := &config.Preset{Name: "Alice", Route: "send"}
	second := &config.Preset{Name: "Bob", Route: "send"}
	m := newRouteSelect(nil)
	m.routes = []routeItem{{Name: "send"}, {Name: "Alice", Preset: first}, {Name: "Bob", Preset: second}}
	m.results = []routeItem{m.routes[2], m.routes[1]}
	m.selectedRoute = 0
	m.clearSearchInput()
	require.Equal(t, 2, m.selectedRoute)
	m.clearSearchInput()
	require.Equal(t, 2, m.selectedRoute)
	m.results = nil
	m.selectedRoute = -1
	m.clearSearchInput()
	require.Equal(t, 0, m.selectedRoute)
}

func TestRouteWarningTruncation(t *testing.T) {
	for _, warning := range []string{"short", strings.Repeat("long warning ", 10), strings.Repeat("警告", 30)} {
		m := newRouteSelect(nil)
		m.width, m.height = 40, 2
		m.results = []routeItem{{Name: "send", Warning: warning}}
		_, view := m.View()
		row := strings.Split(ansi.Strip(view), "\n")[1]
		require.LessOrEqual(t, ansi.StringWidth(row), m.width)
		require.True(t, strings.HasSuffix(strings.TrimSpace(row), ")"))
		if warning != "short" {
			require.Contains(t, row, "...)")
		} else {
			require.Contains(t, row, "(short)")
		}
	}
}
