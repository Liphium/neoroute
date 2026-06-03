package main

import (
	"log"

	"github.com/Liphium/neoroute/client"
)

func registerReceiver(r *client.Receiver) {
	client.Receive(r, "new_pun_submitted", func(c *client.Ctx, req NewPunEvent) error {
		log.Printf("A new was submitted by someone, it is %v\n", req.Pun)
		return nil
	})
}
