package neoroute

import (
	"fmt"
	"slices"
	"sync"

	"github.com/tinylib/msgp/msgp"
)

type TestingAdapter struct {
	eventRegistries []*EventRegistry
	mutex           *sync.Mutex
	sendMutex       *sync.Mutex
	removeFunc      func()
	closed          bool
	transporterType string
	removeOnce      sync.Once

	// Stored events that are sent to the adapter. Used for testing purposes.
	receivedMessages [][]byte
}

func NewTestingAdapter(eventRegistries []*EventRegistry) Adapter {
	adapter := &TestingAdapter{
		transporterType: "TestingAdapter",
		eventRegistries: eventRegistries,
		mutex:           &sync.Mutex{},
		sendMutex:       &sync.Mutex{},
	}
	return adapter
}

func (a *TestingAdapter) GetEvents() ([]event, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	return a.unmarshalEvents()
}

func (a *TestingAdapter) DrainEvents() ([]event, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	events, err := a.unmarshalEvents()
	if err != nil {
		return nil, err
	}
	a.receivedMessages = nil
	return events, nil
}

func (a *TestingAdapter) unmarshalEvents() ([]event, error) {

	// Unmarshal events into struct
	events := []event{}
	for _, messageBytes := range a.receivedMessages {

		var msg message
		_, err := msg.UnmarshalMsg(messageBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal messages: %v", err)
		}

		if msg.Type != MessageTypeEvent {
			return nil, fmt.Errorf("message was not of type event: %v", msg.Type)
		}

		var event event
		_, err = event.UnmarshalMsg(msg.Data)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal events: %v", err)
		}
		events = append(events, event)
	}
	return events, nil
}

func (a *TestingAdapter) send(b []byte) error {
	a.sendMutex.Lock()
	defer a.sendMutex.Unlock()
	a.receivedMessages = append(a.receivedMessages, b)
	return nil
}

func (a *TestingAdapter) isEventRegistered(name string) bool {
	for _, eventRegistry := range a.eventRegistries {
		if slices.Contains(eventRegistry.getEvents(), name) {
			return true
		}
	}
	return false
}

func (a *TestingAdapter) getTransportType() string {
	return a.transporterType
}

func (a *TestingAdapter) setRemoveFunc(removeFunc func()) {
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

func (a *TestingAdapter) Close() {
	a.mutex.Lock()
	a.closed = true
	removeFunc := a.removeFunc
	a.mutex.Unlock()

	if removeFunc != nil {
		a.removeOnce.Do(removeFunc)
	}
}

func UnmarshalEventTesting[E any, EP interface {
	*E
	msgp.Unmarshaler
}](eventData []byte) (E, error) {
	var ev E
	unmarshaler := any(&ev).(msgp.Unmarshaler)
	_, err := unmarshaler.UnmarshalMsg(eventData)
	if err != nil {
		return ev, fmt.Errorf("failed to unmarshal event: %v", err)
	}
	return ev, nil
}
