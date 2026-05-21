package neoroute

import (
	"fmt"
	"io"
	"net/http"
)

type HTTPTransporter struct {
	router *NeoRouter
}

type HTTPHook func(w http.ResponseWriter, r *http.Request)

func NewHTTPTransporter() (HTTPHook, Transporter) {
	transporter := &HTTPTransporter{
		router: nil,
	}
	hook := func(w http.ResponseWriter, r *http.Request) {

		// Route request
		if transporter.router == nil {
			http.Error(w, "Router not set", http.StatusInternalServerError)
			return
		} else {

			// Read body data
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error reading request body: %v", err), http.StatusInternalServerError)
				return
			}

			// Send response
			w.WriteHeader(http.StatusOK)
			_, err = w.Write(transporter.router.handle(body))
			if err != nil {
				logger.Info("failed to send http response: ", err)
			}
		}

	}

	return hook, transporter
}

func (t *HTTPTransporter) SetRouter(r *NeoRouter) {
	t.router = r
}
