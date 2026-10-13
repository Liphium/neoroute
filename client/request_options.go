package client

import "time"

// RequestOptions configures how a request is sent. They are used as defaults in Config and can be overridden per request.
type RequestOptions struct {
	// RequestTimeout is the max time to wait for a response. Defaults to 5 seconds if not set.
	RequestTimeout time.Duration
}

// defaultRequestOptions returns the request options with the library defaults applied.
func defaultRequestOptions() RequestOptions {
	return RequestOptions{RequestTimeout: 5 * time.Second}
}

// newRequestOptions returns the options for a single request, based on the client defaults with the given options applied in order.
func (c *Client) newRequestOptions(opts ...func(*RequestOptions)) RequestOptions {
	options := c.getConfig().RequestOptions
	for _, opt := range opts {
		opt(&options)
	}
	return options
}

// WithMaxRequestTimeout sets the max time to wait for a response for a request.
func WithMaxRequestTimeout(timeout time.Duration) func(*RequestOptions) {
	return func(o *RequestOptions) {
		o.RequestTimeout = timeout
	}
}
