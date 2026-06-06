package neoroute

import (
	"errors"
	"fmt"
	"sync"

	"github.com/tinylib/msgp/msgp"
)

func NewTestingResCtx[D any, RS any, PS interface {
	*RS
	msgp.Marshaler
}](neo *NeoRouter[D], route string, session *Session[D]) *ResCtx[D, RS, PS] {
	return &ResCtx[D, RS, PS]{
		Ctx: NewTestingCtx[D](neo, route, session),
	}
}

func NewTestingOkCtx[D any](neo *NeoRouter[D], route string, session *Session[D]) *OkCtx[D] {
	return &OkCtx[D]{
		Ctx: NewTestingCtx[D](neo, route, session),
	}
}

func NewTestingCtx[D any](neo *NeoRouter[D], route string, session *Session[D]) *Ctx[D] {
	return &Ctx[D]{
		neo:     neo,
		id:      1,
		route:   route,
		session: session,
	}
}

func NewTestingSession[D any](data D, sessionId string) *Session[D] {
	return &Session[D]{
		sessionData: data,
		id:          sessionId,
		mutex:       &sync.Mutex{},
	}
}

func EvaluateCtxTesting[D any](c *Ctx[D]) {
	defer func() {
		for _, fn := range c.runAfter {
			fn()
		}
	}()
}

func GetTestingResponse[RQ any, PQ interface {
	*RQ
	msgp.Unmarshaler
}](err error) (RQ, string, error) {
	var respErr response
	var resp RQ
	if errors.As(err, &respErr) {
		if respErr.IsError {
			return resp, string(respErr.Data), nil
		}
		if !respErr.HasData {
			return resp, "", fmt.Errorf("response has no data")
		}

		// Unmarshal response data into struct
		unmarshaler := any(&resp).(msgp.Unmarshaler)
		_, err := unmarshaler.UnmarshalMsg(respErr.Data)
		if err != nil {
			return resp, "", fmt.Errorf("failed to unmarshal response data: %v", err)
		}
		return resp, "", nil
	} else {
		return resp, "", fmt.Errorf("response is not response type")
	}
}

// GetTestingResponseOk return the error message or an error from the handler or an error if the response is not correct.
func GetTestingResponseOk(err error) (string, error) {
	var respErr response

	if errors.As(err, &respErr) {
		if respErr.IsError {
			return string(respErr.Data), nil
		}
		if respErr.HasData {
			return "", fmt.Errorf("ok response without error should not have any data")
		}
		return "", nil
	} else {
		return "", fmt.Errorf("response is not response type")
	}
}
