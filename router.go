package neoroute

import (
	"errors"
)

type Router[D any] interface {
	Group(route string) Router[D]
	AddRouters(router *NeoRouter[D], routers ...*NeoRouter[D]) Router[D]
	Use(route string, middleware func(c *Ctx[D]) bool)
	getRoute() string
	getNeos() []*NeoRouter[D]
}

type NeoRouter[D any] struct {
	neos       []*NeoRouter[D]
	routes     map[string]func(c *Ctx[D]) error
	middleware map[string]func(c *Ctx[D]) bool
	config     Config
}

func NewNeoRouter[D any](config Config) *NeoRouter[D] {
	return &NeoRouter[D]{
		routes:     make(map[string]func(c *Ctx[D]) error),
		middleware: make(map[string]func(c *Ctx[D]) bool),
		config:     config,
	}
}

func (r *NeoRouter[D]) Group(route string) Router[D] {
	return &Group[D]{
		neos:   []*NeoRouter[D]{r},
		prefix: route,
		parent: nil,
	}
}

func (r *NeoRouter[D]) AddRouters(router *NeoRouter[D], routers ...*NeoRouter[D]) Router[D] {
	r.neos = append(r.neos, append([]*NeoRouter[D]{router}, routers...)...)
	return r
}

func (r *NeoRouter[D]) Use(route string, middleware func(c *Ctx[D]) bool) {
	route = cleanRoute(r.getRoute() + string(RouteSeparator) + route)

	neos := r.getNeos()
	for _, neo := range neos {
		neo.middleware[route] = func(c *Ctx[D]) bool {

			// Run middleware
			return middleware(c)
		}
	}
}

func (r *NeoRouter[D]) getRoute() string {
	return ""
}

func (r *NeoRouter[D]) getNeos() []*NeoRouter[D] {
	return []*NeoRouter[D]{r}
}

func (r *NeoRouter[D]) handle(reqData []byte, session *Session[D]) ([]byte, []func()) {

	c := &Ctx[D]{
		neo:     r,
		id:      -1,
		reqData: []byte{},
		route:   "",
		session: session,
	}

	var data request
	_, err := data.UnmarshalMsg(reqData)
	if err != nil {
		logger.Info("failed to unmarshal request", "err", err)
		return messageResponse(c.respondError(ErrInvalidRequestFormat)), nil
	}

	route := cleanRoute(data.Route)

	c.id = data.Id
	c.reqData = data.Data
	c.route = route

	// Check if handler for route exists
	handler, exists := r.routes[route]
	if !exists {
		return messageResponse(c.respondError(ErrRouteNotExists)), nil
	}

	// Run middlewares
	subRoutes := buildSubroutes(route)
	for _, subroute := range subRoutes {
		if middleware, ok := r.middleware[subroute]; ok {
			if !middleware(c) {
				return messageResponse(c.respondError(ErrMiddlewareDenied)), nil
			}
		}
	}

	// Handle request
	err = handler(c) // TODO: add panic protection
	if err == nil {

		// Handlers never should return nil.
		panic("handler should always return something")
	} else if errors.Is(err, responseData{}) {

		// Return response from handler
		respData := err.(responseData)
		resp := response{
			Id:           c.id,
			responseData: respData,
		}
		return messageResponse(resp), c.runAfter
	} else if errors.Is(err, noResponse{}) {

		// Return no response
		return nil, c.runAfter
	} else {

		// Let user handle the error and decide what error message to send back to the client
		return messageResponse(c.respondError(handleError(r.config, err))), c.runAfter
	}
}
