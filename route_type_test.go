package neoroute

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinylib/msgp/msgp"
)

type nonStructRoutePayload int

func (v nonStructRoutePayload) MarshalMsg(b []byte) ([]byte, error) {
	return msgp.AppendInt(b, int(v)), nil
}

func (v *nonStructRoutePayload) UnmarshalMsg(b []byte) ([]byte, error) {
	value, rest, err := msgp.ReadIntBytes(b)
	*v = nonStructRoutePayload(value)
	return rest, err
}

func TestRoutesRejectNonStructPayloadsAtRegistration(t *testing.T) {
	for _, tc := range []struct {
		name     string
		side     string
		register func(*Router[NoData])
	}{
		{"Route request", "request", func(r *Router[NoData]) {
			r.Route("test", func(*ResCtx[NoData, MsgStruct], nonStructRoutePayload) error { return nil })
		}},
		{"Route response", "response", func(r *Router[NoData]) {
			r.Route("test", func(*ResCtx[NoData, nonStructRoutePayload], MsgStruct) error { return nil })
		}},
		{"RouteNoRequest", "response", func(r *Router[NoData]) {
			r.RouteNoRequest("test", func(*ResCtx[NoData, nonStructRoutePayload]) error { return nil })
		}},
		{"RouteOk", "request", func(r *Router[NoData]) {
			r.RouteOk("test", func(*OkCtx[NoData], nonStructRoutePayload) error { return nil })
		}},
		{"RouteNoResponse", "request", func(r *Router[NoData]) {
			r.RouteNoResponse("test", func(*Ctx[NoData], nonStructRoutePayload) {})
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := NewRouter(Config[NoData]{})
			require.PanicsWithValue(t, `route "test" `+tc.side+` type must be a struct or pointer to a struct, got neoroute.nonStructRoutePayload`, func() {
				tc.register(r)
			})
		})
	}
}

func TestRouteTypeValidation(t *testing.T) {
	for _, check := range []func(){
		func() { validateRouteType[map[string]int]("test", "request") },
		func() { validateRouteType[[]int]("test", "request") },
		func() { validateRouteType[msgp.Marshaler]("test", "response") },
		func() { validateRouteType[*int]("test", "response") },
		func() { validateRouteType[**MsgStruct]("test", "response") },
	} {
		require.Panics(t, check)
	}
	require.Nil(t, validateRouteType[MsgStruct]("test", "response"))
	nilCheck := validateRouteType[*MsgStruct]("test", "response")
	require.NotNil(t, nilCheck)
	require.True(t, nilCheck(nil))
	require.False(t, nilCheck(&MsgStruct{}))
}
