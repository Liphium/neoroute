package connector

import (
	"fmt"

	"github.com/Liphium/neoroute/neoschema"
	tea "charm.land/bubbletea/v2"
)

type ClosedMsg struct{}

type ConnectedMsg struct {
	Connection Connection
}

type EventReceivedMsg struct {
	Name string
	Data any
}

// Sent with the ID of the request when a response returns data. When the data is nil, nothing has been returned.
type ResponseReceivedMsg struct {
	ID   string
	Data any
}

type RequestErrorMsg struct {
	ID    string
	Error string
}

type ErrorMsg struct {
	Error string
}

func e(format string, a ...any) tea.Msg {
	return ErrorMsg{fmt.Sprintf(format, a...)}
}

func reqErr(id string, format string, a ...any) tea.Msg {
	return RequestErrorMsg{id, fmt.Sprintf(format, a...)}
}

func Connect(transporter neoschema.TransporterSchema) tea.Cmd {

	// TODO: Set the logger inside of neoroute to one that just pipes everything right into our history or sth

	return func() tea.Msg {
		switch transporter.Type {
		case neoschema.TransporterHTTP:
			return ErrorMsg{"HTTP request transporters are currently not supported."}
		case neoschema.TransporterWebTransport:
			return connectWebsocket()
		}
		return nil
	}
}

type Connection interface {

	// This should send to a route and return a proper message (failed / result)
	Send(id string, route string, request any) tea.Cmd

	// Should wait for events / responses, once one is received we return a message here containing it.
	WaitForEvent() tea.Cmd

	// Should close the connection with the server (in case doable)
	Close()
}
