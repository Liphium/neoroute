package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	client "github.com/Liphium/neoroute/client/sdk/v2"
	clienthttp "github.com/Liphium/neoroute/client/transporter/http/v2"
	clientws "github.com/Liphium/neoroute/client/transporter/websocket/v2"
	"github.com/Liphium/neoroute/v2"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/require"
)

func TestSend(t *testing.T) {
	t.Run("send message success", func(t *testing.T) {
		session := neoroute.NewTestingSession(neoroute.NoData{}, "connection1")

		ret := Send(session.NewTestingOkCtx("send"), SendRequest{
			Text:   "some text",
			Sender: "sender1",
		})
		neoroute.AssertResponseOk(t, ret)
	})

	t.Run("send empty message", func(t *testing.T) {
		session := neoroute.NewTestingSession(neoroute.NoData{}, "connection1")
		ret := Send(session.NewTestingOkCtx("send"), SendRequest{
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
		ctx := session1.NewTestingOkCtx("send")
		ret := Send(ctx, SendRequest{
			Text:   "some text",
			Sender: session1.Id(),
		})
		neoroute.AssertResponseOk(t, ret)

		// Evaluate RunAfter functions
		ctx.BaseCtx().EvaluateCtxTesting()

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

// THIS IS AN INTERNAL NEOROUTE TEST.
//
// Don't use this as a reference for how to test your applications, this is here to make sure neoroute's underlying encoding and stuff works. We use it to verify sending and receiving end-to-end across both HTTP and WebSocket.
func TestClientTransportRoundTrip(t *testing.T) {
	for _, transport := range []string{"http", "ws"} {
		t.Run(transport, func(t *testing.T) {

			// Start the chat server with both transporters
			router := neoroute.NewRouter(neoroute.Config[neoroute.NoData]{})
			router.RouteOk("send", Send)

			server := httptest.NewServer(InitTransporter(router))
			defer server.Close()

			// Connect the client using the selected transporter
			c := client.NewClient(client.Config{RequestTimeout: time.Second})
			u, err := url.Parse(server.URL + "/" + transport)
			require.NoError(t, err)

			events := make(chan MessageEvent, 1)
			if transport == "http" {
				clienthttp.ApplyHTTP(c, http.MethodPost, u)
			} else {
				u.Scheme = "ws"
				c.Receive("message", func(_ *client.Ctx, event MessageEvent) {
					events <- event
				})

				ws := clientws.NewWebSocketTransporter(c)
				_, err := ws.Connect(u)
				require.NoError(t, err)
				defer ws.Close()
			}

			// Verify success and error responses reach the client
			err = c.SendOk("send", &SendRequest{Text: "hello", Sender: "test"})
			require.NoError(t, err)

			err = c.SendOk("send", &SendRequest{Sender: "test"})
			require.EqualError(t, err, "text is required")

			// Verify the WebSocket client also receives the broadcast
			if transport == "ws" {
				select {
				case event := <-events:
					require.Equal(t, "hello", event.Text)
					require.Equal(t, "test", event.Sender)
					require.False(t, event.Timestamp.IsZero())
				case <-time.After(time.Second):
					require.FailNow(t, "timed out waiting for broadcast")
				}
			}
		})
	}
}
