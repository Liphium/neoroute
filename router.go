package neoroute

import (
	"fmt"
)

type Router interface {
	Use()
	Mount(route string) Router
	getRoute() string
	getNeo() *NeoRouter
}

type NeoRouter struct {
	routes map[string]func(c *Ctx) error
}

// Route saves a handler for a given route.
// Be aware that only a-z, A-Z, 0-9, "-", "_", "~" can be used as characters for a route.
// To separate subroutes use "."
// Example routes: "", "route1", "route1.route2", "route1.route3"
// If characters are used that are not allowed, they will be striped, this can lead to unwanted behavior.
func Route[A any](r Router, route string, handler func(c *Ctx, req A) error) {
	route = r.getRoute() + route

	neo := r.getNeo()
	neo.routes[route] = func(c *Ctx) error {

		// Parse the request
		var req Request
		_, err := req.UnmarshalMsg(c.data)
		if err != nil {
			return fmt.Errorf("failed to unmarshal struct: %v", err)
		}

		var data A
		_, err = req.UnmarshalMsg(req.Data)
		if err != nil {
			return fmt.Errorf("failed to unmarshal struct: %v", err)
		}

		// Let the handler handle it
		return handler(c, data)
	}
}

func (r *NeoRouter) Use() {

}

func (r *NeoRouter) Mount(route string) Router {
	return &Mount{
		neo:    r,
		prefix: route,
		parent: nil,
	}
}

func (r *NeoRouter) getRoute() string {
	return ""
}

func (r *NeoRouter) getNeo() *NeoRouter {
	return r
}

func (r *NeoRouter) handle(route string) error {
	return nil
}
