package main

import (
	"log"
	"web_socket/definitions"
)

func registerReceiver(r *definitions.MainConnector) {
	r.ReceiveNewPunSubmitted(func(event definitions.NewPunEvent) {
		log.Printf("A new pun was submitted by someone, it is %v\n", event.Pun)
	})
}
