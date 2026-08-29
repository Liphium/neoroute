package main

import (
	"bufio"
	"fmt"
	"log"
	"log/slog"
	"net/url"
	"os"
	"strings"

	"github.com/Liphium/neoroute/client"
	"github.com/Liphium/neoroute/examples/chat_client_ws/generated"
)

func main() {

	// Create a ws connector from the generated code
	wsConn := generated.NewWsConnector(client.Config{
		ErrorHandler: func(err error) {
			slog.Error("ws error", "error", err)
			os.Exit(1)
		},
	})

	// Receive handler for all incoming MessageEvents
	wsConn.ReceiveMessage(func(event generated.MessageEvent) {
		fmt.Println("----------------")
		if sender := event.Sender; sender != "" {
			fmt.Println("From: ", event.Sender)
		}
		fmt.Println("Text: ", event.Text)
		fmt.Printf("Sent at: %d\n", event.Timestamp)
		fmt.Println("----------------")
	})

	// Collect input from console
	inputChan := make(chan string)
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Err() != nil {
			log.Println("error with input scanner", scanner.Err())
			return
		}
		for scanner.Scan() {
			inputChan <- scanner.Text()
		}
	}()

	fmt.Println("Type \"exit\" to quit the program.")
	fmt.Println("Type \"name \" followed by your name to set your username.")
	fmt.Println("Type \"msg \" followed by a message to broadcast it to all connected clients.")

	// Connect to server using WebSocket
	u, _ := url.Parse("ws://localhost:6121/ws")
	done, err := wsConn.Connect(u)
	if err != nil {
		log.Println("error: ", err)
		return
	}

	var name string = "guest"

	func() {
		for {
			select {
			case <-done:
				return
			case input := <-inputChan:
				if newName, isName := strings.CutPrefix(input, "name "); isName {
					name = newName
				} else if msg, isMsg := strings.CutPrefix(input, "msg "); isMsg {

					// Send the message to the server
					go wsConn.SendSend(generated.SendRequest{
						Sender: name,
						Text:   msg,
					})
				} else if input == "exit" {
					if err := wsConn.Close(); err != nil {
						slog.Error("failed to close connection", "error", err)
						os.Exit(1)
					}
					<-done
					return
				}
			}
		}
	}()
}
