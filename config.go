package neoroute

type Config[D any] struct {

	// If an error is returned from a route, this function will be called with the error and the context (which includes the session) and should return a string that will be sent back to the client.
	//
	// If nil, a default error message will be sent and the error will be logged.
	//
	// c can be nil when called from places where there is no context.
	ErrorHandler func(err error, c *Ctx[D]) string
}

// RunErrorHandler calls the configured ErrorHandler, or returns a default message if none is set.
func (cfg Config[D]) RunErrorHandler(err error, c *Ctx[D]) string {
	if cfg.ErrorHandler == nil {
		Logger.Info("ErrorHandler is not set in config. A route returned an error", "error", err)
		return "Internal Server Error"
	}
	return cfg.ErrorHandler(err, c)
}
