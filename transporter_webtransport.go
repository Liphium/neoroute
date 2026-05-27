package neoroute

import (
	"io"
	"net/http"
	"sync"

	"github.com/quic-go/webtransport-go"
)

type WebTransportTransporter[D any] struct {
	router   *NeoRouter[D]
	config   WTTConfig[D]
	mutex    *sync.Mutex
	sessions map[string]*wttSession[D]
}

type UpgradeFunc func(w http.ResponseWriter, r *http.Request) (*webtransport.Session, error)

type WTTConfig[D any] struct {
	UpgradeFunc          UpgradeFunc
	OverwriteSessionFunc func(id string) bool

	HandshakeFunc     func(r *http.Request) (*Session[D], bool)
	EnterNetworkFunc  func(session *Session[D])
	DisconnectHandler func(session *Session[D])

	WantReliableSteam        bool
	WantUnreliableConnection bool
}

type wttSession[D any] struct {
	mutex     *sync.Mutex
	wtSession *webtransport.Session
	session   *Session[D]
}

var _ Transporter[any] = &WebTransportTransporter[any]{}

func NewWebTransportTransporter[D any](config WTTConfig[D]) (http.HandlerFunc, *WebTransportTransporter[D]) {
	transporter := &WebTransportTransporter[D]{
		router:   nil,
		config:   config,
		sessions: make(map[string]*wttSession[D]),
		mutex:    &sync.Mutex{},
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

		// Upgrade to WebTransport session
		transportSession, err := transporter.config.UpgradeFunc(w, r)
		if err != nil {
			logger.Info("Upgrade to WebTransport failed", "err", err)
			return
		}

		// Add session to transporter
		session := transporter.addSession(userSession.Id(), userSession, transportSession)
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
func (t *WebTransportTransporter[D]) SetRouter(r *NeoRouter[D]) {
	t.router = r
}

func (t *WebTransportTransporter[D]) addSession(id string, userSession *Session[D], transportSession *webtransport.Session) *wttSession[D] {

	// Check if session already exists and if it should be overwritten
	t.mutex.Lock()
	if oldSession, exists := t.sessions[userSession.id]; exists {

		if t.config.OverwriteSessionFunc(userSession.id) {

			oldSession.mutex.Lock()

			// Close existing session before overwriting
			oldSession.wtSession.CloseWithError(0, "New session established.")

			oldSession.mutex.Unlock()

		} else {
			t.mutex.Unlock()
			return nil
		}
	}

	// Create new session entry
	session := &wttSession[D]{
		mutex:     &sync.Mutex{},
		wtSession: transportSession,
		session:   userSession,
	}
	t.sessions[userSession.id] = session
	t.mutex.Unlock()
	return session
}

func (t *WebTransportTransporter[D]) removeSession(id string) {
	t.mutex.Lock()
	delete(t.sessions, id)
	t.mutex.Unlock()
}

func (t *WebTransportTransporter[D]) Adapt() {

}

func (t *WebTransportTransporter[D]) AdaptUnreliable() {

}

func (t *WebTransportTransporter[D]) handleSession(session *wttSession[D]) {

	defer func() {
		t.config.DisconnectHandler(session.session)
		session.mutex.Lock()
		t.removeSession(session.session.id)
		session.mutex.Unlock()
	}()

	var doneStream chan struct{}
	var doneDatagram chan struct{}

	t.config.EnterNetworkFunc(session.session)

	if t.config.WantReliableSteam {
		doneStream = make(chan struct{})
		go t.listenStream(session, doneStream, doneDatagram)
	}

	if t.config.WantUnreliableConnection {
		doneDatagram = make(chan struct{})
		go t.listenDatagram(session, doneDatagram, doneStream)
	}

	// Check if any connection method is used
	if doneStream == nil && doneDatagram == nil {
		return
	}

	// Wait for either stream or datagram connection to close
	select {
	case <-doneStream:
		return
	case <-doneDatagram:
		return
	}

}

func (t *WebTransportTransporter[D]) listenStream(session *wttSession[D], done chan struct{}, otherDone chan struct{}) {
	defer func() {
		session.mutex.Lock()
		if session.wtSession != nil {
			session.wtSession.CloseWithError(0, "Connection closed.")
		}
		session.mutex.Unlock()
		close(done)
	}()

	session.mutex.Lock()
	wtSession := session.wtSession
	userSession := session.session
	session.mutex.Unlock()

	for {

		// Check if datagram connection is closed
		select {
		case <-otherDone:
			return
		default:
		}

		st, err := wtSession.AcceptStream(wtSession.Context())
		if err != nil {
			logger.Info("failed to accept stream", "err", err)
			return // session closed / error
		}

		// Collect request data
		reqBytes, err := io.ReadAll(st)
		if err != nil {
			logger.Info("failed to read stream request", "err", err)
			return
		}

		// Handle request and send response
		resp := t.router.handle(reqBytes, userSession)
		if resp != nil {
			_, err = st.Write(resp)
			if err != nil {
				logger.Info("failed to send response", "err", err)
				st.Close()
				return
			}
		}
		st.Close()

	}

}

func (t *WebTransportTransporter[D]) listenDatagram(session *wttSession[D], done chan struct{}, otherDone chan struct{}) {
	defer func() {
		session.mutex.Lock()
		if session.wtSession != nil {
			session.wtSession.CloseWithError(0, "Closing session.")
		}
		session.mutex.Unlock()
		close(done)
	}()

	session.mutex.Lock()
	wtSession := session.wtSession
	userSession := session.session
	session.mutex.Unlock()

	for {

		// Check if datagram connection is closed
		select {
		case <-otherDone:
			return
		default:
		}

		// Receive request data
		data, err := wtSession.ReceiveDatagram(wtSession.Context())
		if err != nil {
			logger.Info("error receiving datagram", "err", err)
			return
		}

		// Handle request and send response
		resp := t.router.handle(data, userSession)
		if resp != nil {
			err = wtSession.SendDatagram(resp)
			if err != nil {
				logger.Info("error sending datagram", "err", err)
				return
			}
		}

	}
}
