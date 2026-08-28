package connector

import (
	"net/url"

	tea "charm.land/bubbletea/v2"
	"github.com/Liphium/neoroute/client"
	"github.com/Liphium/neoroute/client/transporter/websocket"
	"github.com/Liphium/neoroute/neoschema"
	"github.com/Liphium/neoroute/pkg/neodebug/config"
	"github.com/Liphium/neoroute/pkg/neodebug/model"
)

func connectWebsocket(schema neoschema.TransporterSchema) tea.Msg {
	msgChan := make(chan tea.Msg)

	// Create the actual transporter
	c := client.NewClient(client.Config{}) // We don't really need the error handler as we handle errors below
	transporter := websocket.NewWebSocketTransporter(c)

	// Connect to the transporter using the URL in the config
	url, err := url.Parse(config.Config.TransporterURL)
	if err != nil {
		return withClose(model.Error("Failed to parse URL: %v", err))
	}
	doneChan, err := transporter.Connect(url)
	if err != nil {
		return withClose(model.Error("%s", err.Error()))
	}
	go func() {
		<-doneChan
		msgChan <- ClosedMsg{}
	}()

	// Listen for all events and make it emit messages
	for event := range schema.Events {
		c.Receive(event, func(c *client.Ctx, ev PackedAny) {
			msgChan <- model.Event(event, ev.Value)
		})
	}

	return model.Multiple(
		ConnectedMsg{
			Connection: WebSocketConnection{
				transporter: transporter,
				client:      c,
				schema:      schema,
				msgChan:     msgChan,
			},
		},
		model.Info("Connected to transporter."),
	)
}

var _ Connection = WebSocketConnection{}

type WebSocketConnection struct {
	schema      neoschema.TransporterSchema
	transporter *websocket.WebSocketTransporter
	client      *client.Client
	msgChan     chan tea.Msg
}

// Send implements [Connection].
func (w WebSocketConnection) Send(endpoint string, request any) tea.Cmd {
	return func() tea.Msg {
		route, ok := w.schema.Routes[endpoint]
		if !ok {
			return model.Error("couldn't find matching route")
		}

		// Send to the server based on the send type
		var res PackedAny = PackedAny{nil}
		var err error
		switch route.GetSendType() {
		case neoschema.SendOK:
			err = w.client.SendOk(endpoint, PackedAny{request})
		case neoschema.SendOKNoRequest:
			err = w.client.SendOkNoRequest(endpoint)
		case neoschema.SendRequestResponse:
			res, err = w.client.Send[PackedAny](endpoint, PackedAny{request})
		case neoschema.SendNoResponse:
			err = w.client.SendNoResponse(endpoint, PackedAny{request})
		case neoschema.SendPing:
			err = w.client.SendPing(endpoint)
		case neoschema.SendNoRequest:
			res, err = w.client.SendNoRequest[PackedAny](endpoint)
		}

		// Return an error / the response when we get back stuff from the server
		if err != nil {
			if u, ok := err.(*client.UserError); ok {
				return model.Error("%s returned error response: %s", endpoint, u.Error())
			}
			return model.Error("Sending %s failed: %v", endpoint, err)
		}
		if res.Value == nil { // Handler without response thingy
			return model.Response(endpoint, request, nil)
		}
		return model.Response(endpoint, request, res.Value)
	}
}

// WaitForEvent implements [Connection].
func (w WebSocketConnection) WaitForEvent() tea.Cmd {
	return func() tea.Msg {
		defer func() {
			recover() // This happens when no events come in any more (closed)
		}()

		// Wait for a message and also wait for the next one after
		msg := <-w.msgChan
		return model.Multiple(msg, DoWaitMsg{})
	}
}

// Close implements [Connection].
func (w WebSocketConnection) Close() {
	close(w.msgChan)
	w.transporter.Close()
}
