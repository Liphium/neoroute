package main

import (
	"bufio"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/Liphium/neoroute/client"
)

func main() {
	r := client.NewReceiver(client.Config{
		ErrorHandler: func(err error) {
			log.Println("error with receiver", err)
		},
		RequestTimeout: time.Second * 1,
	})

	registerReceiver(r)

	t := client.NewWebSocketTransporter(r)

	u, err := url.Parse("ws://localhost:6121/")
	if err != nil {
		log.Fatalf("failed to parse url: %v", err)
	}
	done, err := t.Connect(u)
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}

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

	fmt.Println("Type exit to quit the program.")
	fmt.Println("Type \"Pun \" followed by a pun to submit it to the server.")
	fmt.Println("Type \"Echo \" followed by a message to receive a echo message from the server.")

	func() {
		for {
			select {
			case <-done:
				return
			case input := <-inputChan:
				if pun, isPun := strings.CutPrefix(input, "Pun "); isPun {
					fmt.Println("Sending pun to server")
					go SendSubmitPunRequest(r, pun)
				} else if msg, isEcho := strings.CutPrefix(input, "Echo "); isEcho {
					fmt.Println("Sending echo to server")
					go SendEchoRequest(r, msg)
				} else if input == "exit" {
					if err := t.Close(); err != nil {
						log.Println("error: ", err)
					}
					<-done
					return
				}
			}
		}
	}()

	log.Println("Connection closed")
}
