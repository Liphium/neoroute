package client

import (
	"fmt"

	"github.com/tinylib/msgp/msgp"
)

func marshalRequestData[RQ any, RQP interface {
	*RQ
	msgp.Marshaler
}](req RQ) ([]byte, error) {
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
}](r Sender, respBytes []byte) (RS, error) {
	var data RS
	unmarshaler := any(&data).(msgp.Unmarshaler)
	_, err := unmarshaler.UnmarshalMsg(respBytes)
	if err != nil {
		return data, fmt.Errorf("failed to unmarshal response data: %v", err)
	}
	return data, nil
}
