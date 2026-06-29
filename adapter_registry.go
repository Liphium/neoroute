package neoroute

import (
	"errors"
	"fmt"
	"os"
	"sync"
)

type AdapterRegistry struct {
	mutex    *sync.RWMutex
	adapters map[string]Adapter
}

func NewAdapterRegistry() *AdapterRegistry {
	return &AdapterRegistry{
		mutex:    &sync.RWMutex{},
		adapters: make(map[string]Adapter),
	}
}

func (r *AdapterRegistry) Register(name string, adapter Adapter) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.adapters[name] = adapter
	adapter.SetRemoveFunc(func() {
		r.unregisterIfSame(name, adapter)
	})
}

func (r *AdapterRegistry) Unregister(name string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	_, exists := r.adapters[name]
	if !exists {
		return
	}
	delete(r.adapters, name)
}

func (r *AdapterRegistry) Disconnect(name string) {
	r.mutex.RLock()
	adapter, exists := r.adapters[name]
	r.mutex.RUnlock()
	if !exists {
		return
	}
	adapter.Disconnect()

}

func (r *AdapterRegistry) unregisterIfSame(name string, adapter Adapter) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if current, exists := r.adapters[name]; exists && current == adapter {
		delete(r.adapters, name)
	}
}

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
