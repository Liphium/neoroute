package client

import "log/slog"

type Ctx struct {
	data []byte // data field from Request struct
	name string // event name
}

func (c *Ctx) Data() []byte {
	return c.data
}

func (c *Ctx) Name() string {
	return c.name
}

// Global logger used for logging in the client package. It can be set using SetLogger function.
var Logger *slog.Logger = slog.Default()

// SetLogger sets the global logger for the client package. If not set, it defaults to slog.Default().
func SetLogger(l *slog.Logger) {
	Logger = l
}
