package http

import (
	"net/http"

	"github.com/Liphium/neoroute"
	"github.com/Liphium/neoroute/neoschema"
	"github.com/google/uuid"
)

var _ neoschema.Transporter = &HTTPTransporter[any]{}

type HTTPTransporter[D any] struct {
	router *neoroute.NeoRouter[D]
}

// GetRegistries implements neoschema.Transporter.
func (h *HTTPTransporter[D]) GetRegistries() []neoroute.IEventRegistry {
	return []neoroute.IEventRegistry{} // No events over HTTP
}

// GetSchema implements neoschema.Transporter.
func (h *HTTPTransporter[D]) GetSchema() map[string]neoschema.RequestResponse {
	return neoschema.ToRouteSchema(h.router.GetRoutes())
}

// Type implements neoschema.Transporter.
func (h *HTTPTransporter[D]) Type() neoschema.TransporterType {
	return neoschema.TransporterHTTP
}

// NewHTTPTransporter creates a new HTTP transporter with the given handshake function and returns it along with an http.HandlerFunc that can be used to handle incoming HTTP requests.
//
// If session returned by the handshake function is nil, a new session will be created with a unique id. The data can then be set in the EnterNetworkFunc.
// If the bool is false, the handshake will be considered failed and the connection will be rejected.
func NewHTTPTransporter[D any](router *neoroute.NeoRouter[D], handshake neoroute.HandshakeFunc[D]) (http.HandlerFunc, *HTTPTransporter[D]) {
	transporter := &HTTPTransporter[D]{
		router: router,
	}
	hook := func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if rec := recover(); rec != nil {
				neoroute.PrintRecoveredPanic("HTTP", rec)
			}
		}()

		// Perform handshake to get session data
		sessionData, ok := handshake(r)
		if !ok {
			http.Error(w, neoroute.ErrHandshakeFailed, http.StatusUnauthorized)
			return
		}

		// Create session with handshake data
		session := neoroute.NewSession(uuid.NewString(), sessionData, neoroute.SessionTransporterCallbacks[D]{})

		defer r.Body.Close()

		// Send response
		w.WriteHeader(http.StatusOK)
		resp, runAfter := transporter.router.Handle(r.Body, session)
		defer func() {
			for _, fn := range runAfter {
				fn()
			}
		}()
		if resp != nil {
			_, err := w.Write(resp)
			if err != nil {
				neoroute.Logger.Info("failed to send http response", "err", err)
			}
		}
	}

	return hook, transporter
}
