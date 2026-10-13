package neoroute

// RouteOptions configures a single route. Options are applied on top of the defaults when a route is registered.
type RouteOptions struct {
	// MaxRequestSize is the max size of the request data in bytes. 0 means unlimited.
	//
	// A negative value means "not set" and the router default from Config is used.
	MaxRequestSize int64
}

// inheritRouteOption marks a route option as not set, so the router default is used.
const inheritRouteOption = -1

// newRouteOptions returns the default route options with all given options applied in order.
func newRouteOptions(opts ...func(*RouteOptions)) RouteOptions {
	options := RouteOptions{MaxRequestSize: inheritRouteOption}
	for _, opt := range opts {
		opt(&options)
	}
	return options
}

// WithMaxRequestSize sets the max request size in bytes for a route. 0 means unlimited.
func WithMaxRequestSize(size int64) func(*RouteOptions) {
	return func(o *RouteOptions) {
		o.MaxRequestSize = size
	}
}
