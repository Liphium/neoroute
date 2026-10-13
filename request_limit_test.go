package neoroute

import (
	"bytes"
	"errors"
	"testing"
)

func encodeTestRequest(t *testing.T, route string, data []byte) *bytes.Reader {
	t.Helper()
	b, err := (&request{Id: 1, Route: route, Data: data}).MarshalMsg(nil)
	if err != nil {
		t.Fatal(err)
	}
	return bytes.NewReader(b)
}

func TestRequestSizeLimit(t *testing.T) {
	data := make([]byte, 100)
	tests := []struct {
		name   string
		config Config[NoData]
		opts   []func(*RouteOptions)
		tooBig bool
	}{
		{name: "route limit exceeded", opts: []func(*RouteOptions){WithMaxRequestSize(10)}, tooBig: true},
		{name: "route limit fits", opts: []func(*RouteOptions){WithMaxRequestSize(1000)}},
		{name: "route unlimited overrides router default", config: Config[NoData]{RouteOptions: RouteOptions{MaxRequestSize: 10}}, opts: []func(*RouteOptions){WithMaxRequestSize(0)}},
		{name: "router default inherited", config: Config[NoData]{RouteOptions: RouteOptions{MaxRequestSize: 10}}, tooBig: true},
		{name: "unlimited by default", config: Config[NoData]{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRouter(tt.config)
			r.RoutePing("ping", func(c *Ctx[NoData]) {}, tt.opts...)
			r.Init()

			_, _, err := r.HandleRequest(encodeTestRequest(t, "ping", data), NewSession("s", NoData{}, SessionTransporterCallbacks[NoData]{}))
			if tooBig := errors.Is(err, ErrRequestTooLarge); tooBig != tt.tooBig {
				t.Fatalf("got too large %v, want %v (err %v)", tooBig, tt.tooBig, err)
			}
		})
	}
}

func TestMaxRequestSize(t *testing.T) {
	r := NewRouter(Config[NoData]{RouteOptions: RouteOptions{MaxRequestSize: 10}})
	r.RoutePing("small", func(c *Ctx[NoData]) {}, WithMaxRequestSize(20))
	r.RoutePing("default", func(c *Ctx[NoData]) {})
	r.Init()
	if got := r.MaxRequestSize(); got != 20 {
		t.Fatalf("got %d, want 20", got)
	}

	r = NewRouter(Config[NoData]{RouteOptions: RouteOptions{MaxRequestSize: 10}})
	r.RoutePing("unlimited", func(c *Ctx[NoData]) {}, WithMaxRequestSize(0))
	r.Init()
	if got := r.MaxRequestSize(); got != 0 {
		t.Fatalf("got %d, want 0", got)
	}
}
