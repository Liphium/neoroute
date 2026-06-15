package client

import (
	"github.com/tinylib/msgp/msgp"
)

// Receive binds the event name to a handle function.
// If the server sends a event with that name to this client
// the provided function will handle it.
func Receive[E any, EP interface {
	*E
	msgp.Unmarshaler
}](r *Receiver, eventName string, handleFunc func(c *Ctx, data E)) {
	r.setEvent(eventName, func(c *Ctx) {

		// Parse request data into struct
		var data E
		unmarshaler := any(&data).(msgp.Unmarshaler)

		_, err := unmarshaler.UnmarshalMsg(c.data)
		if err != nil {
			Logger.Info("failed to unmarshal struct event", "err", err)
			return
		}

		// Let the handler handle it
		handleFunc(c, data)
	})
}
