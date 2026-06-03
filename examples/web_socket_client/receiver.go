package main

import (
	"log"

	"github.com/Liphium/neoroute/client"
)

func registerReceiver(r *client.Receiver) {
	client.Receive(r, "new_pun_submitted", func(c *client.Ctx, req NewPunEvent) {
		log.Printf("A new pun was submitted by someone, it is %v\n", req.Pun)
	})
}
