package neoroute

type Group[D any] struct {
	neos   []*NeoRouter[D]
	prefix string
	parent Router[D]
}

func (m *Group[D]) Group(route string) Router[D] {
	return &Group[D]{
		neos:   m.neos,
		prefix: route,
		parent: m,
	}
}

func (m *Group[D]) AddRouters(router *NeoRouter[D], routers ...*NeoRouter[D]) Router[D] {
	m.neos = append(m.neos, append([]*NeoRouter[D]{router}, routers...)...)
	return m
}

func (m *Group[D]) Use(route string, middleware func(c *Ctx[D]) bool) {
	for _, neo := range m.neos {
		neo.Use(route, middleware)
	}
}

func (m *Group[D]) getRoute() string {
	if m.parent == nil {
		return m.prefix
	}
	return m.parent.getRoute() + string(RouteSeparator) + m.prefix
}

func (m *Group[D]) getNeos(_ ...*NeoRouter[D]) []*NeoRouter[D] {
	return m.neos
}
