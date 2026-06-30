// Code generated with neogen-generated v1 schema by neogen. DO NOT EDIT.
package main

import (
	"github.com/Liphium/neoroute/client"
	"github.com/Liphium/neoroute/client/transporter/websocket"
)

type PunsConnector struct {
	*websocket.WebSocketTransporter
	receiver *client.Receiver
}

func NewPunsConnector(config client.Config) *PunsConnector {
	r := client.NewReceiver(config)

	return &PunsConnector{
		WebSocketTransporter: websocket.NewWebSocketTransporter(r),
		receiver:             r,
	}
}

func (c *PunsConnector) ReceiveNewPunSubmitted(handler func(event NewPunEvent)) {
	client.Receive[NewPunEvent, *NewPunEvent](c.receiver, "new_pun_submitted", func(c *client.Ctx, event NewPunEvent) {
		handler(event)
	})
}

func (c *PunsConnector) SendEcho(payload EchoRequest) /* (EchoResponse, error) */ {
	// TODO: Implement
}

func (c *PunsConnector) SendSubmitPun(payload SubmitPunRequest) /* error */ {
	// TODO: Implement
}
