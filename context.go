package neoroute

import (
	"fmt"

	"github.com/tinylib/msgp/msgp"
)

// Wraps context for handlers that return a response.

type ResCtx[RS any, PS interface {
	*RS
	msgp.Marshaler
}, D any] struct {
	Ctx[D]
}

func (c *ResCtx[RS, PS, D]) Respond(resp RS) error {

	// Marshal response data
	marshaler := any(&resp).(msgp.Marshaler)
	respData, err := marshaler.MarshalMsg(nil)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %v", err)
	}

	return response{
		Id:      c.id,
		HasData: true,
		IsError: false,
		Data:    respData,
	}
}

func (c *ResCtx[RS, PS, D]) RespondError(err string) error {
	return c.respondError(err)
}

// Wraps context for handlers that don't return any data, only success or error.

type OkCtx[D any] struct {
	Ctx[D]
}

func (c *OkCtx[D]) RespondOk() error {
	return response{
		Id:      c.id,
		HasData: false,
		IsError: false,
		Data:    []byte{},
	}
}

func (c *OkCtx[D]) RespondError(err string) error {
	return c.respondError(err)
}

// Default context that contains the request information.

type Ctx[D any] struct {
	neo     *NeoRouter[D]
	id      int    // request id, used for responses
	data    []byte // data field from Request struct
	route   string // the route that matched the request
	session *Session[D]
}

func (c *Ctx[D]) Id() int {
	return c.id
}

func (c *Ctx[D]) Data() []byte {
	return c.data
}

func (c *Ctx[D]) Route() string {
	return c.route
}

func (c *Ctx[D]) Session() *Session[D] {
	return c.session
}

func (c *Ctx[D]) respondError(err string) response {
	return response{
		Id:      c.id,
		HasData: true,
		IsError: true,
		Data:    []byte(err),
	}
}
