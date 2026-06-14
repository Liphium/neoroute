package neoroute

import (
	"log/slog"

	"github.com/google/uuid"
)

type Config struct {

	// If an error is returned from a route, this function will be called with the error
	// and should return a string that will be sent back to the client.
	// If nil, a default error message will be sent and the error will be logged.
	ErrorHandler func(err error) string

	// Is used to generate unique ids for sessions if no session
	// is provided in handshake. Must be thread safe.
	// Default is uuid.NewString() from https://github.com/google/uuid.
	UUIDGenerator func() string
}

func (cfg Config) runErrorHandler(err error) string {
	if cfg.ErrorHandler == nil {
		slog.Info("ErrorHandler is not set in config. An route returned an error", "error", err)
		return "Internal Server Error"
	}
	return cfg.ErrorHandler(err)
}

func (cfg Config) runUUIDGenerator() string {
	if cfg.UUIDGenerator == nil {
		return uuid.NewString()
	}
	return cfg.UUIDGenerator()
}
