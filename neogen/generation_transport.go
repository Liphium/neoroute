package neogen

import "reflect"

type RequestResponse struct {
	HasRequest bool
	Request    reflect.Type

	HasResponse    bool
	CanReturnError bool
	Response       reflect.Type
}

type RequestResponseSchema interface {
	GetSchema() map[string]RequestResponse
}
