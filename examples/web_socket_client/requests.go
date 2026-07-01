package main

import (
	"log"

	"github.com/Liphium/neoroute/client"
)

func SendSubmitPunRequest(r *client.Receiver, pun string) {
	sendErr := client.SendOk(r, "submit_pun", SubmitPunRequest{Pun: pun})
	if sendErr != nil {
		if reqErr, ok := sendErr.(*client.UserError); ok {
			log.Println("Couldn't submit pun because", reqErr)
		} else {
			log.Println("failed to send submit pun request: ", sendErr)
		}
	}
}

func SendEchoRequest(r *client.Receiver, message string) {
	resp, sendErr := client.Send[EchoResponse](r, "echo", EchoRequest{Message: message})
	if sendErr != nil {
		if reqErr, ok := sendErr.(*client.UserError); ok {
			log.Println("Echo failed because", reqErr)
		}
		log.Println("failed to send submit pun request: ", sendErr)
	} else {
		log.Printf("Received %v. echo: %v\n", resp.RequestNumber, resp.Message)
	}
}
