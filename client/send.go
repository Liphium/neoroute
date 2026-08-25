package client

import (
	"fmt"
	"time"

	"github.com/tinylib/msgp/msgp"
)

// Send sends a request to the server and returns a response.
//
// The error returned can be a *client.UserError in case the server returned a error message.
func (c *Client) Send[RS any, RSP msgp.UnmarshalPtr[RS], RQ msgp.Marshaler](route string, req RQ) (RS, error) {
	var resp RS

	reqBytes, err := marshalRequestData(req)
	if err != nil {
		return resp, err
	}

	respChan, reqId, err := c.sendRequest(route, reqBytes, true)
	if err != nil {
		return resp, err
	}

	// Wait for time out duration for a response and remove chan after
	select {
	case res := <-respChan:
		c.removeResponseWaiter(reqId)

		if res.IsError {
			return resp, NewUserError(string(res.Data))
		}

		resp, err = unmarshalResponseData[RS, RSP](res.Data)
		return resp, err

	case <-time.After(c.getConfig().RequestTimeout):
		c.removeResponseWaiter(reqId)
		return resp, fmt.Errorf("waiting for response timed out after %v", c.getConfig().RequestTimeout)
	}
}

// SendOk sends a request to the server and waits for any error returned by the server.
//
// The error returned can be a *client.UserError in case the server returned a error message.
func (c *Client) SendOk[RQ msgp.Marshaler](route string, req RQ) error {
	reqBytes, err := marshalRequestData(req)
	if err != nil {
		return err
	}

	respChan, reqId, err := c.sendRequest(route, reqBytes, true)
	if err != nil {
		return err
	}

	// Wait for time out duration for a response and remove chan after
	select {
	case res := <-respChan:
		c.removeResponseWaiter(reqId)

		if res.IsError {
			return NewUserError(string(res.Data))
		}

		return nil

	case <-time.After(c.getConfig().RequestTimeout):
		c.removeResponseWaiter(reqId)
		return fmt.Errorf("waiting for response timed out after %v", c.getConfig().RequestTimeout)
	}
}

// SendOkNoRequest sends to the server with no request and waits for any error returned by the server.
//
// The error returned can be a *client.UserError in case the server returned a error message.
func (c *Client) SendOkNoRequest(route string) error {
	respChan, reqId, err := c.sendRequest(route, []byte{}, true)
	if err != nil {
		return err
	}

	// Wait for time out duration for a response and remove chan after
	select {
	case res := <-respChan:
		c.removeResponseWaiter(reqId)

		if res.IsError {
			return NewUserError(string(res.Data))
		}

		return nil

	case <-time.After(c.getConfig().RequestTimeout):
		c.removeResponseWaiter(reqId)
		return fmt.Errorf("waiting for response timed out after %v", c.getConfig().RequestTimeout)
	}
}

// SendNoRequest sends a request to the server and waits for any error or response returned by the server.
//
// The error returned can be a *client.UserError in case the server returned a error message.
func (c *Client) SendNoRequest[RS any, RSP msgp.UnmarshalPtr[RS]](route string) (RS, error) {
	var resp RS

	respChan, reqId, err := c.sendRequest(route, []byte{}, true)
	if err != nil {
		return resp, err
	}

	// Wait for time out duration for a response and remove chan after
	select {
	case res := <-respChan:
		c.removeResponseWaiter(reqId)

		if res.IsError {
			return resp, NewUserError(string(res.Data))
		}

		resp, err = unmarshalResponseData[RS, RSP](res.Data)
		return resp, err

	case <-time.After(c.getConfig().RequestTimeout):
		c.removeResponseWaiter(reqId)
		return resp, fmt.Errorf("waiting for response timed out after %v", c.getConfig().RequestTimeout)
	}
}

// SendNoResponse sends a request to the server and does not wait for any response.
func (c *Client) SendNoResponse[RQ msgp.Marshaler](route string, req RQ) error {
	reqBytes, err := marshalRequestData[RQ](req)
	if err != nil {
		return err
	}

	_, _, err = c.sendRequest(route, reqBytes, false)
	return err
}

// SendPing sends to the server with not request and does not wait for any response.
func (c *Client) SendPing(route string) error {
	_, _, err := c.sendRequest(route, []byte{}, false)
	return err
}
