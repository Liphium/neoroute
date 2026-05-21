package neoroute

type Mount struct {
	neo    Router
	prefix string
	parent *Mount
}

func (m *Mount) Use() {

}

func (m *Mount) Mount(route string) Router {
	return &Mount{
		neo:    m.neo,
		prefix: route,
		parent: m,
	}

}

func (r *Mount) getRoute() string {
	if r.parent == nil {
		return r.prefix
	}
	return r.parent.getRoute() + r.prefix
}

func (m *Mount) getNeo() *NeoRouter {
	return m.neo.getNeo()
}
