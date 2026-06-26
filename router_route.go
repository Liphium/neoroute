package neoroute

type RouteRouter[D any] struct {
	neos  []*NeoRouter[D]
	route string
}

func (r *RouteRouter[D]) Group(route string) Router[D] {
	return &Group[D]{
		neos:   r.neos,
		prefix: route,
		parent: r,
	}

}

func (r *RouteRouter[D]) BuildSchema() map[string]RouteData[D] {
	return nil
}

func (r *RouteRouter[D]) AddRouters(router *NeoRouter[D], routers ...*NeoRouter[D]) Router[D] {
	r.neos = append(r.neos, append([]*NeoRouter[D]{router}, routers...)...)
	return r
}

func (r *RouteRouter[D]) Use(route string, middleware func(c *Ctx[D]) bool) {
	for _, neo := range r.neos {
		neo.Use(route, middleware)
	}
}

func (r *RouteRouter[D]) getRoute() string {
	return r.route
}

func (r *RouteRouter[D]) getNeos() []*NeoRouter[D] {
	return r.neos
}
