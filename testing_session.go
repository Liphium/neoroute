package neoroute

import "github.com/tinylib/msgp/msgp"

// NewTestingSession creates a new session for testing with the given data and session id.
func NewTestingSession[D any](data D, sessionId string) *Session[D] {
	return &Session[D]{
		sessionData: data,
		id:          sessionId,
	}
}

// ResCtx creates a new ResCtx with the provided route and session.
func (s *Session[D]) ResCtx[RS msgp.Marshaler](route string) *ResCtx[D, RS] {
	return &ResCtx[D, RS]{
		Ctx: s.Ctx(route),
	}
}

// OkCtx creates a new OkCtx for testing with the provided route and session.
func (s *Session[D]) OkCtx(route string) *OkCtx[D] {
	return &OkCtx[D]{
		Ctx: s.Ctx(route),
	}
}

// Ctx creates a new context for testing with the provided route and session.
func (s *Session[D]) Ctx(route string) *Ctx[D] {
	return &Ctx[D]{
		id:      1,
		route:   route,
		session: s,
	}
}

// EvaluateCtxTesting runs all the functions that were added to the context with RunAfter.
func EvaluateCtxTesting[D any](c *Ctx[D]) {
	for _, fn := range c.runAfter {
		fn()
	}
}
