package client

import (
	"fmt"
	"sync"
	"time"
)

type Sender interface {
	getRequestId() int
	SetSendFunc(sendFunc func(data []byte) error)
	Handle(reqData []byte)
	handleResponse(respBytes []byte)
	handleEvent(eventBytes []byte)
	removeResponseWaiter(reqId int)
	sendRequest(route string, reqData []byte, wantResponse bool) (chan response, int, error)
	getConfig() Config
	setReceiver(eventName string, receiveFunc func(c *Ctx))
}

type sender struct {
	mutex     sync.Mutex
	config    Config
	sendFunc  func(data []byte) error
	requestId int
	waiters   map[int]chan response
	receiver  map[string]func(*Ctx) // Only used with receiver
}

// NewSender returns an initialized sender
// Make sure to use a different sender for every transporter
func NewSender(config Config) Sender {
	if config.RequestTimeout == 0 {
		config.RequestTimeout = 5 * time.Second
	}
	return &sender{
		config:    config,
		requestId: 0,
		waiters:   make(map[int]chan response),
		receiver:  make(map[string]func(*Ctx)),
	}
}

func (s *sender) getRequestId() int {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.requestId++
	return s.requestId
}

func (s *sender) SetSendFunc(sendFunc func(data []byte) error) {
	s.mutex.Lock()
	s.sendFunc = sendFunc
	s.mutex.Unlock()
}

// Handle should be called by a transporter method when it receives a message.
// Make sure to call handle in a new go routine to avoid blocking.
// ONLY USE THIS WHEN IMPLEMENTING A TRANSPORTER.
func (s *sender) Handle(reqData []byte) {

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
		s.handleEvent(message.Data)
	case MessageTypeResponse:
		s.handleResponse(message.Data)
	default:
		Logger.Info("received unsupported message type", "type", message.Type)
		return
	}
}

func (s *sender) handleResponse(respBytes []byte) {
	var resp response
	_, err := resp.UnmarshalMsg(respBytes)
	if err != nil {
		Logger.Info("failed to unmarshal response", "err", err)
		return
	}

	if resp.Id == -1 && resp.IsError {
		s.config.ErrorHandler(fmt.Errorf("%s", string(resp.Data)))
		return
	}

	s.mutex.Lock()
	waiter, ok := s.waiters[resp.Id]
	s.mutex.Unlock()
	if !ok {
		Logger.Info("received response for non existing waiter", "response id", resp.Id)
		return
	}

	waiter <- resp
}

func (s *sender) handleEvent(eventBytes []byte) {
	var ev event
	_, err := ev.UnmarshalMsg(eventBytes)
	if err != nil {
		Logger.Info("failed to unmarshal event", "err", err)
		return
	}

	// Check if receiver exists
	s.mutex.Lock()
	receiver, ok := s.receiver[ev.Name]
	s.mutex.Unlock()
	if !ok {
		Logger.Info("event receiver doesn't exist", "receiver", ev.Name)
		return
	}

	c := &Ctx{
		data: ev.Data,
		name: ev.Name,
	}

	receiver(c)
}

func (s *sender) removeResponseWaiter(reqId int) {
	s.mutex.Lock()
	delete(s.waiters, reqId)
	s.mutex.Unlock()
}

func (s *sender) sendRequest(route string, reqData []byte, wantResponse bool) (chan response, int, error) {
	reqId := s.getRequestId()
	s.mutex.Lock()
	sendFunc := s.sendFunc
	s.mutex.Unlock()

	req := request{
		Id:    reqId,
		Route: route,
		Data:  reqData,
	}

	reqBytes, err := req.MarshalMsg(nil)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to marshal request: %v", err)
	}

	// Set response channel
	var respChan chan response
	if wantResponse {
		respChan = make(chan response)
		s.mutex.Lock()
		s.waiters[reqId] = respChan
		s.mutex.Unlock()
	}

	return respChan, reqId, sendFunc(reqBytes)
}

func (s *sender) getConfig() Config {
	return s.config
}

func (s *sender) setReceiver(eventName string, receiveFunc func(c *Ctx)) {
	s.receiver[eventName] = receiveFunc
}
