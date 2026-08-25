package neoroute

import "fmt"

// Init initializes the router by building the map of actual routes. You don't have to call this, it's already done by the transporter.
func (r *Router[D]) Init() {
	r.actualRoutes = r.buildMap(map[string]exportedRoute[D]{})
}

// buildMap buildes a map of extracted routes, it makes sure all middlewares are added to the routes and that no duplicate routes exist.
func (r *Router[D]) buildMap(current map[string]exportedRoute[D]) map[string]exportedRoute[D] {
	for route, routers := range r.children {
		for _, router := range routers {
			subMap := router.buildMap(current)

			for subRoute, routeData := range subMap {
				fullRoute := cleanRoute(route + string(RouteSeparator) + subRoute)

				// When a duplicate route is detected we want to panic as it may lead to unexpected behavior
				if _, ok := current[fullRoute]; ok {
					panic(fmt.Sprintf("duplicate route: %s", fullRoute))
				}

				// Find all the middlewares we have locally until the route we have now and add them
				middlewares := append(r.getMiddlewaresFor(fullRoute), routeData.middlewares...)

				current[fullRoute] = exportedRoute[D]{
					middlewares: middlewares,
					RouteData:   routeData.RouteData,
				}
			}
		}
	}

	// Add own route (if there even is one)
	if r.prefix != "" {
		ownRoute := cleanRoute(r.prefix)
		current[ownRoute] = exportedRoute[D]{
			middlewares: r.getMiddlewaresFor(ownRoute),
			RouteData:   r.route,
		}
	}

	return current
}

// getMiddlewaresFor returns all middlewares for a given route, including the middlewares of its parent routes (from the local router).
func (r Router[D]) getMiddlewaresFor(route string) []MiddlewareFunc[D] {
	middlewares := []MiddlewareFunc[D]{}
	for _, route := range buildSubroutes(route) {
		if middlewaresForRoute, ok := r.middlewares[route]; ok {
			middlewares = append(middlewaresForRoute, middlewares...)
		}
	}
	return middlewares
}
