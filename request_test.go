package neoroute

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tinylib/msgp/msgp"
)

type timeoutRequestReader struct {
	*bytes.Reader
	timeout time.Duration
}

func (r *timeoutRequestReader) SetReadTimeout(timeout time.Duration) error {
	r.timeout = timeout
	return nil
}

func TestRouteReadTimeout(t *testing.T) {
	for _, tc := range []struct {
		name string
		base time.Duration
		opts []func(*RouteOptions)
		want time.Duration
	}{
		{"default", 0, nil, time.Minute},
		{"router default inherited by group", 2 * time.Second, nil, 2 * time.Second},
		{"route override", 2 * time.Second, []func(*RouteOptions){WithReadTimeout(5 * time.Second)}, 5 * time.Second},
		{"disabled", 2 * time.Second, []func(*RouteOptions){WithReadTimeout(-1)}, -1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := NewRouter(Config[bool]{RouteOptions: RouteOptions{ReadTimeout: tc.base}})
			called := false
			r.Group("group").RoutePing("x", func(*Ctx[bool]) { called = true }, tc.opts...)
			r.Init()
			req := request{Id: 1, Route: "group/x"}
			data, err := req.MarshalMsg(nil)
			require.NoError(t, err)
			reader := &timeoutRequestReader{Reader: bytes.NewReader(data)}
			_, _, err = r.Handle(reader, NewTestingSession(false, "s"))
			require.NoError(t, err)
			require.True(t, called)
			require.Equal(t, tc.want, reader.timeout)
		})
	}
}

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
			require.NoError(t, err, "size %d: makeRequest", size)
			executeRunAfters(runAfters)

			require.True(t, called, "size %d: handler not called", size)
			require.Truef(t, bytes.Equal(got, body), "size %d: body mismatch: got %d bytes, want %d", size, len(got), len(body))
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
	_, _, err := makeRequest(r, NewTestingSession(false, "s"), "limited", make([]byte, 10))
	require.NoError(t, err)
	require.True(t, called && len(got) == 10, "body at limit not handled: called=%v len=%d", called, len(got))

	// One over the limit is rejected before the body is read.
	called = false
	respBytes, _, err := makeRequest(r, NewTestingSession(false, "s"), "limited", make([]byte, 11))
	assert.Error(t, err)
	assert.Nil(t, respBytes)
	assert.False(t, called, "handler called for oversized request")
}

func TestRequestHeaderTooLong(t *testing.T) {
	var got []byte
	called := false
	r := NewRouter(Config[bool]{})
	captureRoute(r, strings.Repeat("a", 2000), &got, &called)

	respBytes, _, err := makeRequest(r, NewTestingSession(false, "s"), strings.Repeat("a", 2000), []byte("hi"))
	require.Error(t, err)
	require.False(t, called, "handler called for request with oversized header")
	require.Nil(t, respBytes)
}

func TestRequestMalformedHeader(t *testing.T) {
	r := NewRouter(Config[bool]{})
	var got []byte
	called := false
	captureRoute(r, "x", &got, &called)
	r.Init()

	// Not a 3-tuple: array of 2 elements.
	bad := []byte{0x92, 0x01, 0xa1, 'x'}
	respBytes, _, err := r.Handle(bytes.NewReader(bad), NewTestingSession(false, "s"))
	require.False(t, called, "handler called for malformed header")
	require.Error(t, err, "handler called without error")
	require.Nil(t, respBytes)
}

func TestRequestHeaderRejectsOversizedDeclaredRoute(t *testing.T) {
	r := NewRouter(Config[bool]{
		ErrorHandler: func(err error, _ *Ctx[bool]) string { return err.Error() },
	})
	r.Init()

	// Tuple with id=1 and a str32 route declaring 16 MiB, without string contents.
	// Reject the declared length before attempting to read or allocate the string.
	bad := []byte{0x93, 0x01, 0xdb, 0x01, 0x00, 0x00, 0x00}
	resp, _, err := r.Handle(bytes.NewReader(bad), NewTestingSession(false, "s"))
	require.ErrorContains(t, err, msgp.ErrLimitExceeded.Error())
	require.Nil(t, resp)
}

func TestRequestBodyTrailingBytesIsError(t *testing.T) {
	body := bytes.Repeat([]byte{0x11}, 50)
	mReq := request{Id: 1, Route: "trail", Data: body}
	reqBytes, err := mReq.MarshalMsg(nil)
	require.NoError(t, err)

	for _, extra := range []int{1, 67, 2000} {
		var got []byte
		called := false
		r := NewRouter(Config[bool]{})
		captureRoute(r, "trail", &got, &called)
		r.Init()

		stream := append(append([]byte{}, reqBytes...), bytes.Repeat([]byte{0x22}, extra)...)
		respBytes, _, err := r.Handle(bytes.NewReader(stream), NewTestingSession(false, "s"))
		require.False(t, called, "extra=%d: handler called despite trailing bytes", extra)
		require.Error(t, err, "no error even though extra body")
		require.Nil(t, respBytes)
	}
}

// Ensure the header decoder agrees with msgp for a plain request encoding.
func TestReadRequestHeaderMatchesMsgp(t *testing.T) {
	mReq := request{Id: 7, Route: "some/route", Data: make([]byte, 300)}
	reqBytes, err := mReq.MarshalMsg(nil)
	require.NoError(t, err)
	dc := msgp.NewReaderSize(bytes.NewReader(reqBytes), 1024)
	id, route, size, leftover, err := readRequestHeader(dc)
	require.NoError(t, err)
	require.Equal(t, 7, id, "id mismatch")
	require.Equal(t, "some/route", route, "route mismatch")
	require.Equal(t, int64(300), size, "size mismatch")
	require.NotEmpty(t, leftover, "expected buffered leftover bytes after header")
}
