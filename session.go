package neoroute

import "sync"

type Session[D any] struct {
	mutex       *sync.Mutex
	id          string
	sessionData D
}

func NewSession(id string) *Session[any] {
	return &Session[any]{
		id:          id,
		sessionData: nil,
	}
}

func (s *Session[D]) GetId() string {
	return s.id
}

func (s *Session[D]) GetData() D {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.sessionData
}

func (s *Session[D]) SetData(data D) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.sessionData = data
}
