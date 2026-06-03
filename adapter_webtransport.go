package neoroute

import (
	"context"
	"fmt"
	"slices"
	"sync"

	"github.com/quic-go/webtransport-go"
)

type WebTransportAdapter struct {
	session         *webtransport.Session
	mutex           *sync.Mutex
	transporterType string
	eventRegistries []*EventRegistry
	isUnreliable    bool
	removeFunc      func()
	closed          bool
	removeOnce      sync.Once
}

func (a *WebTransportAdapter) send(b []byte) error {
	if a.isUnreliable {
		return a.sendUnreliable(b)
	}
	return a.sendReliable(b)
}

func (a *WebTransportAdapter) isEventRegistered(name string) bool {
	for _, eventRegistry := range a.eventRegistries {
		if slices.Contains(eventRegistry.getEvents(), name) {
			return true
		}
	}
	return false
}

func (a *WebTransportAdapter) getTransportType() string {
	return a.transporterType
}

func (a *WebTransportAdapter) sendReliable(b []byte) error {
	if a.session == nil {
		return fmt.Errorf("webtransport session not set")
	}

	ctx := a.session.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	stream, err := a.session.OpenStreamSync(ctx)
	if err != nil {
		return fmt.Errorf("open webtransport stream: %w", err)
	}

	_, err = stream.Write(b)
	if err != nil {
		_ = stream.Close()
		return fmt.Errorf("write webtransport stream: %w", err)
	}

	if err := stream.Close(); err != nil {
		return fmt.Errorf("close webtransport stream: %w", err)
	}

	return nil
}

func (a *WebTransportAdapter) sendUnreliable(b []byte) error {
	if a.session == nil {
		return fmt.Errorf("webtransport session not set")
	}

	if err := a.session.SendDatagram(b); err != nil {
		return fmt.Errorf("send webtransport datagram: %w", err)
	}

	return nil
}

func (a *WebTransportAdapter) setRemoveFunc(removeFunc func()) {
	if removeFunc == nil {
		return
	}

	a.mutex.Lock()
	a.removeFunc = removeFunc
	closed := a.closed
	a.mutex.Unlock()

	if closed {
		a.removeOnce.Do(removeFunc)
	}
}

func (a *WebTransportAdapter) waitClosed() {
	if a.session == nil {
		return
	}

	<-a.session.Context().Done()

	a.mutex.Lock()
	a.closed = true
	removeFunc := a.removeFunc
	a.mutex.Unlock()

	if removeFunc != nil {
		a.removeOnce.Do(removeFunc)
	}
}
