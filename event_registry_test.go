package neoroute

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinylib/msgp/msgp"
)

type nonStructEventMap map[string]string

func (v nonStructEventMap) MarshalMsg(b []byte) ([]byte, error) {
	return msgp.AppendMapStrStr(b, v), nil
}

func TestEventRegistrationRejectsNonStructTypes(t *testing.T) {
	for _, tc := range []struct {
		name     string
		register func(*EventRegistry)
	}{
		{"scalar", func(r *EventRegistry) { r.Register[nonStructRoutePayload]("test") }},
		{"pointer to scalar", func(r *EventRegistry) { r.Register[*nonStructRoutePayload]("test") }},
		{"slice", func(r *EventRegistry) { r.Register[msgp.Raw]("test") }},
		{"map", func(r *EventRegistry) { r.Register[nonStructEventMap]("test") }},
		{"interface", func(r *EventRegistry) { r.Register[msgp.Marshaler]("test") }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			registry := NewEventRegistry()
			require.Panics(t, func() { tc.register(registry) })
			require.Empty(t, registry.GetEvents())
			require.Empty(t, registry.GetSchemas())
		})
	}
}

func TestEventStructAndPointerPayloads(t *testing.T) {
	registry := NewEventRegistry()
	valueBuilder := registry.Register[MsgStruct]("value")
	pointerBuilder := registry.Register[*MsgStruct]("pointer")
	want := MsgStruct{Val: "test"}
	for _, event := range []event{valueBuilder(want), pointerBuilder(&want)} {
		var got MsgStruct
		_, err := got.UnmarshalMsg(event.Data)
		require.NoError(t, err)
		require.Equal(t, want, got)
	}
	require.PanicsWithValue(t, `event "pointer" data must never be nil`, func() {
		pointerBuilder(nil)
	})
}
