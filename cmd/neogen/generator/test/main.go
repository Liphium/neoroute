// This file is for testing only and embedded in tests.
package main

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/Liphium/neoroute"
	"github.com/Liphium/neoroute/neoschema"
)

type GenericData struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type Recursive struct {
	Next *Recursive `json:"next"`
}

type Normal struct {
	Value int    `json:"value"`
	Text  string `json:"text"`
}

type MockTransporter struct {
	tType      neoschema.TransporterType
	routes     map[string]neoschema.RequestResponse
	registries []*neoroute.EventRegistry
}

func (m *MockTransporter) Type() neoschema.TransporterType { return m.tType }
func (m *MockTransporter) GetSchema() map[string]neoschema.RequestResponse {
	return m.routes
}
func (m *MockTransporter) GetRegistries() []neoroute.IEventRegistry {
	return nil
}

type CustomTransporter struct {
	MockTransporter
	registries []neoroute.IEventRegistry
}

func (c *CustomTransporter) GetRegistries() []neoroute.IEventRegistry {
	return c.registries
}

type MockEventRegistry struct {
	events  []string
	schemas []func() reflect.Type
}

func (m *MockEventRegistry) GetEvents() []string { return m.events }
func (m *MockEventRegistry) GetSchemas() []func() reflect.Type {
	return m.schemas
}

func main() {
	gen := neoschema.NewGenerator()

	data := reflect.TypeOf(GenericData{})

	routes := map[string]neoschema.RequestResponse{
		"requestResponse": {
			HasRequest:     true,
			Request:        data,
			HasResponse:    true,
			CanReturnError: true,
			Response:       data,
		},
		"ok": {
			HasRequest:     true,
			Request:        data,
			CanReturnError: true,
		},
		"okNoRequest": {
			CanReturnError: true,
		},
		"noRequest": {
			HasResponse:    true,
			CanReturnError: true,
			Response:       data,
		},
		"noResponse": {
			HasRequest: true,
			Request:    data,
		},
		"signal": {},
	}

	gen.Transporter("http", &MockTransporter{
		tType:  neoschema.TransporterHTTP,
		routes: routes,
	})

	reg := &MockEventRegistry{
		events: []string{"recursive", "normal"},
		schemas: []func() reflect.Type{
			func() reflect.Type { return reflect.TypeOf(Recursive{}) },
			func() reflect.Type { return reflect.TypeOf(Normal{}) },
		},
	}

	gen.Transporter("websocket", &CustomTransporter{
		MockTransporter: MockTransporter{
			tType:  neoschema.TransporterWebSocket,
			routes: routes,
		},
		registries: []neoroute.IEventRegistry{reg},
	})

	schema, err := gen.Generate()
	if err != nil {
		panic(err)
	}

	out, _ := json.MarshalIndent(schema, "", "  ")
	fmt.Println(string(out))
}
