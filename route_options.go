package neoroute

const DefaultMaxRequestSize = 2 * 1024 * 1024 // 2 MB

// RouteOptions configures a single route. Options are applied on top of the defaults when a route is registered.
type RouteOptions struct {
	// MaxRequestSize is the max size of the request data in bytes. 0 means unlimited.
	//
	// A negative value means "not set" and the router default from Config is used.
	MaxRequestSize int64
}

func (r RouteOptions) setDefaultOptions() {
	if r.MaxRequestSize == 0 {
		r.MaxRequestSize = DefaultMaxRequestSize
	}
}

// getRouteOptions returns the base route options with all given options applied in order.
func getRouteOptions(base RouteOptions, opts ...func(*RouteOptions)) RouteOptions {
	for _, opt := range opts {
		opt(&base)
	}
	return base
}

// WithMaxRequestSize sets the max request size in bytes for a route. 0 means unlimited.
func WithMaxRequestSize(size int64) func(*RouteOptions) {
	return func(o *RouteOptions) {
		o.MaxRequestSize = size
	}
}
