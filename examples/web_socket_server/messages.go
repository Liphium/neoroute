package main

//go:generate msgp

// Outgoing

type EchoRequest struct {
	Message string `msg:"message"`
}

type SubmitPunRequest struct {
	Pun string `msg:"pun"`
}

// Response

type EchoResponse struct {
	RequestNumber int    `msg:"requestNumber"`
	Message       string `msg:"message"`
}

// Event
type NewPunEvent struct {
	Pun string `msg:"pun"`
}
