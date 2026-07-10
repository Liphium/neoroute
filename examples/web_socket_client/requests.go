package main

import (
	"log"
	"web_socket/definitions"

	"github.com/Liphium/neoroute/client"
)

func SendSubmitPunRequest(c *definitions.MainConnector, pun string) {
	sendErr := c.SendSubmitPun(definitions.SubmitPunRequest{Pun: pun})
	if sendErr != nil {
		if reqErr, ok := sendErr.(*client.UserError); ok {
			log.Println("Couldn't submit pun because", reqErr)
		} else {
			log.Println("failed to send submit pun request: ", sendErr)
		}
	}
}

func SendEchoRequest(c *definitions.MainConnector, message string) {
	resp, sendErr := c.SendEcho(definitions.EchoRequest{Message: message})
	if sendErr != nil {
		if reqErr, ok := sendErr.(*client.UserError); ok {
			log.Println("Echo failed because", reqErr)
		}
		log.Println("failed to send submit pun request: ", sendErr)
	} else {
		log.Printf("Received %v. echo: %v\n", resp.RequestNumber, resp.Message)
	}
}
