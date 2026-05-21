package neoroute

import (
	"errors"
	"fmt"

	"github.com/tinylib/msgp/msgp"
)

type Router interface {
	Use()
	Group(route string) Router
	getRoute() string
	getNeo() *NeoRouter
}

type NeoRouter struct {
	routes     map[string]func(c *ctx) error
	middleware map[string]func(c *ctx) bool
	config     Config
}

// Route saves a handler for a given route.
// Be aware that only a-z, A-Z, 0-9, "-", "_", "~" can be used as characters for a route.
// To separate subroutes use "."
// Example routes: "", "route1", "route1.route2", "route1.route3"
// If characters are used that are not allowed, they will be striped, this can lead to unwanted behavior.
//
// Make sure the handler never returns nil, otherwise the router will panic.
func Route[RQ any, RS msgp.Marshaler, PQ interface {
	*RQ
	msgp.Unmarshaler
}](r Router, route string, handler func(c *Ctx[RS], req RQ) error) {
	route = cleanRoute(r.getRoute() + string(RouteSeparator) + route)

	neo := r.getNeo()
	neo.routes[route] = func(c *ctx) error {

		// Parse request data into struct
		var data RQ
		unmarshaler := any(&data).(msgp.Unmarshaler)

		_, err := unmarshaler.UnmarshalMsg(c.Data)
		if err != nil {
			return fmt.Errorf("failed to unmarshal struct: %v", err)
		}

		ctx := &Ctx[RS]{
			ctx: *c,
		}

		// Let the handler handle it
		return handler(ctx, data)
	}
}

// RouteWithout is the same as Route but the handler does not receive a request struct, only the context.
// This can be useful if you only want to receive the request and don't want any data.
func RouteWithout[RS msgp.Marshaler](r Router, route string, handler func(c *Ctx[RS]) error) {
	route = cleanRoute(r.getRoute() + string(RouteSeparator) + route)

	neo := r.getNeo()
	neo.routes[route] = func(c *ctx) error {

		ctx := &Ctx[RS]{
			ctx: *c,
		}

		// Let the handler handle it
		return handler(ctx)
	}
}

func (r *NeoRouter) Use() {

}

func (r *NeoRouter) Group(route string) Router {
	return &Group{
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

func (r *NeoRouter) handle(reqData []byte) []byte {

	c := &ctx{
		neo:   r,
		id:    -1,
		Data:  []byte{},
		Route: "",
	}

	var data Request
	_, err := data.UnmarshalMsg(reqData)
	if err != nil {
		logger.Info("failed to unmarshal request: ", err)
		return messageResponse(r, c.respondError(fmt.Errorf("Invalid request format.")))
	}

	c.id = data.Id
	c.Data = data.Data
	c.Route = data.Route

	// Handle request
	if handler, ok := r.routes[data.Route]; ok {
		if err := handler(c); err == nil {

			// Handlers never should return nil.
			panic("handler should always return something")
		} else if errors.Is(err, response{}) {

			// Return response from handler
			return messageResponse(c.neo, err.(response))
		} else {
			// Log error from handler and return generic error message to client
			logger.Info("an error occurred: ", err)
			return messageResponse(r, c.respondError(fmt.Errorf("Internal server error.")))
		}
	} else {
		return messageResponse(r, c.respondError(fmt.Errorf("Route does not exist.")))
	}

}
