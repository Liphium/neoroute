package main

//go:generate msgp

// Outgoing

type EchoRequest struct {
	Message string `msg:"message"`
}

type SubmitPunRequest struct {
	Pun    string `msg:"pun"`
	Number Test   `msg:"number"`
}

type Test struct {
	World  string `msg:"world"`
	Test   string `msg:"test"`
	World2 string `msg:"world2"`
	Test2  string `msg:"test2"`
	World3 string `msg:"world3"`
	Test3  string `msg:"test3"`
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
