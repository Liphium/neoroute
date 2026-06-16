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
	gotDisconnected  bool
}

// NewTestingAdapter creates a new TestingAdapter with the given event registries.
// This adapter can be registered in tests instead of a real adapter to collect
// sent events and check if the correct events were sent.
func NewTestingAdapter(eventRegistries []*EventRegistry) Adapter {
	adapter := &TestingAdapter{
		transporterType: "TestingAdapter",
		eventRegistries: eventRegistries,
		mutex:           &sync.Mutex{},
		sendMutex:       &sync.Mutex{},
	}
	return adapter
}

// GetEvents returns the unmarshaled events that were sent to the adapter.
func (a *TestingAdapter) GetEvents() ([]event, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	return a.unmarshalEvents()
}

// DrainEvents returns the unmarshaled events that were sent to the adapter and clears the stored events.
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

// ConnectionStatus returns if the adapter is closed and if it got disconnected.
// This can be used in tests to check if the adapter was closed and or got disconnected.
func (a *TestingAdapter) ConnectionStatus() (closed bool, gotDisconnected bool) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	return a.closed, a.gotDisconnected
}

// unmarshalEvents unmarshals the stored messages into events. It returns an error
// if the messages could not be unmarshaled or if they are not of type event.
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

// Send stores the given bytes in the adapter.
// This implements the Send function of the Adapter interface
// and would normally Send the bytes to the client.
func (a *TestingAdapter) Send(b []byte) error {
	a.sendMutex.Lock()
	defer a.sendMutex.Unlock()
	a.receivedMessages = append(a.receivedMessages, b)
	return nil
}

// IsEventRegistered checks if the given event name is registered in any of the event registries of the adapter.
// This implements the IsEventRegistered function of the Adapter interface.
func (a *TestingAdapter) IsEventRegistered(name string) bool {
	for _, eventRegistry := range a.eventRegistries {
		if slices.Contains(eventRegistry.GetEvents(), name) {
			return true
		}
	}
	return false
}

// GetTransportType returns the type of the transporter.
// This implements the GetTransportType function of the Adapter interface.
func (a *TestingAdapter) GetTransportType() string {
	return a.transporterType
}

// SetRemoveFunc sets the function that should be called when the adapter is closed,
// in the testing adapter this will only happen when Close() is called.
// This implements the SetRemoveFunc function of the Adapter interface.
func (a *TestingAdapter) SetRemoveFunc(removeFunc func()) {
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

// Disconnect sets the gotDisconnected flag to true and calls Close() to close the adapter.
func (a *TestingAdapter) Disconnect() {
	a.mutex.Lock()
	a.gotDisconnected = true
	a.mutex.Unlock()
	a.Close()
}

// Close sets the closed flag to true and calls the remove function if it is set.
// Use this to simulate a user closing the connection to the server.
// This would normally be called by the transporter when the connection is closed.
func (a *TestingAdapter) Close() {
	a.mutex.Lock()
	a.closed = true
	removeFunc := a.removeFunc
	a.mutex.Unlock()

	if removeFunc != nil {
		a.removeOnce.Do(removeFunc)
	}
}

// UnmarshalEventTesting is a helper function to unmarshal event data from an event struct.
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
