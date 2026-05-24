package neoroute

type RouterGroup[D any] struct {
	neos []*NeoRouter[D]
}

func NewRouterGroup[D any](neo *NeoRouter[D], neos ...*NeoRouter[D]) Router[D] {
	return &RouterGroup[D]{
		neos: append(neos, neo),
	}
}

func (m *RouterGroup[D]) Group(route string) Router[D] {
	return &Group[D]{
		neos:   m.neos,
		prefix: route,
		parent: m,
	}

}

func (r *RouterGroup[D]) getRoute() string {
	return ""
}

func (m *RouterGroup[D]) getNeos() []*NeoRouter[D] {
	return m.neos
}
