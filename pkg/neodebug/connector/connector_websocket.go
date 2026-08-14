package connector

import (
	"net/url"

	"github.com/Liphium/neoroute/client"
	"github.com/Liphium/neoroute/client/transporter/websocket"
	"github.com/Liphium/neoroute/neoschema"
	"github.com/Liphium/neoroute/pkg/neodebug/config"
	tea "github.com/charmbracelet/bubbletea"
)

func connectWebsocket() tea.Msg {
	msgChan := make(chan tea.Msg)

	// Create the actual transporter
	r := client.NewReceiver(client.Config{}) // We don't really need the error handler as we handle errors below
	transporter := websocket.NewWebSocketTransporter(r)

	// Connect to the transporter using the URL in the config
	url, err := url.Parse(config.Config.TransporterURL)
	if err != nil {
		return e("Failed to parse URL: %w", err)
	}
	doneChan, err := transporter.Connect(url)
	if err != nil {
		return e("Failed to connect: %w", err)
	}
	go func() {
		<-doneChan
		msgChan <- ClosedMsg{}
	}()

	return ConnectedMsg{
		Connection: WebSocketConnection{
			transporter: transporter,
			recv:        r,
			msgChan:     msgChan,
		},
	}
}

var _ Connection = WebSocketConnection{}

type WebSocketConnection struct {
	schema      neoschema.TransporterSchema
	transporter *websocket.WebSocketTransporter
	recv        *client.Receiver
	msgChan     chan tea.Msg
}

// Send implements [Connection].
func (w WebSocketConnection) Send(id, endpoint string, request any) tea.Cmd {
	return func() tea.Msg {
		route, ok := w.schema.Routes[endpoint]
		if !ok {
			return reqErr(id, "couldn't find matching route")
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
			return reqErr(id, "something went wrong: %w", err)
		}
		if res == nil { // Handler without response thingy
			return ResponseReceivedMsg{ID: id, Data: nil}
		}
		return ResponseReceivedMsg{ID: id, Data: res.Value}
	}
}

// WaitForEvent implements [Connection].
func (w WebSocketConnection) WaitForEvent() tea.Cmd {
	return func() tea.Msg {
		return <-w.msgChan
	}
}

func (w WebSocketConnection) Close() {
	w.transporter.Close()
}
