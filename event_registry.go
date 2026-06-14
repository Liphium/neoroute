package neoroute

import (
	"fmt"
	"sync"

	"github.com/tinylib/msgp/msgp"
)

type EventRegistry struct {
	mutex            *sync.Mutex
	registeredEvents []string
}

func NewEventRegistry() *EventRegistry {
	return &EventRegistry{
		mutex:            &sync.Mutex{},
		registeredEvents: []string{},
	}
}

// GetEvents returns the registered events in the registry.
// ONLY USE THIS WHEN IMPLEMENTING AN ADAPTER.
func (er *EventRegistry) GetEvents() []string {
	er.mutex.Lock()
	defer er.mutex.Unlock()
	return er.registeredEvents
}

// Register returns a new event builder.
// First register all events only then start creating adapters.
func Register[E any, EM interface {
	*E
	msgp.Marshaler
}](e *EventRegistry, name string) func(ev E) (event, error) {

	e.mutex.Lock()
	e.registeredEvents = append(e.registeredEvents, name)
	e.mutex.Unlock()

	return func(eventData E) (event, error) {

		// Marshal event data
		marshaler := any(&eventData).(msgp.Marshaler)
		respData, err := marshaler.MarshalMsg(nil)
		if err != nil {
			return event{}, fmt.Errorf("failed to marshal event data for event %v: %v", name, err)
		}

		return event{
			Name: name,
			Data: respData,
		}, nil
	}
}
