package neoroute

import "sync"

// NewTestingSession creates a new session for testing with the given data and session id.
func NewTestingSession[D any](data D, sessionId string) *Session[D] {
	return &Session[D]{
		sessionData: data,
		id:          sessionId,
		mutex:       &sync.Mutex{},
	}
}
