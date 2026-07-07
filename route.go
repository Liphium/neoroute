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
func Route[D any, RS any, PS interface {
	*RS
	msgp.Marshaler
}, RQ any, PQ interface {
	*RQ
	msgp.Unmarshaler
}](r Router[D], route string, handler func(c *ResCtx[D, RS, PS], req RQ) error) Router[D] {
	route = cleanRoute(r.getRoute() + string(RouteSeparator) + route)

	neos := r.getNeos()
	for _, neo := range neos {
		neo.routes[route] = func(c *Ctx[D]) error {

			// Parse request data into struct
			var data RQ
			unmarshaler := any(&data).(msgp.Unmarshaler)

			_, err := unmarshaler.UnmarshalMsg(c.reqData)
			if err != nil {
				return fmt.Errorf("failed to unmarshal struct: %v", err)
			}

			ctx := &ResCtx[D, RS, PS]{
				Ctx: c,
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

// RouteNoRequest does the same as Route but the handler does not receive a request struct, only the context.
// This can be useful if you only want to receive the request and don't want any data.
func RouteNoRequest[D any, RS any, PS interface {
	*RS
	msgp.Marshaler
}](r Router[D], route string, handler func(c *ResCtx[D, RS, PS]) error) Router[D] {
	route = cleanRoute(r.getRoute() + string(RouteSeparator) + route)

	neos := r.getNeos()
	for _, neo := range neos {
		neo.routes[route] = func(c *Ctx[D]) error {

			ctx := &ResCtx[D, RS, PS]{
				Ctx: c,
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

// RouteOk does the same as Route but the handler does not have a return type, it can only succeed or error.
// This can be useful if you don't have any return data, but the request can still have an error.
func RouteOk[D any, RQ any, PQ interface {
	*RQ
	msgp.Unmarshaler
}](r Router[D], route string, handler func(c *OkCtx[D], req RQ) error) Router[D] {
	route = cleanRoute(r.getRoute() + string(RouteSeparator) + route)

	neos := r.getNeos()
	for _, neo := range neos {
		neo.routes[route] = func(c *Ctx[D]) error {

			// Parse request data into struct
			var data RQ
			unmarshaler := any(&data).(msgp.Unmarshaler)

			_, err := unmarshaler.UnmarshalMsg(c.reqData)
			if err != nil {
				return fmt.Errorf("failed to unmarshal request data in RouteRequestOk for route %s: %v", route, err)
			}

			ctx := &OkCtx[D]{
				Ctx: c,
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

// RouteOkNoRequest does the same as RouteOk but the handler does not receive a request struct, only the context.
// This can be useful if you don't want to receive any data and the handler can only succeed or error.
func RouteOkNoRequest[D any](r Router[D], route string, handler func(c *OkCtx[D]) error) Router[D] {
	route = cleanRoute(r.getRoute() + string(RouteSeparator) + route)

	neos := r.getNeos()
	for _, neo := range neos {
		neo.routes[route] = func(c *Ctx[D]) error {
			ctx := &OkCtx[D]{
				Ctx: c,
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

// RouteNoResponse does the same as Route but the handler does not return anything.
// This can be useful if you only want to receive the data for example streaming over WebTransport.
func RouteNoResponse[D any, RQ any, PQ interface {
	*RQ
	msgp.Unmarshaler
}](r Router[D], route string, handler func(c *Ctx[D], req RQ)) Router[D] {
	route = cleanRoute(r.getRoute() + string(RouteSeparator) + route)

	neos := r.getNeos()
	for _, neo := range neos {
		neo.routes[route] = func(c *Ctx[D]) error {

			// Parse request data into struct
			var data RQ
			unmarshaler := any(&data).(msgp.Unmarshaler)

			_, err := unmarshaler.UnmarshalMsg(c.reqData)
			if err != nil {
				return fmt.Errorf("failed to unmarshal request data in RouteRequest for route %s: %v", route, err)
			}

			// Let the handler handle it
			handler(c, data)
			return &noResponse{}
		}
	}

	return &RouteRouter[D]{
		neos:  neos,
		route: route,
	}
}

// RoutePing is the same as Route but the handler does not receive a request struct, only the context.
// This can be useful if you only want to receive the request and don't want any data.
func RoutePing[D any](r Router[D], route string, handler func(c *Ctx[D])) Router[D] {
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
