package websocket

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"runtime/debug"
	"sync"
	"time"

	"github.com/Liphium/neoroute/client"

	ws "github.com/coder/websocket"
)

type WebSocketTransporter struct {
	conn      *ws.Conn
	done      chan struct{}
	client    *client.Client
	sendMutex sync.Mutex
}

// NewWebSocketTransporter creates a new transporter for connecting to the server over WebSocket.
func NewWebSocketTransporter(c *client.Client) *WebSocketTransporter {
	return &WebSocketTransporter{
		client:    c,
		sendMutex: sync.Mutex{},
	}
}

// Connect establishes a connection to the server.
//
// The returned channel will get a struct sent through it when the connection is closed.
func (w *WebSocketTransporter) Connect(url *url.URL) (chan struct{}, error) {
	w.done = make(chan struct{})

	// Connect to server
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	conn, resp, err := ws.Dial(ctx, url.String(), nil)
	defer resp.Body.Close()
	if err != nil {

		if resp != nil {
			bodyBytes, readErr := io.ReadAll(resp.Body)
			if readErr != nil {
				return nil, fmt.Errorf("failed to read response body: %v (dialErr: %v)", readErr, err)
			}

			// Check for transporter errors
			if resp.StatusCode != http.StatusOK {
				return nil, errors.New("received non ok status " + resp.Status + ": " + string(bodyBytes))
			}
		}

		return nil, fmt.Errorf("failed to connect to websocket server: %v", err)
	}
	w.conn = conn
	w.client.SetSendFunc(func(data []byte) error {
		w.sendMutex.Lock()
		defer w.sendMutex.Unlock()
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
		defer cancel()
		return w.conn.Write(ctx, ws.MessageBinary, data)
	})
	go w.ws(conn)
	return w.done, nil
}

// Close closes the WebSocket connection.
func (w *WebSocketTransporter) Close() error {
	return w.conn.Close(ws.StatusNormalClosure, "")
}

func (w *WebSocketTransporter) ws(conn *ws.Conn) {

	// Get server URL with any handshake parameters applied.

	defer func() {
		defer close(w.done)
		if err := recover(); err != nil {
			debug.PrintStack()
			w.Close()
			client.Logger.Error("there was an error with the connection", "err", err.(error))
			return
		}

		// Close the connection
		defer w.Close()
	}()

	for {
		messageType, msg, err := conn.Read(context.Background())
		if err != nil {
			if closeErr, ok := errors.AsType[ws.CloseError](err); ok {
				client.Logger.Info("websocket connection closed by remote",
					"code", closeErr.Code,
					"reason", closeErr.Reason,
				)
				return
			}

			client.Logger.Error("error reading message", "err", err)
			return
		}

		if messageType != ws.MessageBinary {
			client.Logger.Info("wrong message type", "type", messageType)
			return
		}

		// Let receiver handle message
		go w.client.Handle(msg)
	}
}
