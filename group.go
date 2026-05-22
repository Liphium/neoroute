package neoroute

type Group[D any] struct {
	neo    Router[D]
	prefix string
	parent Router[D]
}

func (m *Group[D]) Group(route string) Router[D] {
	return &Group[D]{
		neo:    m.neo,
		prefix: route,
		parent: m,
	}

}

func (r *Group[D]) getRoute() string {
	if r.parent == nil {
		return r.prefix
	}
	return r.parent.getRoute() + string(RouteSeparator) + r.prefix
}

func (m *Group[D]) getNeo() *NeoRouter[D] {
	return m.neo.getNeo()
}
