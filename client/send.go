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
}](r Sender, route string, req RQ) (RS, error) {

	var resp RS

	reqBytes, err := marshalRequestData[RQ, RQP](req)
	if err != nil {
		return resp, err
	}

	respChan, reqId, err := r.sendRequest(route, reqBytes, true)
	if err != nil {
		return resp, err
	}

	// Wait for time out duration for a response and remove chan after
	select {
	case res := <-respChan:
		r.removeResponseWaiter(reqId)

		if res.IsError {
			return resp, NewUserError(string(res.Data))
		}

		resp, err = unmarshalResponseData[RS, RSP](r, res.Data)
		return resp, err

	case <-time.After(r.getConfig().RequestTimeout):
		r.removeResponseWaiter(reqId)
		return resp, fmt.Errorf("waiting for response timed out after %v", r.getConfig().RequestTimeout)
	}

}

func SendOk[RQ any, RQP interface {
	*RQ
	msgp.Marshaler
}](r Sender, route string, req RQ) error {

	reqBytes, err := marshalRequestData[RQ, RQP](req)
	if err != nil {
		return err
	}

	respChan, reqId, err := r.sendRequest(route, reqBytes, true)
	if err != nil {
		return err
	}

	// Wait fpr time out duration for a response and remove chan after
	select {
	case res := <-respChan:
		r.removeResponseWaiter(reqId)

		if res.IsError {
			return NewUserError(string(res.Data))
		}

		return nil

	case <-time.After(r.getConfig().RequestTimeout):
		r.removeResponseWaiter(reqId)
		return fmt.Errorf("waiting for response timed out after %v", r.getConfig().RequestTimeout)
	}

}

func SendOkNoRequest(r Sender, route string) error {

	respChan, reqId, err := r.sendRequest(route, []byte{}, true)
	if err != nil {
		return err
	}

	// Wait fpr time out duration for a response and remove chan after
	select {
	case res := <-respChan:
		r.removeResponseWaiter(reqId)

		if res.IsError {
			return NewUserError(string(res.Data))
		}

		return nil

	case <-time.After(r.getConfig().RequestTimeout):
		r.removeResponseWaiter(reqId)
		return fmt.Errorf("waiting for response timed out after %v", r.getConfig().RequestTimeout)
	}

}

func SendNoRequest[RS any, RSP interface {
	*RS
	msgp.Unmarshaler
}](r Sender, route string) (RS, error) {

	var resp RS

	respChan, reqId, err := r.sendRequest(route, []byte{}, true)
	if err != nil {
		return resp, err
	}

	// Wait fpr time out duration for a response and remove chan after
	select {
	case res := <-respChan:
		r.removeResponseWaiter(reqId)

		if res.IsError {
			return resp, NewUserError(string(res.Data))
		}

		resp, err = unmarshalResponseData[RS, RSP](r, res.Data)
		return resp, err

	case <-time.After(r.getConfig().RequestTimeout):
		r.removeResponseWaiter(reqId)
		return resp, fmt.Errorf("waiting for response timed out after %v", r.getConfig().RequestTimeout)
	}
}

func SendNoResponse[RQ any, RQP interface {
	*RQ
	msgp.Marshaler
}](r Sender, route string, req RQ) error {

	reqBytes, err := marshalRequestData[RQ, RQP](req)
	if err != nil {
		return err
	}

	_, _, err = r.sendRequest(route, reqBytes, false)
	return err
}

func SendNoop(r Sender, route string) error {
	_, _, err := r.sendRequest(route, []byte{}, false)
	return err
}
