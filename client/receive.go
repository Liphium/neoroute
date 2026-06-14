package client

import (
	"github.com/tinylib/msgp/msgp"
)

func Receive[E any, EP interface {
	*E
	msgp.Unmarshaler
}](r *Receiver, eventName string, handlerFunc func(c *Ctx, data E)) {
	r.handler[eventName] = func(c *Ctx) {

		// Parse request data into struct
		var data E
		unmarshaler := any(&data).(msgp.Unmarshaler)

		_, err := unmarshaler.UnmarshalMsg(c.data)
		if err != nil {
			Logger.Info("failed to unmarshal struct event", "err", err)
			return
		}

		// Let the handler handle it
		handlerFunc(c, data)
	}
}
