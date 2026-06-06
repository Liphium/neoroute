package neoroute

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/coder/websocket"
)

type WebSocketTransporter[D any] struct {
	eventRegistries []*EventRegistry
	router          *NeoRouter[D]
	config          WSConfig[D]
	mutex           *sync.Mutex
	sessions        map[string]*wsSession[D]
}

type UpgradeFuncWS func(w http.ResponseWriter, r *http.Request, opts *websocket.AcceptOptions) (*websocket.Conn, error)

type WSConfig[D any] struct {
	UpgradeFunc          UpgradeFuncWS
	OverwriteSessionFunc func(id string) bool

	HandshakeFunc     func(r *http.Request) (*Session[D], bool)
	EnterNetworkFunc  func(session *Session[D], t *WebSocketTransporter[D])
	DisconnectHandler func(session *Session[D])
}

type wsSession[D any] struct {
	mutex     *sync.Mutex
	sendMutex *sync.Mutex
	ctx       context.Context
	cancel    context.CancelFunc
	conn      *websocket.Conn
	session   *Session[D]
}

var _ Transporter[any] = &WebSocketTransporter[any]{}

func NewWebSocketTransporter[D any](config WSConfig[D]) (http.HandlerFunc, *WebSocketTransporter[D]) {
	transporter := &WebSocketTransporter[D]{
		router:          nil,
		config:          config,
		sessions:        make(map[string]*wsSession[D]),
		mutex:           &sync.Mutex{},
		eventRegistries: []*EventRegistry{},
	}

	hook := func(w http.ResponseWriter, r *http.Request) {

		if transporter.router == nil {
			http.Error(w, "Router not set.", http.StatusInternalServerError)
			return
		}

		// Perform handshake to get session data
		userSession, ok := transporter.config.HandshakeFunc(r)
		if !ok {
			http.Error(w, "Handshake failed.", http.StatusUnauthorized)
			return
		}

		// Upgrade to WebSocket session
		conn, err := transporter.config.UpgradeFunc(w, r, nil)
		if err != nil {
			logger.Info("Upgrade to WebSocket failed", "err", err)
			return
		}

		// Add session to transporter
		session := transporter.addSession(userSession, conn)
		if session == nil {
			return
		}

		go transporter.handleSession(session)
	}

	return hook, transporter
}

// SetRouter sets the router for the transporter.
// This should be done before starting to listen for connections.
// This should only be done once and not changed later.
func (t *WebSocketTransporter[D]) SetRouter(r *NeoRouter[D]) {
	t.router = r
}

func (t *WebSocketTransporter[D]) AddEventRegistry(e *EventRegistry) {
	t.mutex.Lock()
	t.eventRegistries = append(t.eventRegistries, e)
	t.mutex.Unlock()
}

func (t *WebSocketTransporter[D]) addSession(userSession *Session[D], conn *websocket.Conn) *wsSession[D] {

	// Check if session already exists and if it should be overwritten
	t.mutex.Lock()
	if oldSession, exists := t.sessions[userSession.id]; exists {

		if t.config.OverwriteSessionFunc(userSession.id) {

			oldSession.mutex.Lock()

			// Close existing session before overwriting
			if oldSession.cancel != nil {
				oldSession.cancel() // Cancel old context if overwritten
			}
			if err := oldSession.conn.CloseNow(); err != nil {
				logger.Info("failed to close old connection", "session", userSession.id, "err", err)
			}

			oldSession.mutex.Unlock()

		} else {
			t.mutex.Unlock()
			return nil
		}
	}

	ctx, cancel := context.WithCancel(context.Background())

	// Create new session entry
	session := &wsSession[D]{
		mutex:     &sync.Mutex{},
		conn:      conn,
		session:   userSession,
		ctx:       ctx,
		sendMutex: &sync.Mutex{},
		cancel:    cancel,
	}
	t.sessions[userSession.id] = session
	t.mutex.Unlock()
	return session
}

func (t *WebSocketTransporter[D]) removeSession(id string) {
	t.mutex.Lock()
	delete(t.sessions, id)
	t.mutex.Unlock()
}

func (t *WebSocketTransporter[D]) Adapt(id string) (Adapter, error) {
	session, ok := t.getSession(id)
	if !ok {
		return nil, fmt.Errorf("session %s not found", id)
	}

	session.mutex.Lock()
	conn := session.conn
	sendMutex := session.sendMutex
	ctx := session.ctx
	session.mutex.Unlock()

	if conn == nil {
		return nil, fmt.Errorf("websocket session not set for %s", id)
	}

	adapter := &WebSocketAdapter{
		transporterType: "WebSocket",
		eventRegistries: t.eventRegistries,
		conn:            conn,
		mutex:           &sync.Mutex{},
		sendMutex:       sendMutex,
		ctx:             ctx,
	}
	go adapter.waitClosed()
	return adapter, nil
}

func (t *WebSocketTransporter[D]) getSession(id string) (*wsSession[D], bool) {
	t.mutex.Lock()
	session, ok := t.sessions[id]
	t.mutex.Unlock()
	return session, ok
}

func (t *WebSocketTransporter[D]) handleSession(session *wsSession[D]) {

	session.mutex.Lock()
	conn := session.conn
	userSession := session.session
	session.mutex.Unlock()

	defer func() {
		if session.cancel != nil {
			session.cancel()
		}
		conn.CloseNow()
		t.config.DisconnectHandler(session.session)
		session.mutex.Lock()
		t.removeSession(session.session.id)
		session.mutex.Unlock()
	}()

	t.config.EnterNetworkFunc(session.session, t)

	for {
		messageType, msg, err := conn.Read(context.Background())
		if err != nil {

			// Only log err if it is not due to expected connection closure
			var closeErr websocket.CloseError
			if errors.As(err, &closeErr) {
				logger.Info("websocket connection closed by remote",
					"code", closeErr.Code,
					"reason", closeErr.Reason,
				)
				return
			}

			return
		}

		if messageType != websocket.MessageBinary {
			return
		}

		// Handle request and send response back over the same connection
		resp, runAfter := t.router.handle(msg, userSession)
		defer func() {
			for _, fn := range runAfter {
				fn()
			}
		}()
		if resp != nil {
			go func() {
				session.sendMutex.Lock()
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
				defer cancel()
				err = conn.Write(ctx, websocket.MessageBinary, resp)
				session.sendMutex.Unlock()
				if err != nil {
					logger.Info("failed to send websocket response", "err", err)
					return
				}
			}()

		}
	}
}
