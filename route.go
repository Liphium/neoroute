package neoroute

import (
	"fmt"
	"reflect"

	"github.com/tinylib/msgp/msgp"
)

// Internal struct used to represent routes for schema generation.
type RouteData[D any] struct {
	handler func(c *Ctx[D]) error

	// Should return the request type for schema generation, return false if no request.
	RequestType func() (bool, reflect.Type)

	// Should return the response type for schema generation, return false for custom if no response.
	ResponseType func() (error bool, custom bool, response reflect.Type)

	// hasError indicates if the handler returns an error.
	hasError bool
}

// Route saves a handler for a given route.
//
// Be aware that only a-z, A-Z, 0-9, "-", "_", "~", "." can be used as characters for a route. To separate subroutes use "/". Example for valid routes are: "route1", "route1/route2", "route1/route3".
//
// Any non-allowed characters will simply be stripped.
//
// Make sure the handler never returns nil, otherwise the router will panic.
func (r *Router[D]) Route[RS msgp.Marshaler, RQ any, PQ msgp.UnmarshalPtr[RQ]](route string, handler func(c *ResCtx[D, RS], req RQ) error) *Router[D] {
	panicIfPointer[RS](route)

	return addRoute(r, route, RouteData[D]{
		handler: func(c *Ctx[D]) error {

			// Parse request data into struct
			var data RQ
			unmarshaler := any(&data).(msgp.Unmarshaler)

			_, err := unmarshaler.UnmarshalMsg(c.reqData)
			if err != nil {
				return fmt.Errorf("failed to unmarshal struct: %v", err)
			}

			ctx := &ResCtx[D, RS]{
				Ctx: c,
			}

			// Let the handler handle it
			return handler(ctx, data)
		},
		RequestType: func() (bool, reflect.Type) {
			return true, reflect.TypeFor[RQ]()
		},
		ResponseType: func() (bool, bool, reflect.Type) {
			return true, true, reflect.TypeFor[RS]()
		},
		hasError: true,
	})
}

// RouteNoRequest is the same as Route but the handler does not receive a request struct, only the context.
//
// This can be useful if you only want to receive the request and don't want any data.
func (r *Router[D]) RouteNoRequest[RS msgp.Marshaler](route string, handler func(c *ResCtx[D, RS]) error) *Router[D] {
	panicIfPointer[RS](route)

	return addRoute(r, route, RouteData[D]{
		handler: func(c *Ctx[D]) error {

			ctx := &ResCtx[D, RS]{
				Ctx: c,
			}

			// Let the handler handle it
			return handler(ctx)
		},
		RequestType: func() (bool, reflect.Type) {
			return false, nil
		},
		ResponseType: func() (bool, bool, reflect.Type) {
			return true, true, reflect.TypeFor[RS]()
		},
		hasError: true,
	})
}

// RouteOk is the same as Route but the handler does not have a return type, it can only succeed or error.
//
// This can be useful if you don't have any return data, but the request can still have an error.
func (r *Router[D]) RouteOk[RQ any, PQ msgp.UnmarshalPtr[RQ]](route string, handler func(c *OkCtx[D], req RQ) error) *Router[D] {
	return addRoute(r, route, RouteData[D]{
		handler: func(c *Ctx[D]) error {

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
		},
		RequestType: func() (bool, reflect.Type) {
			return true, reflect.TypeFor[RQ]()
		},
		ResponseType: func() (bool, bool, reflect.Type) {
			return true, false, nil
		},
		hasError: true,
	})
}

// RouteOkNoRequest is the same as RouteOk but the handler does not receive a request struct, only the context.
//
// This can be useful if you don't want to receive any data and the handler can only succeed or error.
func (r *Router[D]) RouteOkNoRequest(route string, handler func(c *OkCtx[D]) error) *Router[D] {
	return addRoute(r, route, RouteData[D]{
		handler: func(c *Ctx[D]) error {
			ctx := &OkCtx[D]{
				Ctx: c,
			}

			// Let the handler handle it
			return handler(ctx)
		},
		RequestType: func() (bool, reflect.Type) {
			return false, nil
		},
		ResponseType: func() (bool, bool, reflect.Type) {
			return true, false, nil
		},
		hasError: true,
	})
}

// RouteNoResponse is the same as Route but the handler does not return anything.
//
// This can be useful if you only want to receive the data for example streaming over WebTransport.
func (r *Router[D]) RouteNoResponse[RQ any, PQ msgp.UnmarshalPtr[RQ]](route string, handler func(c *Ctx[D], req RQ)) *Router[D] {
	return addRoute(r, route, RouteData[D]{
		handler: func(c *Ctx[D]) error {

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
		},
		RequestType: func() (bool, reflect.Type) {
			return true, reflect.TypeFor[RQ]()
		},
		ResponseType: func() (bool, bool, reflect.Type) {
			return false, false, nil
		},
		hasError: false,
	})
}

// RoutePing is the same as Route but the handler does not receive a request struct, only the context.
//
// This can be useful if you only want to receive the request and don't want any data.
func (r *Router[D]) RoutePing(route string, handler func(c *Ctx[D])) *Router[D] {
	return addRoute(r, route, RouteData[D]{
		handler: func(c *Ctx[D]) error {

			// Let the handler handle it
			handler(c)
			return &noResponse{}
		},
		RequestType: func() (bool, reflect.Type) {
			return false, nil
		},
		ResponseType: func() (bool, bool, reflect.Type) {
			return false, false, nil
		},
		hasError: false,
	})
}

func panicIfPointer[T any](route string) {
	t := reflect.TypeFor[T]()
	if t.Kind() == reflect.Pointer {
		panic(fmt.Sprintf("%s: pointers are not allowed in routes due to nil not being encodable, use a struct instead of a pointer", route))
	}
}

func addRoute[D any](r *Router[D], route string, routeData RouteData[D]) *Router[D] {
	route = cleanRoute(route)
	if route == "" {
		panic("route cannot be empty")
	}

	router := r.Group(route)
	router.hasRoute = true
	router.route = routeData
	return router
}
