package client

import (
	"fmt"
	"sync"
	"time"
)

type Handler struct {
	mutex     *sync.Mutex
	config    Config
	sendFunc  func(data []byte) error
	requestId int
	waiters   map[int]chan response
	handler   map[string]func(*Ctx) // Only used with receiver
}

// NewHandler returns an initialized handler
// Make sure to use a different handler for every transporter
func NewHandler(config Config) *Handler {
	if config.RequestTimeout == 0 {
		config.RequestTimeout = 5 * time.Second
	}
	return &Handler{
		mutex:     &sync.Mutex{},
		config:    config,
		requestId: 0,
		waiters:   make(map[int]chan response),
		handler:   make(map[string]func(*Ctx)),
	}
}

func (h *Handler) getRequestId() int {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.requestId++
	return h.requestId
}

func (h *Handler) SetSendFunc(sendFunc func(data []byte) error) {
	h.mutex.Lock()
	h.sendFunc = sendFunc
	h.mutex.Unlock()
}

// handle should be called by a transporter method when it receives a message.
// Make sure to call handle in a new go routine to avoid blocking.
// ONLY USE THIS WHEN IMPLEMENTING A TRANSPORTER.
func (h *Handler) Handle(reqData []byte) {

	// Unmarshal request data to message
	var message message
	_, err := message.UnmarshalMsg(reqData)
	if err != nil {
		Logger.Info("failed to unmarshal message", "err", err)
		return
	}

	// Check how to handle message
	switch message.Type {
	case MessageTypeEvent:
		h.handleEvent(message.Data)
	case MessageTypeResponse:
		h.handleResponse(message.Data)
	default:
		Logger.Info("received unsupported message type", "type", message.Type)
		return
	}
}

func (h *Handler) handleResponse(respBytes []byte) {
	var resp response
	_, err := resp.UnmarshalMsg(respBytes)
	if err != nil {
		Logger.Info("failed to unmarshal response", "err", err)
		return
	}

	if resp.Id == -1 && resp.IsError {
		h.config.ErrorHandler(fmt.Errorf("%s", string(resp.Data)))
		return
	}

	h.mutex.Lock()
	waiter, ok := h.waiters[resp.Id]
	h.mutex.Unlock()
	if !ok {
		Logger.Info("received response for non existing waiter", "response id", resp.Id)
		return
	}

	waiter <- resp
}

func (h *Handler) handleEvent(eventBytes []byte) {
	var ev event
	_, err := ev.UnmarshalMsg(eventBytes)
	if err != nil {
		Logger.Info("failed to unmarshal event", "err", err)
		return
	}

	// Check if handler exists
	h.mutex.Lock()
	handler, ok := h.handler[ev.Name]
	h.mutex.Unlock()
	if !ok {
		Logger.Info("event handler doesn't exist", "handler", ev.Name)
		return
	}

	c := &Ctx{
		data: ev.Data,
		name: ev.Name,
	}

	handler(c)
}
