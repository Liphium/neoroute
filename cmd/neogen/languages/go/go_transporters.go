package go_gen

import (
	"fmt"
	"os"
	"strings"

	"github.com/Liphium/neoroute/cmd/neogen/util"
	"github.com/Liphium/neoroute/neoschema"
)

func GenerationLine(schema neoschema.Schema) string {
	return fmt.Sprintf("// Code generated with %s-generated v%d schema by neogen. DO NOT EDIT.", schema.Generator, schema.Version)
}

const transporterStart = `%s
package %s

import "fmt"

type %s struct{}

func New%s() *%s {
	return &%s{}
}

func (c *%s) Connect() {
	fmt.Println("Hello, neogen!")
}

`

func GenerateTransporters(schema neoschema.Schema) (map[string]string, error) {
	transporterFiles := map[string]string{}
	for name, transporter := range schema.Transporters {
		transporterName := util.ToCamelCase(name+".Connector", true)
		file := fmt.Sprintf(transporterStart, GenerationLine(schema), os.Getenv("GOPACKAGE"), transporterName, transporterName, transporterName, transporterName, transporterName)

		// Generate the stuff for all the events
		for event, packed := range transporter.Events {
			generated, err := GenerateEvent(transporterName, event, packed)
			if err != nil {
				return transporterFiles, fmt.Errorf("Couldn't generate event %s: %v", name, err)
			}

			file += generated + "\n\n"
		}

		// Generate the stuff for all route schemas
		for name, schema := range transporter.Routes {
			generated, err := GenerateRoutes(transporterName, name, schema)
			if err != nil {
				return transporterFiles, fmt.Errorf("Couldn't generate route %s: %v", name, err)
			}

			file += generated + "\n\n"
		}

		transporterFiles["connector_"+strings.ToLower(name)+".go"] = file
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
