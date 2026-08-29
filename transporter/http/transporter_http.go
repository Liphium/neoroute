package http_transporter

import (
	"net/http"

	"github.com/Liphium/neoroute"
	"github.com/Liphium/neoroute/neoschema"
	"github.com/google/uuid"
)

var _ neoschema.Transporter = &Transporter[any]{}

// Transporter is an HTTP transporter that handles incoming HTTP requests.
type Transporter[D any] struct {
	router *neoroute.Router[D]
}

// Config holds the configuration for the HTTP transporter.
type Config[D any] struct {

	// If session returned by the handshake function is nil, a new session will be created with a unique id. The data can then be set in the EnterNetworkFunc.
	//
	// If the bool is false, the handshake will be considered failed and the connection will be rejected.
	HandshakeFunc func(*http.Request) (D, bool)
}

// GetRegistries implements [neoschema.Transporter].
func (h *Transporter[D]) GetRegistries() []neoroute.IEventRegistry {
	return []neoroute.IEventRegistry{} // No events over HTTP
}

// GetSchema implements [neoschema.Transporter].
func (h *Transporter[D]) GetSchema() map[string]neoschema.RequestResponse {
	return neoschema.ToRouteSchema(h.router.GetRoutes())
}

// Type implements [neoschema.Transporter].
func (h *Transporter[D]) Type() neoschema.TransporterType {
	return neoschema.TransporterHTTP
}

// NewTransporter creates a new HTTP transporter with the given handshake function and returns it along with an http.HandlerFunc that can be used to handle incoming HTTP requests.
func NewTransporter[D any](router *neoroute.Router[D], config Config[D]) (http.HandlerFunc, *Transporter[D]) {
	transporter := &Transporter[D]{
		router: router,
	}
	router.Init()
	hook := func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if rec := recover(); rec != nil {
				neoroute.PrintRecoveredPanic("HTTP", rec)
			}
		}()

		// Perform handshake to get session data
		sessionData, ok := config.HandshakeFunc(r)
		if !ok {
			http.Error(w, router.Config().RunErrorHandler(neoroute.ErrHandshakeFailed, nil), http.StatusUnauthorized)
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
