package neoroute

import (
	"errors"
	"io"

	"github.com/tinylib/msgp/msgp"
)

type MiddlewareFunc[D any] = func(c *Ctx[D]) bool

type Router[D any] struct {
	prefix       string
	actualRoutes map[string]exportedRoute[D]
	config       Config[D]

	// temporary during tree building phase
	middlewares map[string][]MiddlewareFunc[D]
	route       RouteData[D]
	children    map[string][]*Router[D]
}

func (r *Router[D]) Config() Config[D] {
	return r.config
}

// GetRoutes is used for schema generation.
func (r Router[D]) GetRoutes() map[string]RouteData[D] {
	routes := map[string]RouteData[D]{}
	for route, exported := range r.actualRoutes {
		routes[route] = exported.RouteData
	}
	return routes
}

type exportedRoute[D any] struct {
	middlewares []MiddlewareFunc[D]
	RouteData[D]
}

// Group creates a new route group for the given prefix.
func (n *Router[D]) Group(subroute string) *Router[D] {
	subroute = cleanRoute(subroute)

	r := &Router[D]{
		prefix:       subroute,
		actualRoutes: make(map[string]exportedRoute[D]),
		children:     make(map[string][]*Router[D]),
		middlewares:  make(map[string][]MiddlewareFunc[D]),
		route:        RouteData[D]{},
	}
	n.children[subroute] = append(n.children[subroute], r)
	return r
}

func (r *Router[D]) AddRouters(router *Router[D], routers ...*Router[D]) *Router[D] {
	r.children[""] = append(r.children[""], append(routers, router)...)
	return r
}

func (r *Router[D]) Use(subroute string, middleware func(c *Ctx[D]) bool) {
	subroute = cleanRoute(subroute)
	r.middlewares[subroute] = append(r.middlewares[subroute], middleware)
}

// Handle is called by transporters to handle incoming requests.
//
// ONLY USE THIS IN A TRANSPORTER IMPLEMENTATION, THIS IS NOT MEANT TO BE USED BY USERS OF THE LIBRARY.
func (r *Router[D]) Handle(reqReader io.Reader, session *Session[D]) ([]byte, []func()) {
	c := &Ctx[D]{
		id:      -1,
		reqData: []byte{},
		route:   "",
		session: session,
	}

	var data request
	err := data.DecodeMsg(msgp.NewReader(reqReader))
	if err != nil {
		return messageResponse(c.respondError(r.config.RunErrorHandler(err, c))), nil
	}

	route := cleanRoute(data.Route)

	c.id = data.Id
	c.reqData = data.Data
	c.route = route

	// Check if handler for route exists
	routeData, exists := r.actualRoutes[route]
	if !exists {
		return messageResponse(c.respondError(r.config.RunErrorHandler(ErrRouteDoesntExist, c))), nil
	}

	// Run middlewares
	for _, middleware := range routeData.middlewares {
		if !middleware(c) {
			return messageResponse(c.respondError(r.config.RunErrorHandler(ErrMiddlewareDenied, c))), nil
		}
	}

	// Handle request
	err = routeData.handler(c)
	if err == nil {

		// Handlers never should return nil.
		panic("handler should always return something")
	}

	if respData, ok := errors.AsType[*responseData](err); ok {

		// Return response from handler
		resp := response{
			Id:      c.id,
			HasData: respData.HasData,
			IsError: respData.IsError,
			Data:    respData.Data,
		}
		return messageResponse(resp), c.runAfter
	}

	if _, ok := errors.AsType[*noResponse](err); ok {

		// Return no response
		return nil, c.runAfter
	}

	// Let user handle the error and decide what error message to send back to the client
	return messageResponse(c.respondError(r.config.RunErrorHandler(err, c))), c.runAfter
}
