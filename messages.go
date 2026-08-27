package neoroute

import (
	"fmt"
)

const (
	messageTypeResponse = 0
	messageTypeEvent    = 1
)

//go:generate msgp -unexported

type message struct {
	Type int    `msg:"type"` // Response or event
	Data []byte `msg:"data"`
}

// messageEvent marshals an event into a message.
func messageEvent(event event) []byte {

	// Marshal event data
	eventData, err := event.MarshalMsg(nil)
	if err != nil {
		// This should never happen. If this fails for you please open a bug report on our github.
		panic(fmt.Sprintf("failed to marshal event data: %v, This should never happen. If this fails for you please open a bug report on our github.", err))
	}

	msg := message{
		Type: messageTypeEvent,
		Data: eventData,
	}

	msgBytes, err := msg.MarshalMsg(nil)
	if err != nil {
		// This should never happen. If this fails for you please open a bug report on our github.
		panic(fmt.Sprintf("failed to marshal event message: %v, This should never happen. If this fails for you please open a bug report on our github.", err))
	}

	return msgBytes
}

// messageResponse marshals a response into a message.
func messageResponse(resp response) []byte {

	// Marshal response data
	respData, err := resp.MarshalMsg(nil)
	if err != nil {
		// This should never happen. If this fails for you please open a bug report on our github.
		panic(fmt.Sprintf("failed to marshal response: %v, This should never happen. If this fails for you please open a bug report on our github.", err))
	}

	msg := message{
		Type: messageTypeResponse,
		Data: respData,
	}

	msgBytes, err := msg.MarshalMsg(nil)
	if err != nil {

		// This should never happen. If this fails for you please open a bug report on our github.
		panic(fmt.Sprintf("failed to marshal response message: %v, This should never happen. If this fails for you please open a bug report on our github.", err))
	}

	return msgBytes
}
