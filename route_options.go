package neoroute

import "time"

const DefaultMaxRequestSize = 2 * 1024 * 1024 // 2 MB

const DefaultReadTimeout = time.Minute

// RouteOptions configures a single route. Options are applied on top of the defaults when a route is registered.
type RouteOptions struct {
	// ReadTimeout limits body reads after the route header is decoded. Zero defaults
	// to one minute; a negative value disables the body timeout.
	ReadTimeout time.Duration

	// MaxRequestSize is the max size of the request data in bytes.
	//
	// -1 means unlimited.
	MaxRequestSize int64
}

// setDefaultOptions fills in the values that were not set.
func (r *RouteOptions) setDefaultOptions() {
	if r.ReadTimeout == 0 {
		r.ReadTimeout = DefaultReadTimeout
	}
	if r.MaxRequestSize == 0 {
		r.MaxRequestSize = DefaultMaxRequestSize
	}
}

// getRouteOptions returns the base route options with all given options applied in order.
func getRouteOptions(base RouteOptions, opts ...func(*RouteOptions)) RouteOptions {
	for _, opt := range opts {
		opt(&base)
	}
	base.setDefaultOptions()
	return base
}

// WithReadTimeout overrides the time allowed to read a route's body.
// Zero uses the default of one minute; a negative value disables the timeout.
func WithReadTimeout(timeout time.Duration) func(*RouteOptions) {
	return func(o *RouteOptions) { o.ReadTimeout = timeout }
}

// WithMaxRequestSize sets the max request size in bytes for a route. -1 means unlimited.
func WithMaxRequestSize(size int64) func(*RouteOptions) {
	return func(o *RouteOptions) {
		o.MaxRequestSize = size
	}
}
