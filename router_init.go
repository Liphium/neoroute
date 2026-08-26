package neoroute

import (
	"fmt"
)

// Init initializes the router by building the map of actual routes. You don't have to call this, it's already done by the transporter.
func (r *Router[D]) Init() {
	// When this router has a route, init can not be called on it
	if r.hasRoute {
		panic("init of router that is not a root router (has a route), if you are seeing this please make sure you are not mounting a router returned by Route() into a transporter, please create a new router instead or use a group")
	}

	r.actualRoutes = r.buildMap()
}

// buildMap buildes a map of extracted routes, it makes sure all middlewares are added to the routes and that no duplicate routes exist.
func (r *Router[D]) buildMap() map[string]exportedRoute[D] {
	current := map[string]exportedRoute[D]{}
	for route, routers := range r.children {
		for _, router := range routers {
			subMap := router.buildMap()

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
	if r.hasRoute {
		current[""] = exportedRoute[D]{
			middlewares: r.middlewares[""], // Only middlewares at "" apply to this anyway
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
