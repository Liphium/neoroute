package neoroute

import (
	"context"
	"slices"
	"sync"
	"time"

	"github.com/coder/websocket"
)

type WebSocketAdapter struct {
	conn            *websocket.Conn
	mutex           *sync.Mutex
	sendMutex       *sync.Mutex
	ctx             context.Context
	transporterType string
	eventRegistries []*EventRegistry
	removeFunc      func()
	closed          bool
	removeOnce      sync.Once
}

func (a *WebSocketAdapter) send(b []byte) error {
	a.sendMutex.Lock()
	defer a.sendMutex.Unlock()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	return a.conn.Write(ctx, websocket.MessageBinary, b)
}

func (a *WebSocketAdapter) isEventRegistered(name string) bool {
	for _, eventRegistry := range a.eventRegistries {
		if slices.Contains(eventRegistry.getEvents(), name) {
			return true
		}
	}
	return false
}

func (a *WebSocketAdapter) getTransportType() string {
	return a.transporterType
}

func (a *WebSocketAdapter) setRemoveFunc(removeFunc func()) {
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

func (a *WebSocketAdapter) waitClosed() {
	if a.conn == nil {
		return
	}

	<-a.ctx.Done()

	a.mutex.Lock()
	a.closed = true
	removeFunc := a.removeFunc
	a.mutex.Unlock()

	if removeFunc != nil {
		a.removeOnce.Do(removeFunc)
	}
}
