package main

import (
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"time"

	"github.com/Liphium/neoroute/client"
	"github.com/Liphium/neoroute/examples/chat_client_http/generated"
)

func main() {
	if length := len(os.Args); length < 2 || length > 3 || (length == 2 && os.Args[1] == "--help") {
		fmt.Println("Usage: go run . <message> <name  (optional)>")
		return
	}

	message := os.Args[1]
	name := ""
	if len(os.Args) == 3 {
		name = os.Args[2]
	}

	u, _ := url.Parse("http://localhost:6121/http")

	// Create a new HTTP connector from the generated code
	httpConn := generated.NewHttpConnector(client.Config{
		ErrorHandler: func(err error) {
			slog.Error("error", "message", err)
			os.Exit(1)
		},
		RequestTimeout: 10 * time.Second,
	}, "POST", u)

	// Use the generated connector to send the message
	sendErr := httpConn.SendSend(generated.SendRequest{
		Sender: name,
		Text:   message,
	})
	if sendErr != nil {
		slog.Error("error", "message", sendErr)
		os.Exit(1)
	}
}
