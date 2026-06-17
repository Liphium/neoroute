package transporter

import (
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/Liphium/neoroute"
	"github.com/google/uuid"
	"github.com/quic-go/webtransport-go"
)

type WebTransportTransporter[D any] struct {
	eventRegistriesUnreliable []*neoroute.EventRegistry
	eventRegistriesReliable   []*neoroute.EventRegistry
	router                    *neoroute.NeoRouter[D]
	config                    WTTConfig[D]
	mutex                     *sync.Mutex
	sessions                  map[string]*wttSession[D]
}

type UpgradeFuncWTT func(w http.ResponseWriter, r *http.Request) (*webtransport.Session, error)

type WTTConfig[D any] struct {
	UpgradeFunc          UpgradeFuncWTT
	OverwriteSessionFunc func(id string) bool

	// If session is nil, a new session will be created with a unique id. The data can then be set in the EnterNetworkFunc.
	// If the bool is false, the handshake will be considered failed and the connection will be rejected.
	HandshakeFunc     neoroute.HandshakeFunc[D]
	EnterNetworkFunc  func(session *neoroute.Session[D])
	DisconnectHandler func(session *neoroute.Session[D])

	WantReliableSteam        bool
	WantUnreliableConnection bool
}

type wttSession[D any] struct {
	mutex     *sync.Mutex
	wtSession *webtransport.Session
	session   *neoroute.Session[D]
}

func NewWebTransportTransporter[D any](router *neoroute.NeoRouter[D], config WTTConfig[D]) (http.HandlerFunc, *WebTransportTransporter[D]) {
	transporter := &WebTransportTransporter[D]{
		router:                    router,
		config:                    config,
		sessions:                  make(map[string]*wttSession[D]),
		mutex:                     &sync.Mutex{},
		eventRegistriesReliable:   []*neoroute.EventRegistry{},
		eventRegistriesUnreliable: []*neoroute.EventRegistry{},
	}

	hook := func(w http.ResponseWriter, r *http.Request) {

		// Perform handshake to get session data
		sessionData, ok := transporter.config.HandshakeFunc(r)
		if !ok {
			http.Error(w, "Handshake failed.", http.StatusUnauthorized)
			return
		}

		// Upgrade to WebTransport session
		transportSession, err := transporter.config.UpgradeFunc(w, r)
		if err != nil {
			neoroute.Logger.Info("Upgrade to WebTransport failed", "err", err)
			return
		}

		// Add session to transporter
		session := transporter.addSession(sessionData, transportSession)
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
func (t *WebTransportTransporter[D]) SetRouter(r *neoroute.NeoRouter[D]) {
	t.router = r
}

func (t *WebTransportTransporter[D]) AddEventRegistryReliable(e *neoroute.EventRegistry) {
	t.mutex.Lock()
	t.eventRegistriesReliable = append(t.eventRegistriesReliable, e)
	t.mutex.Unlock()
}

func (t *WebTransportTransporter[D]) AddEventRegistryUnreliable(e *neoroute.EventRegistry) {
	t.mutex.Lock()
	t.eventRegistriesUnreliable = append(t.eventRegistriesUnreliable, e)
	t.mutex.Unlock()
}

func (t *WebTransportTransporter[D]) addSession(sessionData D, transportSession *webtransport.Session) *wttSession[D] {

	// Check if session already exists and if it should be overwritten
	t.mutex.Lock()

	// Create session with unique id if handshake did not return one
	var userSession *neoroute.Session[D]
	for {
		id := uuid.NewString()
		if _, exists := t.sessions[id]; !exists {
			userSession = neoroute.NewSession(id, sessionData)
			break
		}
	}

	// Create new session entry
	session := &wttSession[D]{
		mutex:     &sync.Mutex{},
		wtSession: transportSession,
		session:   userSession,
	}
	t.sessions[userSession.Id()] = session
	t.mutex.Unlock()
	return session
}

func (t *WebTransportTransporter[D]) removeSession(id string) {
	t.mutex.Lock()
	delete(t.sessions, id)
	t.mutex.Unlock()
}

func (t *WebTransportTransporter[D]) Adapt(session *neoroute.Session[D]) (neoroute.Adapter, error) {
	return t.newAdapter(session.Id(), false)
}

func (t *WebTransportTransporter[D]) AdaptUnreliable(session *neoroute.Session[D]) (neoroute.Adapter, error) {
	return t.newAdapter(session.Id(), true)
}

func (t *WebTransportTransporter[D]) newAdapter(id string, unreliable bool) (neoroute.Adapter, error) {
	session, ok := t.getSession(id)
	if !ok {
		return nil, fmt.Errorf("session %s not found", id)
	}

	session.mutex.Lock()
	wtSession := session.wtSession
	session.mutex.Unlock()

	if wtSession == nil {
		return nil, fmt.Errorf("webtransport session not set for %s", id)
	}

	transporterType := "WebTransport Reliable"
	if unreliable {
		transporterType = "WebTransport Unreliable"
	}

	eventRegistries := t.eventRegistriesReliable
	if unreliable {
		eventRegistries = t.eventRegistriesUnreliable
	}

	adapter := &WebTransportAdapter{
		transporterType: transporterType,
		eventRegistries: eventRegistries,
		session:         wtSession,
		mutex:           &sync.Mutex{},
		isUnreliable:    unreliable,
	}
	go adapter.waitClosed()
	return adapter, nil
}

func (t *WebTransportTransporter[D]) getSession(id string) (*wttSession[D], bool) {
	t.mutex.Lock()
	session, ok := t.sessions[id]
	t.mutex.Unlock()
	return session, ok
}

func (t *WebTransportTransporter[D]) handleSession(session *wttSession[D]) {

	defer func() {
		t.config.DisconnectHandler(session.session)
		session.mutex.Lock()
		t.removeSession(session.session.Id())
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
			neoroute.Logger.Info("failed to accept stream", "err", err)
			return // session closed / error
		}

		// Collect request data
		reqBytes, err := io.ReadAll(st)
		if err != nil {
			neoroute.Logger.Info("failed to read stream request", "err", err)
			return
		}

		// Handle request and send response
		resp, runAfter := t.router.Handle(reqBytes, userSession)
		defer func() {
			for _, fn := range runAfter {
				fn()
			}
		}()
		if resp != nil {
			_, err = st.Write(resp)
			if err != nil {
				neoroute.Logger.Info("failed to send response", "err", err)
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
			neoroute.Logger.Info("error receiving datagram", "err", err)
			return
		}

		// Handle request and send response
		resp, runAfter := t.router.Handle(data, userSession)
		defer func() {
			for _, fn := range runAfter {
				fn()
			}
		}()
		if resp != nil {
			err = wtSession.SendDatagram(resp)
			if err != nil {
				neoroute.Logger.Info("error sending datagram", "err", err)
				return
			}
		}

	}
}
