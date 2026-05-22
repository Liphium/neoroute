package neoroute

import "sync"

type Session[D any] struct {
	mutex       *sync.Mutex
	id          string
	sessionData D
}

func NewSession[D any](id string) *Session[D] {
	return &Session[D]{
		mutex: &sync.Mutex{},
		id:    id,
	}
}

func (s *Session[D]) Id() string {
	return s.id
}

func (s *Session[D]) Data() D {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.sessionData
}

func (s *Session[D]) SetData(data D) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.sessionData = data
}
