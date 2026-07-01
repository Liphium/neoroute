package main

import (
	"log"

	"github.com/Liphium/neoroute/client"
)

func SendSubmitPunRequest(c *MainConnector, pun string) {
	sendErr := c.SendSubmitPun(SubmitPunRequest{Pun: pun})
	if sendErr != nil {
		if reqErr, ok := sendErr.(*client.UserError); ok {
			log.Println("Couldn't submit pun because", reqErr)
		} else {
			log.Println("failed to send submit pun request: ", sendErr)
		}
	}
}

func SendEchoRequest(c *MainConnector, message string) {
	resp, sendErr := c.SendEcho(EchoRequest{Message: message})
	if sendErr != nil {
		if reqErr, ok := sendErr.(*client.UserError); ok {
			log.Println("Echo failed because", reqErr)
		}
		log.Println("failed to send submit pun request: ", sendErr)
	} else {
		log.Printf("Received %v. echo: %v\n", resp.RequestNumber, resp.Message)
	}
}
