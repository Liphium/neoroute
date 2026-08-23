package neoroute

import (
	"errors"
	"fmt"
	"os"
	"sync"
)

// AdapterRegistry stores a list of adapters, basically a one-way channel to send messages to a specific connection.
//
// By using adapter registries, you can store a list of connections, send events to them and disconnect all of them when you don't need them anymore.
//
// It allows you to shift the functionality of sending to different parts of your codebase without actually exposing the raw session to any of your service functions. In fact, that's the reason why Neoroute doesn't give you a way to send events directly to different sessions.
type AdapterRegistry struct {
	mutex    sync.RWMutex
	adapters map[string]Adapter
}

// NewAdapterRegistry creates a new AdapterRegistry.
func NewAdapterRegistry() *AdapterRegistry {
	return &AdapterRegistry{
		adapters: make(map[string]Adapter),
	}
}

// Register registers an adapter with the given name.
func (r *AdapterRegistry) Register(name string, adapter Adapter) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.adapters[name] = adapter
	adapter.SetRemoveFunc(func() {
		r.unregisterIfSame(name, adapter)
	})
}

// Unregister unregisters the adapter with the given name.
func (r *AdapterRegistry) Unregister(name string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	_, exists := r.adapters[name]
	if !exists {
		return
	}
	delete(r.adapters, name)
}

// Disconnect disconnects and unregisters the adapter (the underlying connection) with the given name.
func (r *AdapterRegistry) Disconnect(name string) {
	r.mutex.Lock()
	adapter, exists := r.adapters[name]
	delete(r.adapters, name)
	r.mutex.Unlock()
	if !exists {
		return
	}
	adapter.Disconnect()
}

// UnregisterAll unregisters all adapters.
func (r *AdapterRegistry) UnregisterAll() {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.adapters = make(map[string]Adapter)
}

// DisconnectAll disconnects and unregisters all adapters (the underlying connections).
func (r *AdapterRegistry) DisconnectAll() {
	r.mutex.Lock()
	adapters := make([]Adapter, 0, len(r.adapters))
	for _, adapter := range r.adapters {
		adapters = append(adapters, adapter)
	}
	r.adapters = make(map[string]Adapter)
	r.mutex.Unlock()

	for _, adapter := range adapters {
		adapter.Disconnect()
	}
}

// GetAdapters returns a list of all adapter names.
func (r *AdapterRegistry) GetAdapters() []string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	names := make([]string, 0, len(r.adapters))
	for name := range r.adapters {
		names = append(names, name)
	}
	return names
}

// unregisterIfSame unregisters the adapter if it is the same as the one provided.
func (r *AdapterRegistry) unregisterIfSame(name string, adapter Adapter) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if current, exists := r.adapters[name]; exists && current == adapter {
		delete(r.adapters, name)
	}
}

// Send sends an event to the adapter with the given name.
func (r *AdapterRegistry) Send(name string, event event) error {
	r.mutex.RLock()
	adapter, exists := r.adapters[name]
	r.mutex.RUnlock()
	if !exists {
		return fmt.Errorf("adapter with name %s not found", name)
	}
	eventBytes := messageEvent(event)

	// Check the event is registered with transporter or exit
	if ok := adapter.IsEventRegistered(event.Name); !ok {
		Logger.Error("event is not registered with transporter", "transporter", adapter.GetTransportType(), "event", event.Name)
		os.Exit(1)
	}

	return adapter.Send(eventBytes)
}

// Broadcast sends an event to all registered adapters.
func (r *AdapterRegistry) Broadcast(event event) error {

	// Collect adapters to send to
	r.mutex.RLock()
	adapters := make([]Adapter, 0, len(r.adapters))
	for _, adapter := range r.adapters {
		adapters = append(adapters, adapter)
	}
	r.mutex.RUnlock()

	if len(adapters) == 0 {
		return nil
	}

	var wg sync.WaitGroup
	errCh := make(chan error, len(adapters))

	eventBytes := messageEvent(event)

	// Send event to all adapters concurrently
	for _, adapter := range adapters {
		wg.Add(1)
		go func(a Adapter) {
			defer wg.Done()

			// Check the event is registered with transporter or exit
			if ok := a.IsEventRegistered(event.Name); !ok {
				Logger.Error("event is not registered with transporter", "transporter", a.GetTransportType(), "event", event.Name)
				os.Exit(1)
			}

			if err := a.Send(eventBytes); err != nil {
				errCh <- err
			}
		}(adapter)
	}

	wg.Wait()
	close(errCh)

	// Collect errors
	var errs []error
	for err := range errCh {
		errs = append(errs, err)
	}
	if len(errs) == 0 {
		return nil
	}
	return errors.Join(errs...)
}
