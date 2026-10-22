package websocket_transporter

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/Liphium/neoroute"
	"github.com/Liphium/neoroute/neoschema"
	"github.com/google/uuid"

	"github.com/coder/websocket"
)

var _ neoschema.Transporter = &Transporter[any]{}

// Transporter is the WebSocket transporter implementation.
type Transporter[D any] struct {
	eventRegistries []neoroute.IEventRegistry
	router          *neoroute.Router[D]
	config          Config[D]
	mutex           sync.Mutex
	sessions        map[string]*wsSession[D]
}

// GetRegistries implements [neoschema.Transporter].
func (t *Transporter[D]) GetRegistries() []neoroute.IEventRegistry {
	return t.eventRegistries
}

// GetSchema implements [neoschema.Transporter].
func (t *Transporter[D]) GetSchema() map[string]neoschema.RequestResponse {
	return neoschema.ToRouteSchema(t.router.GetRoutes())
}

// Type implements [neoschema.Transporter].
func (t *Transporter[D]) Type() neoschema.TransporterType {
	return neoschema.TransporterWebSocket
}

// Config holds the configuration for the WebSocket transporter.
type Config[D any] struct {
	// HeaderReadTimeout limits reading the Neoroute envelope after a data message arrives.
	//
	// Zero defaults to ten seconds; a negative value disables it. Idle time between messages is unlimited. Body timeouts use neoroute.RouteOptions.ReadTimeout.
	HeaderReadTimeout time.Duration

	// If session is nil, a new session will be created with a unique id. The data can then be set in the EnterNetworkFunc.
	//
	// If the bool is false, the handshake will be considered failed and the connection will be rejected.
	HandshakeFunc func(*http.Request) (D, bool)

	// Options for accepting the websocket connection, if not set, default options from the library are used.
	AcceptOptions *websocket.AcceptOptions

	EnterNetworkFunc  func(session *neoroute.Session[D])
	DisconnectHandler func(session *neoroute.Session[D])
}

type wsSession[D any] struct {
	mutex     sync.Mutex
	sendMutex *sync.Mutex
	ctx       context.Context
	cancel    context.CancelFunc
	conn      *websocket.Conn
	session   *neoroute.Session[D]
}

// NewTransporter returns a new WebSocket transporter.
func NewTransporter[D any](router *neoroute.Router[D], config Config[D]) (http.HandlerFunc, *Transporter[D]) {
	if config.HeaderReadTimeout == 0 {
		config.HeaderReadTimeout = 10 * time.Second
	}
	transporter := &Transporter[D]{
		router:          router,
		config:          config,
		sessions:        make(map[string]*wsSession[D]),
		eventRegistries: []neoroute.IEventRegistry{},
	}

	router.Init()

	hook := func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if rec := recover(); rec != nil {
				neoroute.PrintRecoveredPanic("WebSocket", rec)
			}
		}()

		// Perform handshake to get session data
		sessionData, ok := transporter.config.HandshakeFunc(r)
		if !ok {
			http.Error(w, router.Config().RunErrorHandler(neoroute.ErrHandshakeFailed, nil), http.StatusUnauthorized)
			return
		}

		// Upgrade to WebSocket session
		conn, err := websocket.Accept(w, r, config.AcceptOptions)
		if err != nil {
			neoroute.Logger.Info("Upgrade to WebSocket failed", "err", err)
			return
		}

		// The connection limit has to fit the biggest route, the route specific limit is checked by the router
		if max := router.MaxRequestSize(); max > 0 {
			conn.SetReadLimit(max)
		} else {
			conn.SetReadLimit(-1)
		}

		// Add session to transporter
		session := transporter.addSession(sessionData, conn)
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
func (t *Transporter[D]) SetRouter(r *neoroute.Router[D]) {
	t.router = r
}

// AddEventRegistry adds an event registry to the transporter.
// This allows the transporter to send events registered in the event registry.
func (t *Transporter[D]) AddEventRegistry(e *neoroute.EventRegistry) {
	t.mutex.Lock()
	t.eventRegistries = append(t.eventRegistries, e)
	t.mutex.Unlock()
}

func (t *Transporter[D]) addSession(sessionData D, conn *websocket.Conn) *wsSession[D] {

	// Check if session already exists and if it should be overwritten
	t.mutex.Lock()

	// Create session with unique id and provided session data
	var userSession *neoroute.Session[D]
	for {
		id := uuid.NewString()
		if _, exists := t.sessions[id]; !exists {
			userSession = neoroute.NewSession(id, sessionData, neoroute.SessionTransporterCallbacks[D]{
				Adapt:      t.adaptFunc(id),
				Disconnect: conn.CloseNow,
			})
			break
		}
	}

	ctx, cancel := context.WithCancel(context.Background())

	// Create new session entry
	session := &wsSession[D]{
		conn:      conn,
		session:   userSession,
		ctx:       ctx,
		sendMutex: &sync.Mutex{},
		cancel:    cancel,
	}
	t.sessions[userSession.Id()] = session
	t.mutex.Unlock()
	return session
}

func (t *Transporter[D]) removeSession(id string) {
	t.mutex.Lock()
	delete(t.sessions, id)
	t.mutex.Unlock()
}

func (t *Transporter[D]) adaptFunc(sessionId string) func() (neoroute.Adapter, error) {
	return func() (neoroute.Adapter, error) {
		wsSession, ok := t.getSession(sessionId)
		if !ok {
			return nil, fmt.Errorf("session %s not found", sessionId)
		}

		wsSession.mutex.Lock()
		conn := wsSession.conn
		sendMutex := wsSession.sendMutex
		ctx := wsSession.ctx
		wsSession.mutex.Unlock()

		if conn == nil {
			return nil, fmt.Errorf("websocket session not set for %s", sessionId)
		}

		adapter := &WebSocketAdapter{
			transporterType: "WebSocket",
			eventRegistries: t.eventRegistries,
			conn:            conn,
			sendMutex:       sendMutex,
			ctx:             ctx,
		}
		go adapter.waitClosed()
		return adapter, nil
	}
}

func (t *Transporter[D]) getSession(id string) (*wsSession[D], bool) {
	t.mutex.Lock()
	session, ok := t.sessions[id]
	t.mutex.Unlock()
	return session, ok
}

func (t *Transporter[D]) handleSession(session *wsSession[D]) {
	session.mutex.Lock()
	conn := session.conn
	session.mutex.Unlock()

	defer func() {
		if rec := recover(); rec != nil {
			neoroute.PrintRecoveredPanic("WebSocket", rec)
		}

		if session.cancel != nil {
			session.cancel()
		}
		conn.CloseNow()
		if handler := t.config.DisconnectHandler; handler != nil {
			handler(session.session)
		}
		session.mutex.Lock()
		t.removeSession(session.session.Id())
		session.mutex.Unlock()
	}()

	if handler := t.config.EnterNetworkFunc; handler != nil {
		handler(session.session)
	}

	for {
		resp, runAfter, err := t.readRequest(session)
		if err != nil {

			// Only log err if it is not due to expected connection closure
			if closeErr, ok := errors.AsType[websocket.CloseError](err); ok {
				neoroute.Logger.Info("websocket connection closed by remote",
					"code", closeErr.Code,
					"reason", closeErr.Reason,
				)
				return
			}

			return
		}

		if resp != nil {
			session.sendMutex.Lock()
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
			err = conn.Write(ctx, websocket.MessageBinary, resp)
			cancel()
			session.sendMutex.Unlock()
			if err != nil {
				neoroute.Logger.Info("failed to send websocket response", "err", err)
				return
			}
		}
		for _, fn := range runAfter {
			fn()
		}
	}
}

// readRequest scopes timeout cleanup to a single message, including partial messages.
func (t *Transporter[D]) readRequest(session *wsSession[D]) ([]byte, []func(), error) {
	ctx, cancel := context.WithCancel(session.ctx)
	defer cancel()
	var timer *time.Timer
	defer func() {
		if timer != nil {
			timer.Stop()
		}
	}()
	setTimeout := func(timeout time.Duration) error {
		if err := ctx.Err(); err != nil {
			return err
		}
		if timer != nil {
			timer.Stop()
		}
		if timeout > 0 {
			if timer == nil {
				timer = time.AfterFunc(timeout, cancel)
			} else {
				timer.Reset(timeout)
			}
		}
		return nil
	}

	// Wait without a deadline so idle connections and ping/pong traffic do not consume the time allowed to read an incoming request's envelope.
	messageType, reader, err := session.conn.Reader(ctx)
	if err != nil {
		return nil, nil, err
	}
	if messageType != websocket.MessageBinary {
		return nil, nil, fmt.Errorf("expected binary websocket message")
	}
	if err := setTimeout(t.config.HeaderReadTimeout); err != nil {
		return nil, nil, err
	}

	// Stop the read timeout as soon as the complete message has been read, before the router executes application handlers. The router checks for EOF itself.
	reader = &messageReader{Reader: reader, cancel: func() {
		if timer != nil {
			timer.Stop()
		}
		cancel()
	}, setTimeout: setTimeout}
	resp, runAfter, err := t.router.Handle(reader, session.session)
	if err != nil {
		return nil, nil, err
	}
	if _, err := io.Copy(io.Discard, reader); err != nil {
		return nil, nil, err
	}
	return resp, runAfter, nil
}

type messageReader struct {
	io.Reader
	cancel     context.CancelFunc
	done       bool
	setTimeout func(time.Duration) error
}

func (r *messageReader) SetReadTimeout(timeout time.Duration) error {
	return r.setTimeout(timeout)
}

func (r *messageReader) Read(p []byte) (int, error) {
	if r.done {
		return 0, io.EOF
	}
	n, err := r.Reader.Read(p)
	if err == io.EOF {
		r.done = true
		r.cancel()
	}
	return n, err
}
