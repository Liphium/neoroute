package client

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"runtime/debug"
	"sync"
	"time"

	"github.com/coder/websocket"
)

type WebSocketTransporter struct {
	conn      *websocket.Conn
	done      chan struct{}
	receiver  *Receiver
	sendMutex *sync.Mutex
}

func NewWebSocketTransporter(r *Receiver) *WebSocketTransporter {

	return &WebSocketTransporter{
		receiver:  r,
		sendMutex: &sync.Mutex{},
	}
}

func (w *WebSocketTransporter) Connect(url *url.URL) (chan struct{}, error) {
	w.done = make(chan struct{})

	// Connect to server
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	conn, _, err := websocket.Dial(ctx, url.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to websocket server: %v", err)
	}
	w.conn = conn
	w.receiver.setSendFunc(func(data []byte) error {
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
			logger.Error("there was an error with the connection", "err", err.(error))
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
				logger.Info("websocket connection closed by remote",
					"code", closeErr.Code,
					"reason", closeErr.Reason,
				)
				return
			}

			logger.Error("error reading message", "err", err)
			return
		}

		if messageType != websocket.MessageBinary {
			logger.Info("wrong message type", "type", messageType)
			return
		}

		// Let receiver handle message
		go w.receiver.handle(msg)
	}
}
