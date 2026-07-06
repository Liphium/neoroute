// Code generated with neogen-generated v1 schema by neogen. DO NOT EDIT.
package main

import (
	"github.com/Liphium/neoroute/client"
	"github.com/Liphium/neoroute/client/transporter/websocket"
)

type MainConnector struct {
	*websocket.WebSocketTransporter
	receiver *client.Receiver
}

func NewMainConnector(config client.Config) *MainConnector {
	r := client.NewReceiver(config)

	return &MainConnector{
		WebSocketTransporter: websocket.NewWebSocketTransporter(r),
		receiver:             r,
	}
}

func (c *MainConnector) ReceiveNewPunSubmitted(handler func(event NewPunEvent)) {
	client.Receive[NewPunEvent, *NewPunEvent](c.receiver, "new_pun_submitted", func(c *client.Ctx, event NewPunEvent) {
		handler(event)
	})
}

func (c *MainConnector) SendSubmitPun(payload SubmitPunRequest) error {

	return client.SendOk(c.receiver, "submit_pun", payload)

}

func (c *MainConnector) SendEcho(payload EchoRequest) (EchoResponse, error) {

	return client.Send[EchoResponse](c.receiver, "echo", payload)

}
