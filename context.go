package neoroute

import (
	"fmt"

	"github.com/tinylib/msgp/msgp"
)

type Ctx[RS msgp.Marshaler, D any] struct {
	ctx[D]
}

func (c *Ctx[RS, D]) Respond(resp RS) error {

	// Marshal response data
	respData, err := resp.MarshalMsg(nil)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %v", err)
	}

	return response{
		Id:      c.id,
		IsError: false,
		Data:    respData,
	}
}

func (c *Ctx[RS, D]) RespondError(err error) error {
	return response{
		Id:      c.id,
		IsError: true,
		Data:    []byte(fmt.Sprintf("%v", err)),
	}
}

type ctx[D any] struct {
	neo     *NeoRouter[D]
	id      int    // request id, used for responses
	data    []byte // data field from Request struct
	route   string // the route that matched the request
	session *Session[D]
}

func (c *ctx[D]) Id() int {
	return c.id
}

func (c *ctx[D]) Data() []byte {
	return c.data
}

func (c *ctx[D]) Route() string {
	return c.route
}

func (c *ctx[D]) Session() *Session[D] {
	return c.session
}

func (c *ctx[D]) respondError(err error) response {
	return response{
		Id:      c.id,
		IsError: true,
		Data:    []byte(fmt.Sprintf("%v", err)),
	}
}
