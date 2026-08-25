package neoroute

import (
	"fmt"
	"maps"
	"testing"
)

const (
	routingTypeGroup     = "group"
	routingTypeRoute     = "route"
	routingTypeAddRouter = "add-router"
	routingTypeUse       = "use"
)

func TestRouter_Init(t *testing.T) {
	middlewareFunc := func(c *Ctx[NoData]) bool { return false }
	tests := []struct {
		name        string
		routingFunc func(r *Router[NoData])
		want        map[string]int // route -> middleware amount
	}{
		{
			name: "only use middleware, with no route",
			routingFunc: func(r *Router[NoData]) {
				r.Use("some-route", middlewareFunc)
			},
			want: map[string]int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRouter(Config[NoData]{})
			tt.routingFunc(r)
			r.Init()
			got := map[string]int{}
			for route, exportedRoute := range r.actualRoutes {
				fmt.Printf("route: '%v'; middlewares: '%v'\n", route, len(exportedRoute.middlewares))
				got[route] = len(exportedRoute.middlewares)
			}

			if !maps.Equal(tt.want, got) {
				t.Errorf("RunErrorHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}
