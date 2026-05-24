package neoroute

type RouteRouter[D any] struct {
	neos  []*NeoRouter[D]
	route string
}

func (m *RouteRouter[D]) Group(route string) Router[D] {
	return &Group[D]{
		neos:   m.neos,
		prefix: route,
		parent: m,
	}

}

func (r *RouteRouter[D]) getRoute() string {
	return r.route
}

func (m *RouteRouter[D]) getNeos() []*NeoRouter[D] {
	return m.neos
}
