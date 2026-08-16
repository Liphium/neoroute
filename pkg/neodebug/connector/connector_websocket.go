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

func connectWebsocket() tea.Msg {
	msgChan := make(chan tea.Msg)

	// Create the actual transporter
	r := client.NewReceiver(client.Config{}) // We don't really need the error handler as we handle errors below
	transporter := websocket.NewWebSocketTransporter(r)

	// Connect to the transporter using the URL in the config
	url, err := url.Parse(config.Config.TransporterURL)
	if err != nil {
		return withClose(model.Error("Failed to parse URL: %w", err))
	}
	doneChan, err := transporter.Connect(url)
	if err != nil {
		return withClose(model.Error(err.Error()))
	}
	go func() {
		<-doneChan
		msgChan <- ClosedMsg{}
	}()

	return model.Multiple(
		ConnectedMsg{
			Connection: WebSocketConnection{
				transporter: transporter,
				recv:        r,
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
	recv        *client.Receiver
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
		var res *PackedAny
		var err error
		switch route.GetSendType() {
		case neoschema.SendOK:
			err = client.SendOk(w.recv, endpoint, PackedAny{request})
		case neoschema.SendOKNoRequest:
			err = client.SendOkNoRequest(w.recv, endpoint)
		case neoschema.SendRequestResponse:
			*res, err = client.Send[PackedAny](w.recv, endpoint, PackedAny{request})
		case neoschema.SendNoResponse:
			err = client.SendNoResponse(w.recv, endpoint, PackedAny{request})
		case neoschema.SendPing:
			err = client.SendPing(w.recv, endpoint)
		case neoschema.SendNoRequest:
			*res, err = client.SendNoRequest[PackedAny](w.recv, endpoint)
		}

		// Return an error / the response when we get back stuff from the server
		if err != nil {
			return model.Error("Sending %s failed: %w", endpoint, err)
		}
		if res == nil { // Handler without response thingy
			return model.Response(endpoint, nil)
		}
		return model.Response(endpoint, res.Value)
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

func (w WebSocketConnection) Close() {
	close(w.msgChan)
	w.transporter.Close()
}
