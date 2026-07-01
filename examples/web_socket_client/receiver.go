package main

import (
	"log"
)

func registerReceiver(r *MainConnector) {
	r.ReceiveNewPunSubmitted(func(event NewPunEvent) {
		log.Printf("A new pun was submitted by someone, it is %v\n", event.Pun)
	})
}
