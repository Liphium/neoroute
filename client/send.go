package client

import (
	"fmt"
	"time"

	"github.com/tinylib/msgp/msgp"
)

func Send[RS any, RSP interface {
	*RS
	msgp.Unmarshaler
}, RQ any, RQP interface {
	*RQ
	msgp.Marshaler
}](r *Handler, route string, req RQ) (RS, string, error) {

	var resp RS

	reqBytes, err := marshalRequestData[RQ, RQP](r, req)
	if err != nil {
		return resp, "", err
	}

	respChan, reqId, err := sendRequest(r, route, reqBytes, true)
	if err != nil {
		return resp, "", err
	}

	// Wait fpr time out duration for a response and remove chan after
	select {
	case res := <-respChan:
		removeResponseWaiter(r, reqId)

		if res.IsError {
			return resp, string(res.Data), nil
		}

		resp, err = unmarshalResponseData[RS, RSP](r, res.Data)
		return resp, "", err

	case <-time.After(r.config.RequestTimeout):
		removeResponseWaiter(r, reqId)
		return resp, "", fmt.Errorf("waiting for response timed out after %v", r.config.RequestTimeout)
	}

}

func SendOk[RQ any, RQP interface {
	*RQ
	msgp.Marshaler
}](r *Handler, route string, req RQ) (string, error) {

	reqBytes, err := marshalRequestData[RQ, RQP](r, req)
	if err != nil {
		return "", err
	}

	respChan, reqId, err := sendRequest(r, route, reqBytes, true)
	if err != nil {
		return "", err
	}

	// Wait fpr time out duration for a response and remove chan after
	select {
	case res := <-respChan:
		removeResponseWaiter(r, reqId)

		if res.IsError {
			return string(res.Data), nil
		}

		return "", nil

	case <-time.After(r.config.RequestTimeout):
		removeResponseWaiter(r, reqId)
		return "", fmt.Errorf("waiting for response timed out after %v", r.config.RequestTimeout)
	}

}

func SendOkNoRequest(r *Handler, route string) (string, error) {

	respChan, reqId, err := sendRequest(r, route, []byte{}, true)
	if err != nil {
		return "", err
	}

	// Wait fpr time out duration for a response and remove chan after
	select {
	case res := <-respChan:
		removeResponseWaiter(r, reqId)

		if res.IsError {
			return string(res.Data), nil
		}

		return "", nil

	case <-time.After(r.config.RequestTimeout):
		removeResponseWaiter(r, reqId)
		return "", fmt.Errorf("waiting for response timed out after %v", r.config.RequestTimeout)
	}

}

func SendNoRequest[RS any, RSP interface {
	*RS
	msgp.Unmarshaler
}](r *Handler, route string) (RS, string, error) {

	var resp RS

	respChan, reqId, err := sendRequest(r, route, []byte{}, true)
	if err != nil {
		return resp, "", err
	}

	// Wait fpr time out duration for a response and remove chan after
	select {
	case res := <-respChan:
		removeResponseWaiter(r, reqId)

		if res.IsError {
			return resp, string(res.Data), nil
		}

		resp, err = unmarshalResponseData[RS, RSP](r, res.Data)
		return resp, "", err

	case <-time.After(r.config.RequestTimeout):
		removeResponseWaiter(r, reqId)
		return resp, "", fmt.Errorf("waiting for response timed out after %v", r.config.RequestTimeout)
	}

}

func SendNoResponse[RQ any, RQP interface {
	*RQ
	msgp.Marshaler
}](r *Handler, route string, req RQ) error {

	reqBytes, err := marshalRequestData[RQ, RQP](r, req)
	if err != nil {
		return err
	}

	_, _, err = sendRequest(r, route, reqBytes, false)
	return err
}

func SendNoop(r *Handler, route string) error {
	_, _, err := sendRequest(r, route, []byte{}, false)
	return err
}
