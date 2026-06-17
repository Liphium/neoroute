package neoroute

import (
	"fmt"

	"github.com/tinylib/msgp/msgp"
)

// Context allows helper functions to accept Ctx, ResCtx, or
// OkCtx interchangeably to extract their underlying data.
type Context[D any] interface {
	BaseCtx() *Ctx[D] // BaseCtx returns the underlying Ctx, allowing access to the session, request data, etc.
}

// --------------------------------------------------------------------------------
// Base Context
// --------------------------------------------------------------------------------

type Ctx[D any] struct {
	neo      *NeoRouter[D]
	id       int         // request id, used for responses
	reqData  []byte      // data field from Request struct
	route    string      // the route that matched the request
	session  *Session[D] // clients session, contains the session data and id
	runAfter []func()    // functions to run after the handler finishes, used for cleanup
}

func (c *Ctx[D]) BaseCtx() *Ctx[D] {
	return c
}

func (c *Ctx[D]) Id() int {
	return c.id
}

func (c *Ctx[D]) Route() string {
	return c.route
}

func (c *Ctx[D]) Session() *Session[D] {
	return c.session
}

// RunAfter allows handlers to register functions that will be executed after the response is sent.
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

type ResCtx[D any, RS any, PS interface {
	*RS
	msgp.Marshaler
}] struct {
	*Ctx[D]
}

func (c *ResCtx[D, RS, PS]) BaseCtx() *Ctx[D] {
	return c.Ctx
}

// Respond sends a successful response with the provided data.
func (c *ResCtx[D, RS, PS]) Respond(resp RS) error {
	marshaler := any(&resp).(msgp.Marshaler)
	respData, err := marshaler.MarshalMsg(nil)
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

type OkCtx[D any] struct {
	*Ctx[D]
}

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
