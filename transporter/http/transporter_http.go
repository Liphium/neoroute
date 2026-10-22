package http_transporter

import (
	"io"
	"net/http"
	"time"

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
	// HeaderReadTimeout limits reading the Neoroute envelope before the route is known. Zero defaults to ten seconds; a negative value disables it. Body timeouts use neoroute.RouteOptions.ReadTimeout.
	//
	// Configure http.Server.ReadHeaderTimeout separately.
	//
	// ResponseWriter middleware must support SetReadDeadline or expose Unwrap; requests are rejected if the deadline cannot be set.
	HeaderReadTimeout time.Duration

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
	if config.HeaderReadTimeout == 0 {
		config.HeaderReadTimeout = 10 * time.Second
	}
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

		reader := &requestReader{Reader: r.Body, controller: http.NewResponseController(w)}
		if config.HeaderReadTimeout > 0 {
			if err := reader.SetReadTimeout(config.HeaderReadTimeout); err != nil {
				neoroute.Logger.Error("failed to set HTTP read deadline", "err", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		// Perform handshake to get session data
		sessionData, ok := config.HandshakeFunc(r)
		if !ok {
			http.Error(w, router.Config().RunErrorHandler(neoroute.ErrHandshakeFailed, nil), http.StatusUnauthorized)
			return
		}

		// Create session with handshake data
		session := neoroute.NewSession(uuid.NewString(), sessionData, neoroute.SessionTransporterCallbacks[D]{})

		defer r.Body.Close()

		resp, runAfter, err := transporter.router.Handle(reader, session)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Send response
		w.WriteHeader(http.StatusOK)
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

type requestReader struct {
	io.Reader
	controller *http.ResponseController
}

func (r *requestReader) SetReadTimeout(timeout time.Duration) error {
	var deadline time.Time
	if timeout > 0 {
		deadline = time.Now().Add(timeout)
	}
	return r.controller.SetReadDeadline(deadline)
}
