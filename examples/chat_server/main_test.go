package main

import (
	"testing"

	"github.com/Liphium/neoroute"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestSend(t *testing.T) {
	t.Run("send message success", func(t *testing.T) {

		session := neoroute.NewTestingSession(neoroute.NoData{}, "connection1")
		ctx := neoroute.NewTestingOkCtx("send", session)

		ret := Send(ctx, SendRequest{
			Text:   "some text",
			Sender: "sender1",
		})
		neoroute.AssertResponseOk(t, ret)
	})

	t.Run("send empty message", func(t *testing.T) {

		session := neoroute.NewTestingSession(neoroute.NoData{}, "connection1")
		ctx := neoroute.NewTestingOkCtx("send", session)

		ret := Send(ctx, SendRequest{
			Text:   "",
			Sender: "sender1",
		})
		neoroute.AssertUserError(t, ret, "text is required")
	})

	t.Run("broadcast message to two connections", func(t *testing.T) {

		connectionId1 := "connection1"
		connectionId2 := "connection2"

		// Register two connections, this simulates the registration that happens in the EnterNetworkFunc of the ws transporter.
		adapter1 := neoroute.NewTestingAdapter(eventRegistry)
		adapter2 := neoroute.NewTestingAdapter(eventRegistry)
		adapterRegistry.Register(connectionId1, adapter1)
		adapterRegistry.Register(connectionId2, adapter2)

		session1 := neoroute.NewTestingSession(neoroute.NoData{}, connectionId1)
		ctx := neoroute.NewTestingOkCtx("send", session1)

		// Send message from connection1
		ret := Send(ctx, SendRequest{
			Text:   "some text",
			Sender: connectionId1,
		})
		neoroute.AssertResponseOk(t, ret)

		// Verify message received by all connections
		for _, adapter := range []*neoroute.TestingAdapter{adapter1, adapter2} {
			events := neoroute.AssertEvents(t, adapter, 1)
			neoroute.AssertEvent(t, events, 0, MessageEvent{
				Sender: connectionId1,
				Text:   "some text",
			}, cmpopts.IgnoreFields(MessageEvent{}, "Timestamp"))
		}

	})
}
