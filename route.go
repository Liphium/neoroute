package neoroute

import (
	"fmt"

	"github.com/tinylib/msgp/msgp"
)

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
}, D any](r Router[D], route string, handler func(c *ResCtx[RS, D], req RQ) error) Router[D] {
	route = cleanRoute(r.getRoute() + string(RouteSeparator) + route)

	neos := r.getNeos()
	for _, neo := range neos {
		neo.routes[route] = func(c *Ctx[D]) error {

			// Parse request data into struct
			var data RQ
			unmarshaler := any(&data).(msgp.Unmarshaler)

			_, err := unmarshaler.UnmarshalMsg(c.data)
			if err != nil {
				return fmt.Errorf("failed to unmarshal struct: %v", err)
			}

			ctx := &ResCtx[RS, D]{
				Ctx: *c,
			}

			// Let the handler handle it
			return handler(ctx, data)
		}
	}

	return &RouteRouter[D]{
		neos:  neos,
		route: route,
	}
}

// RouteResponse does the same as Route but the handler does not receive a request struct, only the context.
// This can be useful if you only want to receive the request and don't want any data.
func RouteResponse[RS msgp.Marshaler, D any](r Router[D], route string, handler func(c *ResCtx[RS, D]) error) Router[D] {
	route = cleanRoute(r.getRoute() + string(RouteSeparator) + route)

	neos := r.getNeos()
	for _, neo := range neos {
		neo.routes[route] = func(c *Ctx[D]) error {

			ctx := &ResCtx[RS, D]{
				Ctx: *c,
			}

			// Let the handler handle it
			return handler(ctx)
		}
	}

	return &RouteRouter[D]{
		neos:  neos,
		route: route,
	}
}

// RouteRequest does the same as Route but the handler does not return anything.
// This can be useful if you only want to receive the data for example streaming over WebTransport.
func RouteRequest[RQ any, RS msgp.Marshaler, PQ interface {
	*RQ
	msgp.Unmarshaler
}, D any](r Router[D], route string, handler func(c *ResCtx[RS, D], req RQ)) Router[D] {
	route = cleanRoute(r.getRoute() + string(RouteSeparator) + route)

	neos := r.getNeos()
	for _, neo := range neos {
		neo.routes[route] = func(c *Ctx[D]) error {

			// Parse request data into struct
			var data RQ
			unmarshaler := any(&data).(msgp.Unmarshaler)

			_, err := unmarshaler.UnmarshalMsg(c.data)
			if err != nil {
				logger.Info("failed to unmarshal request data in RouteRequest", "route", route, "err", err)
				return &noResponse{}
			}

			ctx := &ResCtx[RS, D]{
				Ctx: *c,
			}

			// Let the handler handle it
			handler(ctx, data)
			return &noResponse{}
		}
	}

	return &RouteRouter[D]{
		neos:  neos,
		route: route,
	}
}

// RouteNoop is the same as Route but the handler does not receive a request struct, only the context.
// This can be useful if you only want to receive the request and don't want any data.
func RouteNoop[D any](r Router[D], route string, handler func(c *Ctx[D])) Router[D] {
	route = cleanRoute(r.getRoute() + string(RouteSeparator) + route)

	neos := r.getNeos()
	for _, neo := range neos {
		neo.routes[route] = func(c *Ctx[D]) error {

			// Let the handler handle it
			handler(c)
			return &noResponse{}
		}
	}

	return &RouteRouter[D]{
		neos:  neos,
		route: route,
	}
}
