package tui

import (
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/Liphium/neoroute/pkg/neodebug/v2/config"
	"github.com/Liphium/neoroute/pkg/neodebug/v2/connector"
	"github.com/Liphium/neoroute/v2/neoschema"
	"github.com/stretchr/testify/require"
)

func TestPresetSelectionActions(t *testing.T) {
	for _, warning := range []string{"", "invalid body"} {
		t.Run(warning, func(t *testing.T) {
			m := newRouteSelect(nil)
			preset := &config.Preset{Route: "unknown", Data: "body"}
			m.routes = []routeItem{{Name: "preset", Preset: preset, Warning: warning}}
			m.results = m.routes
			m.input.Focus()
			edited, cmd := m.Update(tea.KeyPressMsg{Code: 'e', Mod: tea.ModAlt})
			if warning == "" {
				require.NotNil(t, cmd)
				require.Equal(t, RouteSelectedMsg{Route: "unknown", Preset: preset, Edit: true}, cmd())
			} else {
				require.Nil(t, cmd)
			}
			require.True(t, edited.handledKey)
			_, cmd = m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
			require.NotNil(t, cmd)
			require.Equal(t, RouteSelectedMsg{Route: "unknown", Preset: preset}, cmd())
		})
	}
}

func TestPresetDraftIsIndependent(t *testing.T) {
	text := &neoschema.BasicType{ActualType: neoschema.TypeString}
	root := &StructNode{children: []StructField{
		{Name: "sender", Node: createNode(text, nil)},
		{Name: "tags", Node: &SliceNode{element: text}},
		{Name: "metadata", Node: &MapNode{keyType: text, valueType: text, selected: -1}},
		{Name: "optional", Node: &NullableNode{null: true, other: createNode(text, nil)}},
		{Name: "date", Node: newDateNode()},
		{Name: "file", Node: newFileNode()},
	}}
	root.Init()
	root.OnUp(func() {})
	root.OnDown(func() {})
	date := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	preset := map[string]any{"sender": "Alice", "tags": []any{"hello"}, "metadata": map[string]any{"key": "value"}, "optional": "present", "date": date, "file": []byte{1, 2, 3}}
	creator := inputRequestCreator{rootNode: root}
	require.NoError(t, creator.loadPreset(preset))
	request := root.Request().(map[string]any)
	require.Equal(t, "Alice", request["sender"])
	require.Equal(t, []any{"hello"}, request["tags"])
	require.Equal(t, map[string]any{"key": "value"}, request["metadata"])
	require.Equal(t, "present", request["optional"])
	require.True(t, date.Equal(request["date"].(time.Time)))
	require.Equal(t, []byte{1, 2, 3}, request["file"])
	root.children[0].Node.(*ValueNode[string]).input.SetValue("Bob")
	request["file"].([]byte)[0] = 9
	require.Equal(t, "Alice", preset["sender"])
	require.Equal(t, byte(1), preset["file"].([]byte)[0])
	require.Equal(t, "Bob", root.Request().(map[string]any)["sender"])
}

type presetTestConnection struct {
	connector.Connection
	route string
	body  any
}

func (c *presetTestConnection) Send(route string, body any) tea.Cmd {
	c.route, c.body = route, body
	return nil
}

func TestPresetSendsWithoutSchemaValidation(t *testing.T) {
	connection := &presetTestConnection{}
	input := newInput(neoschema.TransporterSchema{})
	input.state = stateRouteSelect
	input.history.connection = connection
	preset := &config.Preset{Route: "unknown", Data: "invalid body"}
	updated, _ := input.Update(RouteSelectedMsg{Route: preset.Route, Preset: preset})
	require.Equal(t, "unknown", connection.route)
	require.Equal(t, preset.Data, connection.body)
	require.Equal(t, stateRouteSelect, updated.state)
	require.Len(t, updated.history.requests, 1)
}
