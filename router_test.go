package neoroute

import (
	"testing"
)

func TestNeoRouter_getNeos(t *testing.T) {
	tests := []struct {
		name  string
		setup func() (Router[any], []*NeoRouter[any])
	}{
		{
			name: "a single neo router",
			setup: func() (Router[any], []*NeoRouter[any]) {
				router := NewNeoRouter[any](Config[any]{})
				return router, []*NeoRouter[any]{router}
			},
		},
		{
			name: "multiple neo router with AddRouters",
			setup: func() (Router[any], []*NeoRouter[any]) {
				router := NewNeoRouter[any](Config[any]{})
				router2 := NewNeoRouter[any](Config[any]{})
				router3 := NewNeoRouter[any](Config[any]{})
				router.AddRouters(router2, router3)
				return router, []*NeoRouter[any]{router, router2, router3}
			},
		},
		{
			name: "add routers, that already have routers",
			setup: func() (Router[any], []*NeoRouter[any]) {
				router := NewNeoRouter[any](Config[any]{})
				router2 := NewNeoRouter[any](Config[any]{})
				router3 := NewNeoRouter[any](Config[any]{})
				router.AddRouters(router2, router3)

				router4 := NewNeoRouter[any](Config[any]{})
				router5 := NewNeoRouter[any](Config[any]{})
				router6 := NewNeoRouter[any](Config[any]{})
				router4.AddRouters(router5, router6)

				router7 := NewNeoRouter[any](Config[any]{})
				router8 := NewNeoRouter[any](Config[any]{})
				router9 := NewNeoRouter[any](Config[any]{})
				router7.AddRouters(router8, router9)

				router.AddRouters(router4, router7)

				return router, []*NeoRouter[any]{router, router2, router3, router4, router5, router6, router7, router8, router9}
			},
		},
		{
			name: "add circular routers",
			setup: func() (Router[any], []*NeoRouter[any]) {
				router := NewNeoRouter[any](Config[any]{})
				router2 := NewNeoRouter[any](Config[any]{})
				router.AddRouters(router2)
				router2.AddRouters(router)

				return router, []*NeoRouter[any]{router, router2}
			},
		},
		{
			name: "router adding it self",
			setup: func() (Router[any], []*NeoRouter[any]) {
				router := NewNeoRouter[any](Config[any]{})
				router.AddRouters(router)

				return router, []*NeoRouter[any]{router}
			},
		},
		{
			name: "circular routers with multiple levels",
			setup: func() (Router[any], []*NeoRouter[any]) {
				router := NewNeoRouter[any](Config[any]{})
				router2 := NewNeoRouter[any](Config[any]{})
				router3 := NewNeoRouter[any](Config[any]{})
				router.AddRouters(router2, router3)

				router4 := NewNeoRouter[any](Config[any]{})
				router5 := NewNeoRouter[any](Config[any]{})
				router6 := NewNeoRouter[any](Config[any]{})
				router4.AddRouters(router5, router6)

				router7 := NewNeoRouter[any](Config[any]{})
				router8 := NewNeoRouter[any](Config[any]{})
				router9 := NewNeoRouter[any](Config[any]{})
				router7.AddRouters(router8, router9)

				router.AddRouters(router4, router7)
				router4.AddRouters(router4, router)
				router7.AddRouters(router7, router)

				return router, []*NeoRouter[any]{router, router2, router3, router4, router5, router6, router7, router8, router9}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, want := tt.setup()
			got := router.getNeos()
			if !sameElements(got, want) {
				t.Errorf("getNeos() = %v, want %v", got, want)
			}
		})
	}
}
