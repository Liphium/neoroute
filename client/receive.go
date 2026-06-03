package client

import (
	"fmt"

	"github.com/tinylib/msgp/msgp"
)

func Receive[E any, EP interface {
	*E
	msgp.Unmarshaler
}](r *Receiver, eventName string, handlerFunc func(c *Ctx, req E) error) {
	r.handler[eventName] = func(c *Ctx) error {

		// Parse request data into struct
		var data E
		unmarshaler := any(&data).(msgp.Unmarshaler)

		_, err := unmarshaler.UnmarshalMsg(c.data)
		if err != nil {
			return fmt.Errorf("failed to unmarshal struct: %v", err)
		}

		// Let the handler handle it
		return handlerFunc(c, data)
	}
}
