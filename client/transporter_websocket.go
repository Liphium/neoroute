package client

import (
	"fmt"
	"net/url"
	"runtime/debug"
	"sync"
	"time"

	"github.com/gorilla/websocket"
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
	conn, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to websocket server: %v", err)
	}
	w.conn = conn
	w.receiver.setSendFunc(func(data []byte) error {
		w.sendMutex.Lock()
		defer w.sendMutex.Unlock()
		return w.conn.WriteMessage(websocket.BinaryMessage, data)
	})
	go w.ws(conn)
	return w.done, nil
}

func (w *WebSocketTransporter) Close() error {
	closePayload := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")

	err := w.conn.WriteControl(
		websocket.CloseMessage,
		closePayload,
		time.Now().Add(time.Second*2),
	)

	time.Sleep(500 * time.Millisecond)
	return err
}

func (w *WebSocketTransporter) ws(conn *websocket.Conn) {

	// Get server URL with any handshake parameters applied.

	defer func() {
		defer close(w.done)
		if err := recover(); err != nil {
			debug.PrintStack()
			conn.Close()
			logger.Error("there was an error with the connection", "err", err.(error))
			return
		}

		// Close the connection
		defer conn.Close()
	}()

	for {
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				return
			}

			logger.Info("error reading message", "err", err)
			return
		}

		if messageType != websocket.BinaryMessage {
			logger.Info("wrong message type", "type", messageType)
			return
		}

		// Let receiver handle message
		go w.receiver.handle(msg)
	}
}
