package neoroute

func Use[D any](r Router[D], route string, middleware func(c *Ctx[D]) bool) {
	route = cleanRoute(r.getRoute() + string(RouteSeparator) + route)

	neos := r.getNeos()
	for _, neo := range neos {
		neo.middleware[route] = func(c *Ctx[D]) bool {

			// Run middleware
			return middleware(c)
		}
	}
}
