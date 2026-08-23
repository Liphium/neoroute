package neoroute

import "errors"

var (
	ErrHandshakeFailed      = errors.New("handshake failed")
	ErrReadingBody          = errors.New("failed to read request body")
	ErrInvalidRequestFormat = errors.New("invalid request format")
	ErrRouteDoesntExist     = errors.New("this route does not exist")
	ErrMiddlewareDenied     = errors.New("middleware denied the request")
)
