package main

import (
	"testing"

	"github.com/Liphium/neoroute"
)

func TestSend(t *testing.T) {
	t.Run("send message success", func(t *testing.T) {

		session := neoroute.NewTestingSession(neoroute.NoData{}, "connection1")
		ctx := neoroute.NewTestingOkCtx("send", session)

		ret := Send(ctx, SendRequest{
			Text:   "some text",
			Sender: "sender1",
		})
		if userErr, err := neoroute.GetTestingResponseOk(ret); err != nil {
			t.Errorf("Send failed: %v", err)
		} else if userErr != "" {
			t.Errorf("Send failed: %s", userErr)
		}
	})

	t.Run("send empty message", func(t *testing.T) {

		session := neoroute.NewTestingSession(neoroute.NoData{}, "connection1")
		ctx := neoroute.NewTestingOkCtx("send", session)

		ret := Send(ctx, SendRequest{
			Text:   "",
			Sender: "sender1",
		})
		if userErr, err := neoroute.GetTestingResponseOk(ret); err != nil {
			t.Errorf("Send failed: %v", err)
		} else if userErr != "text is required" {
			t.Errorf("Did not receive user error for empty text: %v", userErr)
		}
	})

	t.Run("broadcast message to two connections", func(t *testing.T) {

		connectionId1 := "connection1"
		connectionId2 := "connection2"

		// Register two connections
		adapter1 := neoroute.NewTestingAdapter([]*neoroute.EventRegistry{eventRegistry})
		adapter2 := neoroute.NewTestingAdapter([]*neoroute.EventRegistry{eventRegistry})
		adapterRegistry.Register(connectionId1, adapter1)
		adapterRegistry.Register(connectionId2, adapter2)

		session1 := neoroute.NewTestingSession(neoroute.NoData{}, connectionId1)
		ctx := neoroute.NewTestingOkCtx("send", session1)

		// Send message from connection1
		ret := Send(ctx, SendRequest{
			Text:   "some text",
			Sender: connectionId1,
		})
		if userErr, err := neoroute.GetTestingResponseOk(ret); err != nil {
			t.Errorf("Send failed: %v", err)
		} else if userErr != "" {
			t.Errorf("Send failed: %s", userErr)
		}

		// Verify message received by all connections
		if events, err := adapter1.GetEvents(); err != nil {
			t.Errorf("GetEvents failed: %v", err)
		} else if len(events) != 1 {
			t.Errorf("Expected 1 event, got %d", len(events))
		} else if events[0].Name != "message" {
			t.Errorf("Expected event name 'message', got '%s'", events[0].Name)
		} else {
			ev, err := neoroute.UnmarshalEventTesting[MessageEvent](events[0].Data)
			if err != nil {
				t.Errorf("UnmarshalEventTesting failed: %v", err)
			}
			if ev.Text != "some text" {
				t.Errorf("Expected text 'some text', got '%s'", ev.Text)
			}
			if ev.Sender != connectionId1 {
				t.Errorf("Expected sender '%s', got '%s'", connectionId1, ev.Sender)
			}
		}

		if events, err := adapter2.GetEvents(); err != nil {
			t.Errorf("GetEvents failed: %v", err)
		} else if len(events) != 1 {
			t.Errorf("Expected 1 event, got %d", len(events))
		} else if events[0].Name != "message" {
			t.Errorf("Expected event name 'message', got '%s'", events[0].Name)
		} else {
			ev, err := neoroute.UnmarshalEventTesting[MessageEvent](events[0].Data)
			if err != nil {
				t.Errorf("UnmarshalEventTesting failed: %v", err)
			}
			if ev.Text != "some text" {
				t.Errorf("Expected text 'some text', got '%s'", ev.Text)
			}
			if ev.Sender != connectionId1 {
				t.Errorf("Expected sender '%s', got '%s'", connectionId1, ev.Sender)
			}
		}
	})
}
