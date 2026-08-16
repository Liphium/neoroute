package connector

import (
	tea "charm.land/bubbletea/v2"
	"github.com/Liphium/neoroute/neoschema"
	"github.com/Liphium/neoroute/pkg/neodebug/model"
)

type ClosedMsg struct{}

// When this message is sent, we're connected, also return connection.WaitForEvent() afterwards and handle DoWaitMsg after that as well.
type ConnectedMsg struct {
	Connection Connection
}

// When this message is sent, you should call connection.WaitForEvent() so we can send more events into the connection.
type DoWaitMsg struct{}

func Connect(transporter neoschema.TransporterSchema) tea.Cmd {

	// TODO: Set the logger inside of neoroute to one that just pipes everything right into our history or sth

	return func() tea.Msg {
		switch transporter.Type {
		case neoschema.TransporterHTTP:
			return withClose(model.Error("HTTP request transporters are currently not supported."))
		case neoschema.TransporterWebSocket:
			return connectWebsocket()
		}
		return withClose(model.Error("This type of transporter is unknown."))
	}
}

func withClose(msg tea.Msg) tea.Msg {
	return model.Multiple(msg, ClosedMsg{})
}

type Connection interface {

	// This should send to a route and return a proper message (failed / result)
	Send(route string, request any) tea.Cmd

	// Should wait for events / responses, once one is received we return a message here containing it.
	WaitForEvent() tea.Cmd

	// Should close the connection with the server (in case doable)
	Close()
}
