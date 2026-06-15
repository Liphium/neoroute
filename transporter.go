package neoroute

import (
	"net/http"
)

type Transporter[D any] interface {
	SetRouter(r *NeoRouter[D])
}

// HandshakeFunc should perform a handshake for new connections.
// If the handshake is successful, it should return a session data and true.
// If the handshake fails, it should return false.
type HandshakeFunc[D any] func(r *http.Request) (D, bool)
