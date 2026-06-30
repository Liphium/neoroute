package go_gen

import (
	"fmt"
	"strings"
	"text/template"

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

var eventHandler = template.Must(template.New("").Parse(`func (c *{{ .transporterName }}) Receive{{ .eventCamelCase }}(handler func(event {{ .eventType }})) {
	client.Receive[{{ .eventType }}, *{{ .eventType }}](c.{{ .receiver }}, "{{ .event }}", func(c *client.Ctx, event {{ .eventType }}) {
		handler(event)
	})
}`))

func GenerateEvent(transporterName, receiverName, event string, packed neoschema.PackedType) (string, error) {
	eventType, err := GetType("", packed)
	if err != nil {
		return "", err
	}

	result := strings.Builder{}
	eventHandler.Execute(&result, map[string]string{
		"transporterName": transporterName,
		"receiver":        receiverName,
		"eventCamelCase":  util.ToCamelCase(event, true),
		"eventType":       eventType,
		"event":           event,
	})
	return result.String(), nil
}

var routeCaller = template.Must(template.New("").Parse(`func (c *{{ .transporterName }}) Send{{ .routeName }}({{ .requestType }}) /* {{ .responseType }} */ {
	// TODO: Implement
}`))

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

	builder := strings.Builder{}
	routeCaller.Execute(&builder, map[string]string{
		"transporterName": transporterName,
		"routeName":       util.ToCamelCase(name, true),
		"requestType":     requestType,
		"responseType":    responseType,
	})
	return builder.String(), nil
}
