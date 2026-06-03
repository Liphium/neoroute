package client

import (
	"fmt"

	"github.com/tinylib/msgp/msgp"
)

func removeResponseWaiter(r *Handler, reqId int) {
	r.mutex.Lock()
	delete(r.waiters, reqId)
	r.mutex.Unlock()
}

func marshalRequestData[RQ any, RQP interface {
	*RQ
	msgp.Marshaler
}](r *Handler, req RQ) ([]byte, error) {
	marshaler := any(&req).(msgp.Marshaler)
	reqData, err := marshaler.MarshalMsg(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request data: %v", err)
	}
	return reqData, nil
}

func unmarshalResponseData[RS any, RSP interface {
	*RS
	msgp.Unmarshaler
}](r *Handler, respBytes []byte) (RS, error) {
	var data RS
	unmarshaler := any(&data).(msgp.Unmarshaler)
	_, err := unmarshaler.UnmarshalMsg(respBytes)
	if err != nil {
		return data, fmt.Errorf("failed to unmarshal response data: %v", err)
	}
	return data, nil
}

func sendRequest(r *Handler, route string, reqData []byte, wantResponse bool) (chan response, int, error) {
	reqId := r.getRequestId()
	r.mutex.Lock()
	sendFunc := r.sendFunc
	r.mutex.Unlock()

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
		r.mutex.Lock()
		r.waiters[reqId] = respChan
		r.mutex.Unlock()
	}

	return respChan, reqId, sendFunc(reqBytes)
}
