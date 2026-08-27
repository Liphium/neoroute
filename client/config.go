package client

import (
	"time"
)

// Config holds the configuration for the client.
type Config struct {
	ErrorHandler   func(err error)
	RequestTimeout time.Duration
}

// RunErrorHandler runs the error handler function if one is set. It logs the error if no handler is set.
func (cfg Config) RunErrorHandler(err error) {
	if cfg.ErrorHandler == nil {
		Logger.Info("ErrorHandler is not set in config. An error occurred", "error", err)
		return
	}
	cfg.ErrorHandler(err)
}
