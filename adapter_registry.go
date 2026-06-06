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
	adapter.setRemoveFunc(func() {
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
	adapter.disconnect()

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
	eventBytes, err := messageEvent(event)
	if err != nil {
		return fmt.Errorf("marshal event for adapter %s: %v", name, err)
	}

	// Check the event is registered with transporter or exit
	if ok := adapter.isEventRegistered(event.Name); !ok {
		logger.Error("event is not registered with transporter", "transporter", adapter.getTransportType(), "event", event.Name)
		os.Exit(1)
	}

	return adapter.send(eventBytes)
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

	eventBytes, err := messageEvent(event)
	if err != nil {
		return fmt.Errorf("marshal event for broadcast: %v", err)
	}

	// Send event to all adapters concurrently
	for _, adapter := range adapters {
		wg.Add(1)
		go func(a Adapter) {
			defer wg.Done()

			// Check the event is registered with transporter or exit
			if ok := a.isEventRegistered(event.Name); !ok {
				logger.Error("event is not registered with transporter", "transporter", a.getTransportType(), "event", event.Name)
				os.Exit(1)
			}

			if err := a.send(eventBytes); err != nil {
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
