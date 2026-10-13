package neoroute

import (
	"bytes"
	"errors"
	"io"
	"sync"

	"github.com/tinylib/msgp/msgp"
)

type MiddlewareFunc[D any] = func(c *Ctx[D]) error

// Router is the main router struct, it holds the routes and middleware functions.
type Router[D any] struct {
	initOnce     sync.Once
	actualRoutes map[string]exportedRoute[D]
	config       Config[D]

	// temporary during tree building phase
	middlewares map[string][]MiddlewareFunc[D]
	hasRoute    bool
	route       RouteData[D]
	children    map[string][]*Router[D]
}

// NewRouter returns a new Router instance with the given config.
// Use this to create a new router instance to pass to the transporters.
func NewRouter[D any](config Config[D]) *Router[D] {

	// Set config default values to make sure we have the proper config in here
	config.setDefaultOptions()

	return &Router[D]{
		config:       config,
		actualRoutes: make(map[string]exportedRoute[D]),
		children:     make(map[string][]*Router[D]),
		middlewares:  make(map[string][]MiddlewareFunc[D]),
	}
}

// Config returns the router's config.
func (r *Router[D]) Config() Config[D] {
	return r.config
}

// GetRoutes is used for schema generation.
func (r *Router[D]) GetRoutes() map[string]RouteData[D] {
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
		actualRoutes: make(map[string]exportedRoute[D]),
		children:     make(map[string][]*Router[D]),
		middlewares:  make(map[string][]MiddlewareFunc[D]),
		route:        RouteData[D]{},
	}
	n.children[subroute] = append(n.children[subroute], r)
	return r
}

// AddRouters adds one or more routers to the current router.
//
// This can be used to add routes to multiple routers at once.
func (r *Router[D]) AddRouters(router *Router[D], routers ...*Router[D]) *Router[D] {
	r.children[""] = append(r.children[""], append(routers, router)...)
	return r
}

// Use adds a middleware function to the router.
//
// If nil is returned, the middleware is considered to have passed and the request is allowed to proceed.
// If an neoroute.NewError or a context response is returned it will be returned to the user.
// An error of a different type will be forwarded to the ErrorHandler it's return will be sent to the user.
//
// Middlewares will be executed in the order of root route to full route.
// And for every specific subroute, the middlewares will be executed in the order they are added.
//
// For for route1/route2/route3, the middlewares will be executed in the order of route1, route2, and route3.
func (r *Router[D]) Use(subroute string, middleware func(c *Ctx[D]) error) *Router[D] {
	subroute = cleanRoute(subroute)
	r.middlewares[subroute] = append(r.middlewares[subroute], middleware)
	return r
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

	// Read only the first bytes of the request header (we want to make sure it's not too long)
	headerReader := msgp.NewReaderSize(io.LimitReader(reqReader, maxRequestHeaderSize), maxRequestHeaderSize)
	id, route, size, leftover, err := readRequestHeader(headerReader)
	if err != nil {
		return messageResponse(c.respondError(r.config.RunErrorHandler(err, c))), nil
	}

	route = cleanRoute(route)
	c.id = id
	c.route = route

	// Check the size of the request before reading the data
	routeData, exists := r.actualRoutes[route]
	if !exists {
		return messageResponse(c.respondError(r.config.RunErrorHandler(ErrRouteDoesntExist, c))), nil
	}
	if limit := r.maxRequestSizeFor(routeData); limit > 0 && size > limit {
		return messageResponse(c.respondError(r.config.RunErrorHandler(ErrRequestTooLarge, c))), nil
	}

	// Body reader: buffered leftover from the header first, then the rest of the request, capped at exactly size bytes
	rawBody := io.MultiReader(bytes.NewReader(leftover), reqReader)
	data, err := readRequestData(io.LimitReader(rawBody, size), size)
	if err != nil {
		return messageResponse(c.respondError(r.config.RunErrorHandler(err, c))), nil
	}

	// The body must end exactly at size, any further byte means the request is malformed
	if err := ensureBodyEnded(rawBody); err != nil {
		return messageResponse(c.respondError(r.config.RunErrorHandler(err, c))), nil
	}
	c.reqData = data

	resp, runAfter := r.handleRoute(c, routeData)
	return resp, runAfter
}

// maxRequestSizeFor returns the max request size of a route, falling back to the router default if not set on the route.
func (r *Router[D]) maxRequestSizeFor(routeData exportedRoute[D]) int64 {
	if routeData.options.MaxRequestSize < 0 {
		return r.config.MaxRequestSize
	}
	return routeData.options.MaxRequestSize
}

// handleRoute runs the middlewares and the handler of a route and returns the response.
func (r *Router[D]) handleRoute(c *Ctx[D], routeData exportedRoute[D]) ([]byte, []func()) {
	// Run middlewares
	for _, middleware := range routeData.middlewares {
		if err := middleware(c); err != nil {

			// If route is of type NoResponse, return nothing.
			if !routeData.hasError {
				return nil, c.runAfter
			}

			// Check if error is user error
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

			return messageResponse(c.respondError(r.config.RunErrorHandler(err, c))), c.runAfter
		}
	}

	// Handle request
	err := routeData.handler(c)
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
