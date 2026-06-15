package neoroute

import (
	"log/slog"
)

type Config struct {

	// If an error is returned from a route, this function will be called with the error
	// and should return a string that will be sent back to the client.
	// If nil, a default error message will be sent and the error will be logged.
	ErrorHandler func(err error) string
}

func (cfg Config) RunErrorHandler(err error) string {
	if cfg.ErrorHandler == nil {
		slog.Info("ErrorHandler is not set in config. An route returned an error", "error", err)
		return "Internal Server Error"
	}
	return cfg.ErrorHandler(err)
}
