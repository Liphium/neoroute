package neoroute

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/tinylib/msgp/msgp"
)

type IEventRegistry interface {
	GetEvents() []string
	GetSchemas() []func() reflect.Type
}

var _ IEventRegistry = &EventRegistry{}

// EventRegistry stores a list of events making it possible to link events directly to different transporters.EventRegistry
//
// This is especially important for schema generation, but also for you to not accidentally sent wrong events through wrong transporters.
type EventRegistry struct {
	mutex             sync.Mutex
	registeredEvents  []string
	registeredSchemas []func() reflect.Type
}

// NewEventRegistry returns a new EventRegistry instance.
func NewEventRegistry() *EventRegistry {
	return &EventRegistry{
		mutex:             sync.Mutex{},
		registeredEvents:  []string{},
		registeredSchemas: []func() reflect.Type{},
	}
}

// GetEvents returns the registered events in the registry.
//
// ONLY USE THIS WHEN IMPLEMENTING AN ADAPTER.
func (er *EventRegistry) GetEvents() []string {
	er.mutex.Lock()
	defer er.mutex.Unlock()
	return er.registeredEvents
}

// GetSchemas returns the registered schemas for the registered events in the registry (same index as event names).
//
// ONLY USE THIS WHEN IMPLEMENTING AN ADAPTER.
func (er *EventRegistry) GetSchemas() []func() reflect.Type {
	er.mutex.Lock()
	defer er.mutex.Unlock()
	return er.registeredSchemas
}

// Register returns a new event builder.
func (e *EventRegistry) Register[E msgp.Marshaler](name string) func(ev E) event {

	e.mutex.Lock()
	e.registeredEvents = append(e.registeredEvents, name)
	e.registeredSchemas = append(e.registeredSchemas, func() reflect.Type {
		return reflect.TypeFor[E]()
	})
	e.mutex.Unlock()

	return func(eventData E) event {

		// Marshal event data
		respData, err := eventData.MarshalMsg(nil)
		if err != nil {
			panic(fmt.Sprintf("failed to marshal event data for event %v: %v", name, err))
		}

		return event{
			Name: name,
			Data: respData,
		}
	}
}
