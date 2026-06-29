// Code generated with neogen-generated v1 schema by neogen. DO NOT EDIT.
package main

import "fmt"

type PunsConnector struct{}

func NewPunsConnector() *PunsConnector {
	return &PunsConnector{}
}

func (c *PunsConnector) Connect() {
	fmt.Println("Hello, neogen!")
}

func (c *PunsConnector) ReceiveNewPunSubmitted(handler func(event NewPunEvent)) {
	fmt.Println("Handling some event!")
}

func (c *PunsConnector) SendEcho(payload EchoRequest) (EchoResponse, error) {
	fmt.Println("Sending some event!")
	// TODO: Return
}

func (c *PunsConnector) SendSubmitPun(payload SubmitPunRequest) error {
	fmt.Println("Sending some event!")
	// TODO: Return
}
