package neoroute

import (
	"fmt"
)

const (
	MessageTypeResponse = 0
	MessageTypeEvent    = 1
)

//go:generate msgp -unexported

type Message struct {
	Type int    `msg:"type"` // Response or event
	Data []byte `msg:"data"`
}

func MessageEvent(event Event) ([]byte, error) {

	// Marshal event data
	eventData, err := event.MarshalMsg(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal event data: %v", err)
	}

	msg := Message{
		Type: MessageTypeEvent,
		Data: eventData,
	}

	msgBytes, err := msg.MarshalMsg(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal message: %v", err)
	}

	return msgBytes, nil
}

func messageResponse[D any](neo *NeoRouter[D], resp response) []byte {

	// Marshal response data
	respData, err := resp.MarshalMsg(nil)
	if err != nil {
		resp := response{
			Id:      resp.Id,
			IsError: true,
			Data:    []byte(neo.config.ErrorHandler(err)),
		}
		respData, err = resp.MarshalMsg(nil)
		if err != nil {
			// This should never happen. If this fails for you please open a bug report on our github.
			panic(fmt.Sprintf("failed to marshal error response: %v", err))
		}

	}

	msg := Message{
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

type Event struct {
	Name string `msg:"name"`
	Data []byte `msg:"data"`
}

type response struct {
	Id      int    `msg:"id"`
	IsError bool   `msg:"error"`
	Data    []byte `msg:"data"`
}

func (r response) Error() string {
	return ""
}

func (r response) Is(target error) bool {
	_, ok := target.(response)
	if ok {
		return true
	}
	_, ok = target.(*response)
	return ok
}

// This type is used for routes that have no response so no error is thrown.
type noResponse struct{}

func (r noResponse) Error() string {
	return ""
}

func (r noResponse) Is(target error) bool {
	_, ok := target.(noResponse)
	if ok {
		return true
	}
	_, ok = target.(*noResponse)
	return ok
}
