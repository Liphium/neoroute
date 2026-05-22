package neoroute

import (
	"errors"
	"fmt"

	"github.com/tinylib/msgp/msgp"
)

type Router[D any] interface {
	Use()
	Group(route string) Router[D]
	getRoute() string
	getNeo() *NeoRouter[D]
}

type NeoRouter[D any] struct {
	routes     map[string]func(c *ctx[D]) error
	middleware map[string]func(c *ctx[D]) bool
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
}, D any](r Router[D], route string, handler func(c *Ctx[RS, D], req RQ) error) {
	route = cleanRoute(r.getRoute() + string(RouteSeparator) + route)

	neo := r.getNeo()
	neo.routes[route] = func(c *ctx[D]) error {

		// Parse request data into struct
		var data RQ
		unmarshaler := any(&data).(msgp.Unmarshaler)

		_, err := unmarshaler.UnmarshalMsg(c.data)
		if err != nil {
			return fmt.Errorf("failed to unmarshal struct: %v", err)
		}

		ctx := &Ctx[RS, D]{
			ctx: *c,
		}

		// Let the handler handle it
		return handler(ctx, data)
	}
}

// RouteWithout is the same as Route but the handler does not receive a request struct, only the context.
// This can be useful if you only want to receive the request and don't want any data.
func RouteWithout[RS msgp.Marshaler, D any](r Router[D], route string, handler func(c *Ctx[RS, D]) error) {
	route = cleanRoute(r.getRoute() + string(RouteSeparator) + route)

	neo := r.getNeo()
	neo.routes[route] = func(c *ctx[D]) error {

		ctx := &Ctx[RS, D]{
			ctx: *c,
		}

		// Let the handler handle it
		return handler(ctx)
	}
}

func (r *NeoRouter[D]) Use() {

}

func (r *NeoRouter[D]) Group(route string) Router[D] {
	return &Group[D]{
		neo:    r,
		prefix: route,
		parent: nil,
	}
}

func (r *NeoRouter[D]) getRoute() string {
	return ""
}

func (r *NeoRouter[D]) getNeo() *NeoRouter[D] {
	return r
}

func (r *NeoRouter[D]) handle(reqData []byte, session *Session[D]) []byte {

	c := &ctx[D]{
		neo:     r,
		id:      -1,
		data:    []byte{},
		route:   "",
		session: session,
	}

	var data Request
	_, err := data.UnmarshalMsg(reqData)
	if err != nil {
		logger.Info("failed to unmarshal request", "err", err)
		return messageResponse(r, c.respondError(fmt.Errorf("Invalid request format.")))
	}

	c.id = data.Id
	c.data = data.Data
	c.route = data.Route

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
			logger.Info("an error occurred", "err", err)
			return messageResponse(r, c.respondError(fmt.Errorf("Internal server error.")))
		}
	} else {
		return messageResponse(r, c.respondError(fmt.Errorf("Route does not exist.")))
	}

}
