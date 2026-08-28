package neoroute

import (
	"fmt"
	"maps"
	"testing"
)

func TestRouter_Init(t *testing.T) {
	middlewareFunc := func(c *Ctx[NoData]) error { return nil }
	tests := []struct {
		name        string
		routingFunc func(r *Router[NoData])
		want        map[string]int // route -> middleware amount
	}{
		{
			name: "no routes, no middleware",
			routingFunc: func(r *Router[NoData]) {
			},
			want: map[string]int{},
		},
		{
			name: "single route",
			routingFunc: func(r *Router[NoData]) {
				r.RoutePing("some-route", func(c *Ctx[NoData]) {})
			},
			want: map[string]int{"some-route": 0},
		},
		{
			name: "only use middleware, with no route",
			routingFunc: func(r *Router[NoData]) {
				r.Use("some-route", middlewareFunc)
			},
			want: map[string]int{},
		},
		{
			name: "route with middleware",
			routingFunc: func(r *Router[NoData]) {
				r.RoutePing("some-route", func(c *Ctx[NoData]) {}).Use("", middlewareFunc)
			},
			want: map[string]int{"some-route": 1},
		},
		{
			name: "route with multiple middlewares",
			routingFunc: func(r *Router[NoData]) {
				r.RoutePing("some-route", func(c *Ctx[NoData]) {}).Use("", middlewareFunc).Use("", middlewareFunc)
			},
			want: map[string]int{"some-route": 2},
		},
		{
			name: "router a group and two routes one has a middleware",
			routingFunc: func(r *Router[NoData]) {
				group := r.Group("some-group")
				group.RoutePing("some-route1", func(c *Ctx[NoData]) {})
				group.RoutePing("some-route2", func(c *Ctx[NoData]) {}).Use("", middlewareFunc)
			},
			want: map[string]int{"some-group/some-route1": 0, "some-group/some-route2": 1},
		},
		{
			name: "middleware on group, 2 routes and one has a middleware",
			routingFunc: func(r *Router[NoData]) {
				group := r.Group("some-group")
				group.Use("", middlewareFunc)
				group.RoutePing("some-route1", func(c *Ctx[NoData]) {})
				group.RoutePing("some-route2", func(c *Ctx[NoData]) {}).Use("", middlewareFunc)
			},
			want: map[string]int{"some-group/some-route1": 1, "some-group/some-route2": 2},
		},
		{
			name: "middleware on group applied at the end, 2 routes and one has a middleware",
			routingFunc: func(r *Router[NoData]) {
				group := r.Group("some-group")
				group.RoutePing("some-route1", func(c *Ctx[NoData]) {})
				group.RoutePing("some-route2", func(c *Ctx[NoData]) {}).Use("", middlewareFunc)
				group.Use("", middlewareFunc)
			},
			want: map[string]int{"some-group/some-route1": 1, "some-group/some-route2": 2},
		},
		{
			name: "middleware on subroute",
			routingFunc: func(r *Router[NoData]) {
				r.RoutePing("some-route1", func(c *Ctx[NoData]) {})
				r.Use("some-route1", middlewareFunc)
			},
			want: map[string]int{"some-route1": 1},
		},
		{
			name: "add route to router, then add it to base router",
			routingFunc: func(r *Router[NoData]) {
				router := NewRouter(Config[NoData]{})
				router.RoutePing("some-route", func(c *Ctx[NoData]) {})
				r.AddRouters(router)
			},
			want: map[string]int{"some-route": 0},
		},
		{
			name: "add sophisticated router to base router",
			routingFunc: func(r *Router[NoData]) {
				router := NewRouter(Config[NoData]{})
				router.RoutePing("some-route", func(c *Ctx[NoData]) {})
				group := router.Group("some-group")
				group.RoutePing("some-route1", func(c *Ctx[NoData]) {})
				group.Use("", middlewareFunc)
				r.AddRouters(router)
			},
			want: map[string]int{"some-route": 0, "some-group/some-route1": 1},
		},
		{
			name: "stack multiple routers",
			routingFunc: func(r *Router[NoData]) {
				router := NewRouter(Config[NoData]{})
				router2 := NewRouter(Config[NoData]{})
				router3 := NewRouter(Config[NoData]{})
				group := router2.Group("some-group")
				router3.RoutePing("some-route1", func(c *Ctx[NoData]) {})
				router3.Use("", middlewareFunc)
				group.AddRouters(router3)
				r.AddRouters(router, router2)
			},
			want: map[string]int{"some-group/some-route1": 1},
		},
		{
			name: "middleware on a group and a route with same name",
			routingFunc: func(r *Router[NoData]) {
				r.Group("group1").Use("", middlewareFunc)
				r.RoutePing("group1", func(c *Ctx[NoData]) {})
			},
			want: map[string]int{"group1": 0},
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
