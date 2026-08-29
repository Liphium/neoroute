package neoroute

import (
	"bytes"
	"fmt"

	"github.com/tinylib/msgp/msgp"
)

// RouteTest can be used to simulate a session making a request to a route.
//
// This will execute middlewares, the route's handler itself and it's runAfter functions.
func RouteTest[D any, RQ msgp.Marshaler, RS any, PS msgp.UnmarshalPtr[RS],
](r *Router[D], session *Session[D], route string, req RQ) (RS, string, error) {
	r.Init()

	var respStruct RS

	respBytes, runAfterFuncs, err := handleRequestData(r, session, route, req)
	if err != nil {
		return respStruct, "", err
	}

	defer executeRunAfters(runAfterFuncs)

	resp, err := handleResponseData(respBytes)
	if err != nil {
		return respStruct, "", err
	}

	respData := responseData{
		HasData: resp.HasData,
		Data:    resp.Data,
		IsError: resp.IsError,
	}

	return GetTestingResponse[RS, PS](&respData)
}

// RouteNoRequestTest can be used to simulate a session making a request to a route.
//
// This will execute middlewares, the route's handler itself and it's runAfter functions.
func RouteNoRequestTest[D any, RS any, PS msgp.UnmarshalPtr[RS],
](r *Router[D], session *Session[D], route string) (RS, string, error) {
	r.Init()

	var respStruct RS

	respBytes, runAfterFuncs, err := makeRequest(r, session, route, []byte{})
	if err != nil {
		return respStruct, "", err
	}

	defer executeRunAfters(runAfterFuncs)

	resp, err := handleResponseData(respBytes)
	if err != nil {
		return respStruct, "", err
	}

	respData := responseData{
		HasData: resp.HasData,
		Data:    resp.Data,
		IsError: resp.IsError,
	}

	return GetTestingResponse[RS, PS](&respData)
}

// RouteOkTest can be used to simulate a session making a request to a route.
//
// This will execute middlewares, the route's handler itself and it's runAfter functions.
//
// The returned string is the user error if empty no user error was returned.
func RouteOkTest[D any, RQ msgp.Marshaler](r *Router[D], session *Session[D], route string, req RQ) (string, error) {
	r.Init()
	respBytes, runAfterFuncs, err := handleRequestData(r, session, route, req)
	if err != nil {
		return "", err
	}

	defer executeRunAfters(runAfterFuncs)

	resp, err := handleResponseData(respBytes)
	if err != nil {
		return "", err
	}

	respData := responseData{
		HasData: resp.HasData,
		Data:    resp.Data,
		IsError: resp.IsError,
	}

	return GetTestingResponseOk(&respData)
}

// RouteOkNoRequestTest can be used to simulate a session making a request to a route.
//
// This will execute middlewares, the route's handler itself and it's runAfter functions.
//
// The returned string is the user error if empty no user error was returned.
func RouteOkNoRequestTest[D any](r *Router[D], session *Session[D], route string) (string, error) {
	r.Init()

	respBytes, runAfterFuncs, err := makeRequest(r, session, route, []byte{})
	if err != nil {
		return "", err
	}

	defer executeRunAfters(runAfterFuncs)

	resp, err := handleResponseData(respBytes)
	if err != nil {
		return "", err
	}

	respData := responseData{
		HasData: resp.HasData,
		Data:    resp.Data,
		IsError: resp.IsError,
	}

	return GetTestingResponseOk(&respData)
}

// RouteNoResponseTest can be used to simulate a session making a request to a route.
//
// This will execute middlewares, the route's handler itself and it's runAfter functions.
func RouteNoResponseTest[D any, RQ msgp.Marshaler](r *Router[D], session *Session[D], route string, req RQ) error {
	r.Init()

	respBytes, runAfterFuncs, err := handleRequestData(r, session, route, req)
	if err != nil {
		return err
	}

	defer executeRunAfters(runAfterFuncs)

	if respBytes != nil {
		return fmt.Errorf("received response from route, make sure it's actually of type RouteNoResponse; responseBytes: %s", respBytes)
	}
	return nil
}

// RoutePingTest can be used to simulate a session making a request to a route.
//
// This will execute middlewares, the route's handler itself and it's runAfter functions.
func RoutePingTest[D any](r *Router[D], session *Session[D], route string) error {
	r.Init()

	respBytes, runAfterFuncs, err := makeRequest(r, session, route, []byte{})
	if err != nil {
		return err
	}

	defer executeRunAfters(runAfterFuncs)

	if respBytes != nil {
		return fmt.Errorf("received response from route, make sure it's actually of type RoutePing; responseBytes: %s", respBytes)
	}
	return nil
}

//
// Helper functions
//

// executeRunAfters executes the run after functions in order.
func executeRunAfters(runAfters []func()) {
	for _, runAfter := range runAfters {
		runAfter()
	}
}

// handleResponseData unmarshals the response data and returns it.
func handleResponseData(respBytes []byte) (response, error) {
	var resp response

	var msg message
	if _, err := msg.UnmarshalMsg(respBytes); err != nil {
		return resp, fmt.Errorf("failed to unmarshal message data: %v", err)
	}

	if _, err := resp.UnmarshalMsg(msg.Data); err != nil {
		return resp, fmt.Errorf("failed to unmarshal response data: %v", err)
	}
	return resp, nil
}

// handleRequestData marshals the request data and prepares it for the router.
func handleRequestData[D any, RQ msgp.Marshaler](r *Router[D], session *Session[D], route string, req RQ) ([]byte, []func(), error) {

	data, err := marshalRequestData(req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal request data: %v", err)
	}

	return makeRequest(r, session, route, data)
}

// makeRequest prepares and sends a request to the router, returning the response and any runAfter functions.
func makeRequest[D any](r *Router[D], session *Session[D], route string, data []byte) ([]byte, []func(), error) {
	mReq := request{
		Id:    1,
		Route: route,
		Data:  data,
	}

	reqBytes, err := mReq.MarshalMsg(nil)
	if err != nil {
		panic(fmt.Sprintf("failed to marshal request data for route %v: %v", route, err))
	}

	reader := bytes.NewReader(reqBytes)

	respBytes, runAfterFuncs := r.Handle(reader, session)

	return respBytes, runAfterFuncs, nil
}

// marshalRequestData marshals the request data using msgp.
func marshalRequestData[RQ msgp.Marshaler](req RQ) ([]byte, error) {
	reqData, err := req.MarshalMsg(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request data: %v", err)
	}
	return reqData, nil
}

// unmarshalResponseData unmarshals the response data using msgp.
func unmarshalResponseData[RS any, RSP msgp.UnmarshalPtr[RS]](respBytes []byte) (RS, error) {
	var data RS
	unmarshaler := any(&data).(msgp.Unmarshaler)
	_, err := unmarshaler.UnmarshalMsg(respBytes)
	if err != nil {
		return data, fmt.Errorf("failed to unmarshal response data: %v", err)
	}
	return data, nil
}
