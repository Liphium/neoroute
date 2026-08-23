package client

import (
	"time"
)

type Config struct {
	ErrorHandler   func(err error)
	RequestTimeout time.Duration
}

func (cfg Config) RunErrorHandler(err error) {
	if cfg.ErrorHandler == nil {
		Logger.Info("ErrorHandler is not set in config. An error occurred", "error", err)
		return
	}
	cfg.ErrorHandler(err)
}
