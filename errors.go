package neoroute

import "log/slog"

const (
	ErrRouteNotExists       = "route does not exist"
	ErrInvalidRequestFormat = "invalid request format"
	ErrMiddlewareDenied     = "middleware denied the request"
)

func handleError(cfg Config, err error) string {
	if cfg.ErrorHandler == nil {
		slog.Info("ErrorHandler is not set in config. An route returned an error", "error", err)
		return "Internal Server Error"
	}
	return cfg.ErrorHandler(err)
}
