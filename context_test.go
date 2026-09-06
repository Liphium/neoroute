package neoroute

import (
	"strings"
	"testing"
)

//go:generate msgp

type MsgStruct struct {
	Val string `msg:"val"`
}

// TestPointerResponse tests that routes with pointer response types return an error
// instead of panicking when a nil response is given.
func TestPointerResponse(t *testing.T) {
	t.Run("valid pointer response", func(t *testing.T) {
		value := MsgStruct{Val: "hi"}
		route := func(c *ResCtx[NoData, *MsgStruct]) error {
			return c.Respond(&value)
		}

		session := NewTestingSession(NoData{}, "test")
		res := route(session.NewTestingResCtx[*MsgStruct](""))
		AssertResponse(t, res, value)
	})

	t.Run("nil pointer response", func(t *testing.T) {
		route := func(c *ResCtx[NoData, *MsgStruct]) error {
			return c.Respond(nil)
		}

		session := NewTestingSession(NoData{}, "test")
		res := route(session.NewTestingResCtx[*MsgStruct](""))

		// Should return an error instead of panicking, as nil is not encodable.
		if res == nil || !strings.Contains(res.Error(), "response must never be nil") {
			t.Errorf("expected nil response error, got %v", res)
		}
	})

	t.Run("value response unaffected", func(t *testing.T) {
		value := MsgStruct{Val: "ok"}
		route := func(c *ResCtx[NoData, MsgStruct]) error {
			return c.Respond(value)
		}

		session := NewTestingSession(NoData{}, "test")
		res := route(session.NewTestingResCtx[MsgStruct](""))
		AssertResponse(t, res, value)
	})
}
