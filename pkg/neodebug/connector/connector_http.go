package connector

import (
	"net/url"

	tea "charm.land/bubbletea/v2"
	"github.com/Liphium/neoroute/client"
	"github.com/Liphium/neoroute/client/transporter/http"
	"github.com/Liphium/neoroute/neoschema"
	"github.com/Liphium/neoroute/pkg/neodebug/config"
	"github.com/Liphium/neoroute/pkg/neodebug/model"
)

func connectHTTP(schema neoschema.TransporterSchema) tea.Msg {
	msgChan := make(chan tea.Msg)

	// Connect to the transporter using the URL in the config
	url, err := url.Parse(config.Config.TransporterURL)
	if err != nil {
		return withClose(model.Error("Failed to parse URL: %v", err))
	}

	// Create the actual transporter
	r := client.NewReceiver(client.Config{}) // We don't really need the error handler as we handle errors below
	transporter := http.NewHTTPTransporter(r, "POST", url)

	return model.Multiple(
		ConnectedMsg{
			Connection: HTTPConnection{
				transporter: transporter,
				recv:        r,
				schema:      schema,
				msgChan:     msgChan,
			},
		},
		model.Info("Initialized transporter."),
	)
}

var _ Connection = HTTPConnection{}

type HTTPConnection struct {
	schema      neoschema.TransporterSchema
	transporter *http.HTTPTransporter
	recv        *client.Receiver
	msgChan     chan tea.Msg
}

// Send implements [Connection].
func (h HTTPConnection) Send(endpoint string, request any) tea.Cmd {
	return func() tea.Msg {
		route, ok := h.schema.Routes[endpoint]
		if !ok {
			return model.Error("couldn't find matching route")
		}

		// Send to the server based on the send type
		var res PackedAny = PackedAny{nil}
		var err error
		switch route.GetSendType() {
		case neoschema.SendOK:
			err = client.SendOk(h.recv, endpoint, PackedAny{request})
		case neoschema.SendOKNoRequest:
			err = client.SendOkNoRequest(h.recv, endpoint)
		case neoschema.SendRequestResponse:
			res, err = client.Send[PackedAny](h.recv, endpoint, PackedAny{request})
		case neoschema.SendNoResponse:
			err = client.SendNoResponse(h.recv, endpoint, PackedAny{request})
		case neoschema.SendPing:
			err = client.SendPing(h.recv, endpoint)
		case neoschema.SendNoRequest:
			res, err = client.SendNoRequest[PackedAny](h.recv, endpoint)
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
func (h HTTPConnection) WaitForEvent() tea.Cmd {
	return func() tea.Msg {
		defer func() {
			recover() // This happens when no events come in any more (closed)
		}()

		// Wait for a message and also wait for the next one after
		msg := <-h.msgChan
		return model.Multiple(msg, DoWaitMsg{})
	}
}

// Close implements [Connection].
func (h HTTPConnection) Close() {
	close(h.msgChan)
}
