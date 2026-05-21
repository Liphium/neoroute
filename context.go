package neoroute

import (
	"fmt"

	"github.com/tinylib/msgp/msgp"
)

type Ctx[RS msgp.Marshaler] struct {
	ctx
}

func (c *Ctx[RS]) Respond(resp RS) error {

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

func (c *Ctx[RS]) RespondError(err error) error {
	return response{
		Id:      c.id,
		IsError: true,
		Data:    []byte(fmt.Sprintf("%v", err)),
	}
}

type ctx struct {
	neo   *NeoRouter
	id    int    // request id, used for responses
	Data  []byte // data field from Request struct
	Route string // the route that matched the request
}

func (c *ctx) respondError(err error) response {
	return response{
		Id:      c.id,
		IsError: true,
		Data:    []byte(fmt.Sprintf("%v", err)),
	}
}
