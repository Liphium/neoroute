package neoroute

import (
	"testing"
)

func TestUse(t *testing.T) {
	type args struct {
		router     Router[any]
		route      string
		middleware func(c *Ctx[any]) bool
	}
	tests := []struct {
		name      string
		args      args
		wantRoute string
	}{
		{
			name: "set middleware for route with one neo",
			args: args{
				router:     NewNeoRouter[any](Config{}),
				route:      "test...route",
				middleware: func(c *Ctx[any]) bool { return true },
			},
			wantRoute: "test.route",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.router.Use(tt.args.route, tt.args.middleware)
			for i, neo := range tt.args.router.getNeos() {
				if _, ok := neo.middleware[tt.wantRoute]; !ok {
					t.Errorf("Use() middleware not set for route %v for neo %d", tt.args.route, i)
				}
			}
		})
	}
}
