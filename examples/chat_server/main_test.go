package main

import (
	"testing"

	"github.com/Liphium/neoroute"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestSend(t *testing.T) {
	t.Run("send message success", func(t *testing.T) {
		session := neoroute.NewTestingSession(neoroute.NoData{}, "connection1")

		ret := Send(neoroute.NewTestingOkCtx(session, "send"), SendRequest{
			Text:   "some text",
			Sender: "sender1",
		})
		neoroute.AssertResponseOk(t, ret)
	})

	t.Run("send empty message", func(t *testing.T) {
		session := neoroute.NewTestingSession(neoroute.NoData{}, "connection1")
		ret := Send(neoroute.NewTestingOkCtx(session, "send"), SendRequest{
			Text:   "",
			Sender: "sender1",
		})
		neoroute.AssertUserError(t, ret, "text is required")
	})

	t.Run("broadcast message to two connections", func(t *testing.T) {
		session1 := neoroute.NewTestingSession(neoroute.NoData{}, "connection1")
		session2 := neoroute.NewTestingSession(neoroute.NoData{}, "connection2")

		// Register two connections, this simulates the registration that happens in the EnterNetworkFunc of the ws transporter.
		adapter1 := neoroute.NewTestingAdapter(eventRegistry)
		adapter2 := neoroute.NewTestingAdapter(eventRegistry)
		adapterRegistry.Register(session1.Id(), adapter1)
		adapterRegistry.Register(session2.Id(), adapter2)

		// Send message from connection1
		ctx := neoroute.NewTestingOkCtx(session1, "send")
		ret := Send(ctx, SendRequest{
			Text:   "some text",
			Sender: session1.Id(),
		})
		neoroute.AssertResponseOk(t, ret)

		// Evaluate RunAfter functions
		neoroute.EvaluateCtxTesting(ctx.BaseCtx())

		// Verify message received by all connections
		for _, adapter := range []*neoroute.TestingAdapter{adapter1, adapter2} {
			events := neoroute.AssertEvents(t, adapter, 1)
			neoroute.AssertEvent(t, events, 0, MessageEvent{
				Sender: session1.Id(),
				Text:   "some text",
			}, cmpopts.IgnoreFields(MessageEvent{}, "Timestamp"))
		}

	})
}
