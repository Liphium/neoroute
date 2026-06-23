package websocket_transporter

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

	"github.com/coder/websocket"
)

type WebSocketTransporter struct {
	conn      *websocket.Conn
	done      chan struct{}
	receiver  *client.Receiver
	sendMutex sync.Mutex
}

func NewWebSocketTransporter(r *client.Receiver) *WebSocketTransporter {

	return &WebSocketTransporter{
		receiver:  r,
		sendMutex: sync.Mutex{},
	}
}

func (w *WebSocketTransporter) Connect(url *url.URL) (chan struct{}, error) {
	w.done = make(chan struct{})

	// Connect to server
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	conn, resp, err := websocket.Dial(ctx, url.String(), nil)
	if err != nil {

		if resp != nil {
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, fmt.Errorf("failed to read response body: %v", err)
			}

			// Check for transporter errors
			if resp.StatusCode != http.StatusOK {
				return nil, errors.New("received non ok status " + resp.Status + ": " + string(bodyBytes))
			}
		}

		return nil, fmt.Errorf("failed to connect to websocket server: %v", err)
	}
	w.conn = conn
	w.receiver.SetSendFunc(func(data []byte) error {
		w.sendMutex.Lock()
		defer w.sendMutex.Unlock()
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
		defer cancel()
		return w.conn.Write(ctx, websocket.MessageBinary, data)
	})
	go w.ws(conn)
	return w.done, nil
}

func (w *WebSocketTransporter) Close() error {
	return w.conn.Close(websocket.StatusNormalClosure, "")
}

func (w *WebSocketTransporter) ws(conn *websocket.Conn) {

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
			var closeErr websocket.CloseError
			if errors.As(err, &closeErr) {
				client.Logger.Info("websocket connection closed by remote",
					"code", closeErr.Code,
					"reason", closeErr.Reason,
				)
				return
			}

			client.Logger.Error("error reading message", "err", err)
			return
		}

		if messageType != websocket.MessageBinary {
			client.Logger.Info("wrong message type", "type", messageType)
			return
		}

		// Let receiver handle message
		go w.receiver.Handle(msg)
	}
}
