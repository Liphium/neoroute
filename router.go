package neoroute

import (
	"errors"
	"fmt"
)

type Router[D any] interface {
	Group(route string) Router[D]
	getRoute() string
	getNeos() []*NeoRouter[D]
}

type NeoRouter[D any] struct {
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

func (r *NeoRouter[D]) getRoute() string {
	return ""
}

func (r *NeoRouter[D]) getNeos() []*NeoRouter[D] {
	return []*NeoRouter[D]{r}
}

func (r *NeoRouter[D]) handle(reqData []byte, session *Session[D]) []byte {

	c := &Ctx[D]{
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

	route := cleanRoute(data.Route)

	c.id = data.Id
	c.data = data.Data
	c.route = route

	// Check if handler for route exists
	handler, exists := r.routes[route]
	if !exists {
		return messageResponse(r, c.respondError(fmt.Errorf("Route does not exist.")))
	}

	// Run middlewares
	subRoutes := buildSubroutes(route)
	for _, subroute := range subRoutes {
		if middleware, ok := r.middleware[subroute]; ok {
			if !middleware(c) {
				return messageResponse(r, c.respondError(fmt.Errorf("Middleware denied access.")))
			}
		}
	}

	// Handle request
	err = handler(c)
	if err == nil {

		// Handlers never should return nil.
		panic("handler should always return something")
	} else if errors.Is(err, response{}) {

		// Return response from handler
		return messageResponse(c.neo, err.(response))
	} else if errors.Is(err, noResponse{}) {

		// Return no response
		return nil
	} else {
		// Log error from handler and return generic error message to client
		logger.Info("an error occurred", "err", err)
		return messageResponse(r, c.respondError(fmt.Errorf("Internal server error.")))
	}

}
