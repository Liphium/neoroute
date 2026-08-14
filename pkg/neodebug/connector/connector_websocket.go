package connector

import (
	"fmt"
	"net/url"

	"github.com/Liphium/neoroute/client"
	"github.com/Liphium/neoroute/client/transporter/websocket"
	"github.com/Liphium/neoroute/pkg/neodebug/config"
	tea "github.com/charmbracelet/bubbletea"
)

func connectWebsocket() tea.Msg {
	msgChan := make(chan tea.Msg)

	// Create the actual transporter
	r := client.NewReceiver(client.Config{
		ErrorHandler: func(err error) {
			msgChan <- ErrorMsg{fmt.Sprintf("Something went wrong: %w", err)}
		},
	})
	transporter := websocket.NewWebSocketTransporter(r)

	// Connect to the transporter using the URL in the config
	url, err := url.Parse(config.Config.TransporterURL)
	if err != nil {
		return e(fmt.Sprintf("Failed to parse URL: %w", err))
	}
	closeChan, err := transporter.Connect(url)
	if err != nil {
		return e(fmt.Sprint("Failed to connect: %w", err))
	}

	return ConnectedMsg{
		Connection: WebSocketConnection{
			msgChan:   msgChan,
			closeChan: closeChan,
		},
	}
}

var _ Connection = WebSocketConnection{}

type WebSocketConnection struct {
	msgChan   chan tea.Msg
	closeChan chan struct{}
}

// Send implements [Connection].
func (w WebSocketConnection) Send(handler string) tea.Cmd {
	return func() tea.Msg {
		// TODO: Actually do some kind of sending
		return nil
	}
}

// WaitForEvent implements [Connection].
func (w WebSocketConnection) WaitForEvent() tea.Cmd {
	return func() tea.Msg {
		return <-w.msgChan
	}
}

func (w WebSocketConnection) Close() {
	w.closeChan <- struct{}{}
}
