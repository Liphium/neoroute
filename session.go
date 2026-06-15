package neoroute

import (
	"sync"
)

type Session[D any] struct {
	mutex       *sync.Mutex
	id          string
	sessionData D
}

// NewSession creates a new session with the given id and returns a pointer to it.
// ONLY USE THIS FUNCTION IF YOU ARE IMPLEMENTING A TRANSPORTER.
func NewSession[D any](id string, data D) *Session[D] {
	return &Session[D]{
		mutex:       &sync.Mutex{},
		sessionData: data,
		id:          id,
	}
}

// Id returns the session's id.
// ONLY USE THIS FUNCTION IF YOU ARE IMPLEMENTING A TRANSPORTER.
func (s *Session[D]) Id() string {
	return s.id
}

func (s *Session[D]) Data() D {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.sessionData
}

// SetData allows you to set the sessions data.
func (s *Session[D]) SetData(data D) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.sessionData = data
}

// UpdateData allows you to update the sessions data.
// While the updateFunc is running, the session will be locked,
// so other goroutines will not be able to access the session data until the updateFunc is done.
func (s *Session[D]) UpdateData(updateFunc func(data *D)) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	updateData := &s.sessionData
	updateFunc(updateData)
	s.sessionData = *updateData
}
