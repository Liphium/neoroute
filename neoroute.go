// This package is the core server package for the Neoroute Web Server Framework.
package neoroute

import "log/slog"

// Central logger that all of neoroute, including the official transporter packages, will use.
var Logger *slog.Logger = slog.Default()

// SetLogger sets the logger neoroute will use. You can put in a logger with custom formatting, etc. here.
func SetLogger(l *slog.Logger) {
	Logger = l
}

// This can be set as the SessionData type if you don't want to use any session data. It is just a struct with no fields.
type NoData struct{}
