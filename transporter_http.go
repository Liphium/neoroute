package neoroute

import (
	"io"
	"net/http"
)

type HTTPTransporter[D any] struct {
	router *NeoRouter[D]
}

var _ Transporter[any] = &HTTPTransporter[any]{}

// NewHTTPTransporter creates a new HTTP transporter with the given handshake function and returns it along with an http.HandlerFunc that can be used to handle incoming HTTP requests.
//
// If session returned by the handshake function is nil, a new session will be created with a unique id. The data can then be set in the EnterNetworkFunc.
// If the bool is false, the handshake will be considered failed and the connection will be rejected.
func NewHTTPTransporter[D any](handshake func(r *http.Request) (*Session[D], bool)) (http.HandlerFunc, *HTTPTransporter[D]) {
	transporter := &HTTPTransporter[D]{
		router: nil,
	}
	hook := func(w http.ResponseWriter, r *http.Request) {

		// Route request
		if transporter.router == nil {
			http.Error(w, "Router not set.", http.StatusInternalServerError)
			return
		}

		// Perform handshake to get session data
		session, ok := handshake(r)
		if !ok {
			http.Error(w, "Handshake failed.", http.StatusUnauthorized)
			return
		}

		// Create session if handshake did not return one
		if session == nil {
			session = NewSession[D](transporter.router.config.runUUIDGenerator())
		}

		// Read body data
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body.", http.StatusInternalServerError)
			return
		}

		// Send response
		w.WriteHeader(http.StatusOK)
		resp, runAfter := transporter.router.handle(body, session)
		defer func() {
			for _, fn := range runAfter {
				fn()
			}
		}()
		if resp != nil {
			_, err = w.Write(resp)
			if err != nil {
				logger.Info("failed to send http response", "err", err)
			}
		}
	}

	return hook, transporter
}

func (t *HTTPTransporter[D]) SetRouter(r *NeoRouter[D]) {
	t.router = r
}
