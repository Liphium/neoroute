package neoroute

type Group struct {
	neo    Router
	prefix string
	parent *Group
}

func (m *Group) Use() {

}

func (m *Group) Group(route string) Router {
	return &Group{
		neo:    m.neo,
		prefix: route,
		parent: m,
	}

}

func (r *Group) getRoute() string {
	if r.parent == nil {
		return r.prefix
	}
	return r.parent.getRoute() + string(RouteSeparator) + r.prefix
}

func (m *Group) getNeo() *NeoRouter {
	return m.neo.getNeo()
}
