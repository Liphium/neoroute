package neoroute

import (
	"fmt"
)

const (
	MessageTypeResponse = 0
	MessageTypeEvent    = 1
)

//go:generate msgp -unexported

type message struct {
	Type int    `msg:"type"` // Response or event
	Data []byte `msg:"data"`
}

func messageEvent(event event) ([]byte, error) {

	// Marshal event data
	eventData, err := event.MarshalMsg(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal event data: %v", err)
	}

	msg := message{
		Type: MessageTypeEvent,
		Data: eventData,
	}

	msgBytes, err := msg.MarshalMsg(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal message: %v", err)
	}

	return msgBytes, nil
}

func messageResponse(resp response) []byte {

	// Marshal response data
	respData, err := resp.MarshalMsg(nil)
	if err != nil {
		// This should never happen. If this fails for you please open a bug report on our github.
		panic(fmt.Sprintf("failed to marshal response: %v", err))
	}

	msg := message{
		Type: MessageTypeResponse,
		Data: respData,
	}

	msgBytes, err := msg.MarshalMsg(nil)
	if err != nil {

		// This should never happen. If this fails for you please open a bug report on our github.
		panic(fmt.Sprintf("failed to marshal response message: %v", err))
	}

	return msgBytes
}
