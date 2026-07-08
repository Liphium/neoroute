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
	eventType, err := GetType(packed)
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

var routeCaller = template.Must(template.New("").Parse(`func (c *{{ .transporterName }}) Send{{ .routeCamelCase }}({{ .requestType }}) {{ .responseType }} {
	{{ if eq .sendType .const.SendRequestResponse }}
	return client.Send[{{ .responseStruct }}]({{ .receiverName }}, "{{ .route }}", payload)
	{{ else if eq .sendType .const.SendOK }}
	return client.SendOk({{ .receiverName }}, "{{ .route }}", payload)
	{{ else if eq .sendType .const.SendOKNoRequest }}
	return client.SendOkNoRequest({{ .receiverName }}, "{{ .route }}")
	{{ else if eq .sendType .const.SendNoRequest }}
	return client.SendNoRequest[{{ .responseStruct }}]({{ .receiverName }}, "{{ .route }}")
	{{ else if eq .sendType .const.SendNoResponse }}
	return client.SendNoResponse({{ .receiverName }}, "{{ .route }}", payload)
	{{ else if eq .sendType .const.SendSignal }}
	return client.SendPing({{ .receiverName }}, "{{ .route }}")
	{{ end }}
}`))

func GenerateRoutes(transporterName, receiverName, name string, schema neoschema.RouteSchema) (string, error) {
	var requestType, requestStruct, responseType, responseStruct string = "", "", "error", ""
	var err error

	if schema.HasRequest {
		requestStruct, err = GetType(schema.Request)
		if err != nil {
			return "", fmt.Errorf("couldn't generate request type: %v", err)
		}
		requestType = "payload " + requestStruct
	}

	if schema.HasResponse {
		responseStruct, err = GetType(schema.Response)
		if err != nil {
			return "", fmt.Errorf("couldn't generate response type: %v", err)
		}
		responseType = "(" + responseStruct + ", error)"
	}

	builder := strings.Builder{}
	routeCaller.Execute(&builder, map[string]any{
		"receiverName":    receiverName,
		"transporterName": transporterName,
		"route":           name,
		"routeCamelCase":  util.ToCamelCase(name, true),
		"requestType":     requestType,
		"requestStruct":   requestStruct,
		"responseType":    responseType,
		"responseStruct":  responseStruct,
		"sendType":        schema.GetSendType(),
		"const":           neoschema.SendTypeMap(),
	})
	return builder.String(), nil
}
