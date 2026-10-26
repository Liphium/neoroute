package http_transporter

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Liphium/neoroute"
	"github.com/stretchr/testify/assert"
)

func TestSecurity(t *testing.T) {

	// Create one router + http transporter + http server for this setup
	var handled atomic.Bool
	handler := func(*neoroute.Ctx[bool]) { handled.Store(true) }
	router := neoroute.NewRouter(neoroute.Config[bool]{})
	router.RoutePing("x", handler, neoroute.WithReadTimeout(50*time.Millisecond))
	router.Group("group").RoutePing("x", handler, neoroute.WithReadTimeout(time.Second))
	hook, _ := NewTransporter(router, Config[bool]{
		HeaderReadTimeout: 100 * time.Millisecond,
		HandshakeFunc:     func(*http.Request) (bool, bool) { return false, true },
	})
	server := httptest.NewServer(hook)
	defer server.Close()

	for _, tc := range []struct {
		name          string
		payload       []byte
		contentLength int
		delay         time.Duration
		remaining     []byte
		wantStatus    int
		minDuration   time.Duration
		wantHandled   bool
	}{
		{
			name: "partial body times out",
			// Declare a 90-byte payload, then stall after its envelope.
			payload:       []byte{0x93, 1, 0xa1, 'x', 0xc4, 90},
			contentLength: 100,
			wantStatus:    http.StatusBadRequest,
			minDuration:   50 * time.Millisecond,
		},
		{
			name: "partial header times out",
			// Send only the array prefix of the Neoroute envelope.
			payload:       []byte{0x93},
			contentLength: 100,
			wantStatus:    http.StatusBadRequest,
			minDuration:   100 * time.Millisecond,
		},
		{
			name:          "route body timeout replaces header timeout",
			payload:       []byte{0x93, 1, 0xa7, 'g', 'r', 'o', 'u', 'p', '/', 'x', 0xc4, 2},
			contentLength: 14,
			// Complete the body after the header deadline, within the route deadline.
			delay:       250 * time.Millisecond,
			remaining:   []byte{0, 0},
			wantStatus:  http.StatusOK,
			wantHandled: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			handled.Store(false)

			// Dial the test server
			conn, err := net.Dial("tcp", server.Listener.Addr().String())
			if !assert.NoError(t, err) {
				return
			}
			defer conn.Close()
			if !assert.NoError(t, conn.SetDeadline(time.Now().Add(3*time.Second))) {
				return
			}
			start := time.Now()

			// Actually send the payload
			_, err = fmt.Fprintf(conn, "POST / HTTP/1.1\r\nHost: example.com\r\nContent-Length: %d\r\n\r\n%s", tc.contentLength, tc.payload)
			if !assert.NoError(t, err) {
				return
			}
			if len(tc.remaining) > 0 {
				time.Sleep(tc.delay)
				_, err = conn.Write(tc.remaining)
				if !assert.NoError(t, err) {
					return
				}
			}

			// Read the response and make sure it's according to the things we got from the test
			resp, err := http.ReadResponse(bufio.NewReader(conn), nil)
			if !assert.NoError(t, err) {
				return
			}
			defer resp.Body.Close()
			assert.Equal(t, tc.wantStatus, resp.StatusCode)
			assert.GreaterOrEqual(t, time.Since(start), tc.minDuration, "request was rejected before the read timeout")
			assert.Equal(t, tc.wantHandled, handled.Load(), "unexpected handler call")
		})
	}
}
