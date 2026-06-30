package go_gen

import (
	"fmt"
	"strings"

	"github.com/Liphium/neoroute/cmd/neogen/util"
	"github.com/Liphium/neoroute/neoschema"
)

func GenerationLine(schema neoschema.Schema) string {
	return fmt.Sprintf("// Code generated with %s-generated v%d schema by neogen. DO NOT EDIT.", schema.Generator, schema.Version)
}

func GenerateTransporters(schema neoschema.Schema) (map[string]string, error) {
	transporterFiles := map[string]string{}
	for name, transporter := range schema.Transporters {
		var err error
		var generated string

		switch transporter.Type {
		case neoschema.TransporterHTTP:
			generated, err = GenerateHTTPTransporter(name, GenerationLine(schema), transporter)
			if err != nil {
				return transporterFiles, fmt.Errorf("Couldn't generate HTTP transporter %s: %v", name, err)
			}

		case neoschema.TransporterWebSocket:
			generated, err = GenerateWebSocketTransporter(name, GenerationLine(schema), transporter)
			if err != nil {
				return transporterFiles, fmt.Errorf("Couldn't generate WebSocket transporter %s: %v", name, err)
			}

		case neoschema.TransporterWebTransport:
			return transporterFiles, fmt.Errorf("WebTransport is not yet supported for Go")
		}

		transporterFiles["connector_"+strings.ToLower(name)+".go"] = generated
	}

	return transporterFiles, nil
}

const eventHandler = `func (c *%s) Receive%s(handler func(event %s)) {
	fmt.Println("Handling some event!")
}`

func GenerateEvent(transporterName, event string, packed neoschema.PackedType) (string, error) {
	eventType, err := GetType("", packed)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(eventHandler, transporterName, util.ToCamelCase(event, true), eventType), nil
}

const routeCaller = `func (c *%s) Send%s(%s) %s {
	fmt.Println("Sending some event!")
	// TODO: Return
}`

func GenerateRoutes(transporterName, name string, schema neoschema.RouteSchema) (string, error) {
	var requestType, responseType string = "", ""
	var err error
	if schema.HasRequest {
		requestType, err = GetType("", schema.Request)
		if err != nil {
			return "", fmt.Errorf("couldn't generate request type: %v", err)
		}
		requestType = "payload " + requestType
	}

	if schema.HasResponse {
		responseType, err = GetType("", schema.Response)
		if err != nil {
			return "", fmt.Errorf("couldn't generate response type: %v", err)
		}
		responseType = "(" + responseType + ", error)"
	}
	if schema.CanReturnError && responseType == "" {
		responseType = "error"
	}

	return fmt.Sprintf(routeCaller, transporterName, util.ToCamelCase(name, true), requestType, responseType), nil
}
