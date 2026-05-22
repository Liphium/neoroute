package neoroute

type RouteRouter[D any] struct {
	neo   Router[D]
	route string
}

func (m *RouteRouter[D]) Group(route string) Router[D] {
	return &Group[D]{
		neo:    m.neo,
		prefix: route,
		parent: m,
	}

}

func (r *RouteRouter[D]) getRoute() string {
	return r.route
}

func (m *RouteRouter[D]) getNeo() *NeoRouter[D] {
	return m.neo.getNeo()
}
