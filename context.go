package neoroute

import (
	"context"
	"fmt"
	"sync"

	"github.com/tinylib/msgp/msgp"
)

// Context allows helper functions to accept Ctx, ResCtx, or OkCtx interchangeably to extract their underlying data.
type Context[D any] interface {
	BaseCtx() *Ctx[D] // BaseCtx returns the underlying Ctx, allowing access to the session, request data, etc.
}

// --------------------------------------------------------------------------------
// Base Context
// --------------------------------------------------------------------------------

type Ctx[D any] struct {
	mu       sync.Mutex      // mu protects access to the ctx field
	id       int             // request id, used for responses
	reqData  []byte          // data field from Request struct
	route    string          // the route that matched the request
	session  *Session[D]     // clients session, contains the session data and id
	runAfter []func()        // functions to run after the handler finishes, used for cleanup
	ctx      context.Context // context is context.Background() by default
}

// BaseCtx returns the underlying base Ctx.
func (c *Ctx[D]) BaseCtx() *Ctx[D] {
	return c
}

// Context returns [context.Context] for use in middlewares and routes.
//
// By default is returns a background context, but is can be set for this specific request using SetContext.
//
// Context is concurrent-safe.
func (c *Ctx[D]) Context() context.Context {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.ctx
}

// SetContext sets the [context.Context] for this request.
//
// SetContext is concurrent-safe.
func (c *Ctx[D]) SetContext(ctx context.Context) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.ctx = ctx
}

// Route returns the route that matched the request.
func (c *Ctx[D]) Route() string {
	return c.route
}

// Session returns the client's session.
func (c *Ctx[D]) Session() *Session[D] {
	return c.session
}

// RunAfter allows handlers to register functions that will be executed after the response is sent.
// This function can be called multiple times, the functions will be executed in the order they are registered.
func (c *Ctx[D]) RunAfter(fn func(), fns ...func()) *Ctx[D] {
	c.runAfter = append(c.runAfter, fn)
	if len(fns) > 0 {
		c.runAfter = append(c.runAfter, fns...)
	}
	return c
}

func (c *Ctx[D]) respondError(msg string) response {
	return response{
		Id: c.id,

		HasData: true,
		IsError: true,
		Data:    []byte(msg),
	}
}

// --------------------------------------------------------------------------------
// Response Context
// --------------------------------------------------------------------------------

// ResCtx is the context for a request with a response type.
type ResCtx[D any, RS msgp.Marshaler] struct {
	*Ctx[D]

	// nilResponseCheck checks if the response is a nil value (pointer or interface).
	// Set by the route registration if the response type can be nil.
	nilResponseCheck func(RS) bool
}

// BaseCtx returns the underlying base Ctx.
func (c *ResCtx[D, RS]) BaseCtx() *Ctx[D] {
	return c.Ctx
}

// Respond sends a successful response with the provided data.
func (c *ResCtx[D, RS]) Respond(resp RS) error {
	// Make sure the response is never nil, as nil is not encodable.
	if c.nilResponseCheck != nil && c.nilResponseCheck(resp) {
		return fmt.Errorf("response must never be nil, as nil is not encodable")
	}

	respData, err := resp.MarshalMsg(nil)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %v", err)
	}

	return &responseData{
		HasData: true,
		IsError: false,
		Data:    respData,
	}
}

// -----------------------------------------------------------------------------
// OK Context (Used by RouteOk / RouteOkNoRequest)
// -----------------------------------------------------------------------------

// OkCtx is the context for request without response type.
type OkCtx[D any] struct {
	*Ctx[D]
}

// BaseCtx returns the underlying base Ctx.
func (c *OkCtx[D]) BaseCtx() *Ctx[D] {
	return c.Ctx
}

// RespondOk sends a successful response.
func (c *OkCtx[D]) RespondOk() error {
	return &responseData{
		HasData: false,
		IsError: false,
		Data:    []byte{},
	}
}
