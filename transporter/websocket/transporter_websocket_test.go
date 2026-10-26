package websocket_transporter

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Liphium/neoroute"
	"github.com/coder/websocket"
	"github.com/stretchr/testify/assert"
)

func TestSecurity(t *testing.T) {

	// Prepare transporter and router for all of this stuff
	var handled atomic.Int32
	router := neoroute.NewRouter(neoroute.Config[bool]{})
	handler := func(c *neoroute.OkCtx[bool]) error {
		handled.Add(1)
		if c.Route() == "s" {
			time.Sleep(150 * time.Millisecond)
		}
		return c.RespondOk()
	}
	router.RouteOkNoRequest("x", handler, neoroute.WithReadTimeout(100*time.Millisecond))
	router.RouteOkNoRequest("s", handler, neoroute.WithReadTimeout(100*time.Millisecond))
	router.RouteOkNoRequest("e", handler, neoroute.WithReadTimeout(time.Second))
	router.RouteOkNoRequest("u", handler, neoroute.WithReadTimeout(-1))
	hook, _ := NewTransporter(router, Config[bool]{
		HeaderReadTimeout: 100 * time.Millisecond,
		HandshakeFunc:     func(*http.Request) (bool, bool) { return false, true },
	})
	server := httptest.NewServer(hook)
	defer server.Close()

	for _, tc := range []struct {
		name      string
		payload   []byte
		remaining []byte
		delay     time.Duration
		idle      time.Duration
		repeat    int
		stall     bool
		wantCalls int32
	}{
		{
			name:    "partial header times out",
			payload: []byte{0x93},
			repeat:  1,
			stall:   true,
		},
		{
			name:    "partial body times out",
			payload: []byte{0x93, 1, 0xa1, 'x', 0xc4, 90},
			repeat:  1,
			stall:   true,
		},
		{
			name:      "idle gaps and slow handlers allow repeated requests",
			payload:   []byte{0x93, 1, 0xa1, 's', 0xc4, 0},
			idle:      250 * time.Millisecond,
			repeat:    3,
			wantCalls: 3,
		},
		{
			name:      "route body timeout replaces header timeout",
			payload:   []byte{0x93, 1, 0xa1, 'e', 0xc4, 2},
			remaining: []byte{0, 0},
			delay:     250 * time.Millisecond,
			repeat:    1,
			wantCalls: 1,
		},
		{
			name:      "route body timeout can be disabled",
			payload:   []byte{0x93, 1, 0xa1, 'u', 0xc4, 2},
			remaining: []byte{0, 0},
			delay:     250 * time.Millisecond,
			repeat:    1,
			wantCalls: 1,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {

			// Create a raw connection we can write individual frames with
			conn, raw := dialSecurityConnection(t, server.URL)
			if conn == nil {
				return
			}
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			callsBefore := handled.Load()

			// Send the request frames and read the response
			for range tc.repeat {
				time.Sleep(tc.idle)

				// Write the frame
				start := time.Now()
				complete := !tc.stall && len(tc.remaining) == 0
				writeSecurityFrame(t, raw, 0x02, complete, tc.payload)

				// Send the rest of the frame if the test is simulating a slow client
				if len(tc.remaining) > 0 {
					time.Sleep(tc.delay)
					writeSecurityFrame(t, raw, 0x00, true, tc.remaining)
				}

				// Read from the connection to see if the server closed it or sent a response
				kind, _, err := conn.Read(ctx)
				if tc.stall {
					assert.Error(t, err, "server should close the stalled request")
					assert.NoError(t, ctx.Err(), "server should enforce its timeout before the client deadline")
					assert.GreaterOrEqual(t, time.Since(start), 100*time.Millisecond, "connection closed before the read timeout")
				} else {
					if assert.NoError(t, err) {
						assert.Equal(t, websocket.MessageBinary, kind)
					}
				}
			}
			assert.Equal(t, tc.wantCalls, handled.Load()-callsBefore, "unexpected handler calls")
		})
	}
}

// Keep the raw connection so tests can send small non-final frames without the WebSocket writer buffering an incomplete Neoroute header until more data arrives.
func dialSecurityConnection(t *testing.T, serverURL string) (*websocket.Conn, net.Conn) {
	t.Helper()
	dialed := make(chan net.Conn, 1)
	transport := &http.Transport{DialContext: func(ctx context.Context, network, address string) (net.Conn, error) {
		conn, err := (&net.Dialer{}).DialContext(ctx, network, address)
		if err == nil {
			dialed <- conn
		}
		return conn, err
	}}
	t.Cleanup(transport.CloseIdleConnections)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	conn, _, err := websocket.Dial(ctx, "ws"+strings.TrimPrefix(serverURL, "http"), &websocket.DialOptions{
		HTTPClient: &http.Client{Transport: transport},
	})
	if !assert.NoError(t, err) {
		return nil, nil
	}
	t.Cleanup(func() { conn.CloseNow() })
	return conn, <-dialed
}

// All test payloads fit in a short frame. A zero mask leaves the payload unchanged.
func writeSecurityFrame(t *testing.T, conn net.Conn, opcode byte, final bool, payload []byte) {
	t.Helper()
	if !assert.Less(t, len(payload), 126) {
		return
	}
	if final {
		opcode |= 0x80
	}
	frame := append([]byte{opcode, 0x80 | byte(len(payload)), 0, 0, 0, 0}, payload...)
	n, err := conn.Write(frame)
	assert.NoError(t, err)
	assert.Equal(t, len(frame), n)
}
