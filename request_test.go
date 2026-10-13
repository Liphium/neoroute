package neoroute

import (
	"bytes"
	"strings"
	"testing"

	"github.com/tinylib/msgp/msgp"
)

// captureRoute registers a ping route that stores the request body it receives.
func captureRoute(r *Router[bool], route string, got *[]byte, called *bool, opts ...func(*RouteOptions)) {
	r.RoutePing(route, func(c *Ctx[bool]) {
		*called = true
		*got = append([]byte(nil), c.reqData...)
	}, opts...)
}

func TestRequestBodySizes(t *testing.T) {
	// Sizes around the header buffer boundaries so the body spans leftover and the raw reader.
	sizes := []int{0, 1, 3, 100, 900, 1000, 1023, 1024, 1025, 4096, 70000}

	for _, size := range sizes {
		t.Run("body", func(t *testing.T) {
			body := bytes.Repeat([]byte{0xAB}, size)

			var got []byte
			called := false
			r := NewRouter(Config[bool]{})
			captureRoute(r, "body-route", &got, &called)

			r.Init()
			_, runAfters, err := makeRequest(r, NewTestingSession(false, "s"), "body-route", body)
			if err != nil {
				t.Fatalf("size %d: makeRequest: %v", size, err)
			}
			executeRunAfters(runAfters)

			if !called {
				t.Fatalf("size %d: handler not called", size)
			}
			if !bytes.Equal(got, body) {
				t.Fatalf("size %d: body mismatch: got %d bytes, want %d", size, len(got), len(body))
			}
		})
	}
}

func TestRequestMaxSize(t *testing.T) {
	var got []byte
	called := false
	r := NewRouter(Config[bool]{})
	captureRoute(r, "limited", &got, &called, WithMaxRequestSize(10))

	// Exactly at the limit is accepted.
	r.Init()
	if _, _, err := makeRequest(r, NewTestingSession(false, "s"), "limited", make([]byte, 10)); err != nil {
		t.Fatal(err)
	}
	if !called || len(got) != 10 {
		t.Fatalf("body at limit not handled: called=%v len=%d", called, len(got))
	}

	// One over the limit is rejected before the body is read.
	called = false
	respBytes, _, err := makeRequest(r, NewTestingSession(false, "s"), "limited", make([]byte, 11))
	if err != nil {
		t.Fatal(err)
	}
	if called {
		t.Fatal("handler called for oversized request")
	}
	resp, err := handleResponseData(respBytes)
	if err != nil {
		t.Fatal(err)
	}
	if !resp.IsError {
		t.Fatal("expected error response for oversized request")
	}
}

func TestRequestHeaderTooLong(t *testing.T) {
	var got []byte
	called := false
	r := NewRouter(Config[bool]{})
	captureRoute(r, strings.Repeat("a", 2000), &got, &called)

	respBytes, _, err := makeRequest(r, NewTestingSession(false, "s"), strings.Repeat("a", 2000), []byte("hi"))
	if err != nil {
		t.Fatal(err)
	}
	if called {
		t.Fatal("handler called for request with oversized header")
	}
	resp, err := handleResponseData(respBytes)
	if err != nil {
		t.Fatal(err)
	}
	if !resp.IsError {
		t.Fatal("expected error response for oversized header")
	}
}

func TestRequestMalformedHeader(t *testing.T) {
	r := NewRouter(Config[bool]{})
	var got []byte
	called := false
	captureRoute(r, "x", &got, &called)
	r.Init()

	// Not a 3-tuple: array of 2 elements.
	bad := []byte{0x92, 0x01, 0xa1, 'x'}
	respBytes, _ := r.Handle(bytes.NewReader(bad), NewTestingSession(false, "s"))
	if called {
		t.Fatal("handler called for malformed header")
	}
	resp, err := handleResponseData(respBytes)
	if err != nil {
		t.Fatal(err)
	}
	if !resp.IsError {
		t.Fatal("expected error response for malformed header")
	}
}

func TestRequestBodyTrailingBytesIsError(t *testing.T) {
	body := bytes.Repeat([]byte{0x11}, 50)
	mReq := request{Id: 1, Route: "trail", Data: body}
	reqBytes, err := mReq.MarshalMsg(nil)
	if err != nil {
		t.Fatal(err)
	}

	for _, extra := range []int{1, 67, 2000} {
		var got []byte
		called := false
		r := NewRouter(Config[bool]{})
		captureRoute(r, "trail", &got, &called)
		r.Init()

		stream := append(append([]byte{}, reqBytes...), bytes.Repeat([]byte{0x22}, extra)...)
		respBytes, _ := r.Handle(bytes.NewReader(stream), NewTestingSession(false, "s"))
		if called {
			t.Fatalf("extra=%d: handler called despite trailing bytes", extra)
		}
		resp, err := handleResponseData(respBytes)
		if err != nil {
			t.Fatal(err)
		}
		if !resp.IsError {
			t.Fatalf("extra=%d: expected error response for trailing bytes", extra)
		}
	}
}

// Ensure the header decoder agrees with msgp for a plain request encoding.
func TestReadRequestHeaderMatchesMsgp(t *testing.T) {
	mReq := request{Id: 7, Route: "some/route", Data: make([]byte, 300)}
	reqBytes, err := mReq.MarshalMsg(nil)
	if err != nil {
		t.Fatal(err)
	}
	dc := msgp.NewReaderSize(bytes.NewReader(reqBytes), 1024)
	id, route, size, leftover, err := readRequestHeader(dc)
	if err != nil {
		t.Fatal(err)
	}
	if id != 7 || route != "some/route" || size != 300 {
		t.Fatalf("header mismatch: id=%d route=%q size=%d", id, route, size)
	}
	if len(leftover) == 0 {
		t.Fatal("expected buffered leftover bytes after header")
	}
}
