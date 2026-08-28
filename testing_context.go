package neoroute

import (
	"errors"
	"fmt"

	"github.com/tinylib/msgp/msgp"
)

// GetTestingResponse return the response data or an error message for the user or
// an error from the handler or an error if the response is not correct.
func GetTestingResponse[RQ any, PQ msgp.UnmarshalPtr[RQ]](err error) (RQ, string, error) {
	var resp RQ
	if respData, ok := errors.AsType[*responseData](err); ok {
		if respData.IsError {
			return resp, string(respData.Data), nil
		}
		if !respData.HasData {
			return resp, "", fmt.Errorf("response has no data")
		}

		// Unmarshal response data into struct
		unmarshaler := any(&resp).(msgp.Unmarshaler)
		_, err := unmarshaler.UnmarshalMsg(respData.Data)
		if err != nil {
			return resp, "", fmt.Errorf("failed to unmarshal response data: %v", err)
		}
		return resp, "", nil
	} else {

		// The handler returned an error, this would not be returned,
		// but provided to the error handler and the response from it would be returned.
		return resp, "", err
	}
}

// GetTestingResponseOk return the error message or an error from the handler or
// an error from the handler or an error if the response is not correct.
func GetTestingResponseOk(err error) (string, error) {

	if respData, ok := errors.AsType[*responseData](err); ok {
		if respData.IsError {
			return string(respData.Data), nil
		}
		if respData.HasData {
			return "", fmt.Errorf("ok response without error should not have any data")
		}
		return "", nil
	} else {
		// The handler returned an error, this would not be returned,
		// but provided to the error handler and the response from it would be returned.
		return "", err
	}
}
