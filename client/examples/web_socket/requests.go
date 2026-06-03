package main

import (
	"log"

	"github.com/Liphium/neoroute/client"
)

func SendSubmitPunRequest(r *client.Receiver, pun string) {
	reqErr, sendErr := client.SendOk(&r.Handler, "submit_pun", SubmitPunRequest{Pun: pun})
	if sendErr != nil {
		log.Println("failed to send submit pun request: ", sendErr)
	} else if reqErr != "" {
		log.Println("Couldn't submit pun because", reqErr)
	}
}

func SendEchoRequest(r *client.Receiver, message string) {
	resp, reqErr, sendErr := client.Send[EchoResponse](&r.Handler, "echo", EchoRequest{Message: message})
	if sendErr != nil {
		log.Println("failed to send submit pun request: ", sendErr)
	} else if reqErr != "" {
		log.Println("Echo failed because", reqErr)
	} else {
		log.Printf("Received %v. echo: %v\n", resp.RequestNumber, resp.Message)
	}
}
