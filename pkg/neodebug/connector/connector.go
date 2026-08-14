package connector

import (
	"github.com/Liphium/neoroute/neoschema"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tinylib/msgp/msgp"
)

type ConnectedMsg struct {
	Connection Connection
}

type ErrorMsg struct {
	Error string
}

func e(err string) tea.Msg {
	return ErrorMsg{err}
}

func Connect(transporter neoschema.TransporterSchema) tea.Cmd {
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

	// This should send a handler and return a proper message (failed / result)
	Send(handler string) tea.Cmd

	// Should wait for events / responses, once one is received we return a message here containing it.
	WaitForEvent() tea.Cmd

	// Should close the connection with the server (in case doable)
	Close()
}

type AnyUnmarshaler struct {
	Value any
}

func (u *AnyUnmarshaler) UnmarshalMsg(b []byte) ([]byte, error) {
	v, rest, err := msgp.ReadIntfBytes(b)
	if err != nil {
		return rest, err
	}

	u.Value = v
	return rest, nil
}
