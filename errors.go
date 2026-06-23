package neoroute

const (
	ErrRouteNotExists       = "route does not exist"
	ErrHandshakeFailed      = "handshake failed"
	ErrReadingBody          = "failed to read request body"
	ErrInvalidRequestFormat = "invalid request format"
	ErrMiddlewareDenied     = "middleware denied the request"
)
