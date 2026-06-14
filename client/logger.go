package client

import "log/slog"

var Logger *slog.Logger = slog.Default()

func SetLogger(l *slog.Logger) {
	Logger = l
}
