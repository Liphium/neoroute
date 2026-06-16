package neoroute

import "log/slog"

// ONLY USE THIS IN TRANSPORTER IMPLEMENTATIONS IF YOU
// DON'T WANT YOUR NEOROUTE LOGGING TO BE MIXED WITH OTHER LOGGING
var Logger *slog.Logger = slog.Default()

func SetLogger(l *slog.Logger) {
	Logger = l
}
