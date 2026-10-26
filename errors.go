package neoroute

import "errors"

var (
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
