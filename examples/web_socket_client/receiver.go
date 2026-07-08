package main

import (
	"log"
)

func registerReceiver(r *definitions.MainConnector) {
	r.ReceiveNewPunSubmitted(func(event definitions.NewPunEvent) {
		log.Printf("A new pun was submitted by someone, it is %v\n", event.Pun)
	})
}
