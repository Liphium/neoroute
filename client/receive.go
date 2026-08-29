package client

import (
	"github.com/tinylib/msgp/msgp"
)

// Receive binds the event name to a handle function.
//
// If the server sends any event with that name to this client, the provided function will be called to handle it.
func (c *Client) Receive[E any, EP msgp.UnmarshalPtr[E]](eventName string, handleFunc func(c *Ctx, data E)) {
	c.setReceiver(eventName, func(c *Ctx) {

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
