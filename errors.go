package neoroute

import (
	"errors"
	"fmt"
)

// TransporterError is returned to the error handler of routers when there is some kind of failure on the transporter level (e.g. WebSocket request is not an upgrade).
type TransporterError struct {
	ActualError error
}

func (e *TransporterError) Error() string {
	return fmt.Sprintf("transporter failure: %v", e.ActualError)
}

// IsTransporterFailure checks if the given error is a TransporterError.
func IsTransporterFailure(err error) (*TransporterError, bool) {
	if tErr, ok := errors.AsType[*TransporterError](err); ok {
		return tErr, true
	}
	return nil, false
}

var (
	// ErrHandshakeFailed is returned to the error handler of routers when the handshake method of the transporter fails.
	ErrHandshakeFailed = errors.New("handshake failed")
	// ErrReadingBody is returned to the error handler of routers when the router fails to read the request body.
	ErrReadingBody = errors.New("failed to read request body")
	// ErrInvalidRequestFormat is returned to the error handler of routers when a request is received with an invalid format.
	ErrInvalidRequestFormat = errors.New("invalid request format")
	// ErrRouteDoesntExist is returned to the error handler of routers when the router receives a message that is completely invalid in format.
	ErrRouteDoesntExist = errors.New("this route does not exist")
	// ErrRequestTooLarge is returned by Router.Handle when the request exceeds the max request size of its route. The transporter should close the connection.
	ErrRequestTooLarge = errors.New("request too large")
)
