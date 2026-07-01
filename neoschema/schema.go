package neoschema

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

const (
	CurrentVersion int    = 1
	GeneratorName  string = "neogen"
)

type Schema struct {
	Version      int                          `json:"version"`
	Generator    string                       `json:"generator"`
	Transporters map[string]TransporterSchema `json:"transporters"`
}

const (
	TransporterHTTP = iota
	TransporterWebTransport
	TransporterWebSocket
)

type TransporterSchema struct {
	Type   int                    `json:"type"`
	Events map[string]PackedType  `json:"events"`
	Routes map[string]RouteSchema `json:"routes"`
}

type RouteSchema struct {
	HasRequest bool       `json:"has_request"`
	Request    PackedType `json:"request"`

	HasResponse    bool       `json:"has_response"`
	CanReturnError bool       `json:"can_return_error"`
	Response       PackedType `json:"response"`
}

type SendType int

const (
	SendRequestResponse SendType = iota // Has a typed request and response, so the client can send a request and get back a typed response or an error
	SendOK                              // Typed request with the possibility to get back an error, but no typed response
	SendOKNoRequest                     // Same as SendNoRequest, but can only get back an error
	SendNoRequest                       // Has no typed request, just sends a signal and gets back a typed response or an error
	SendNoResponse                      // Has a typed request but no response, so nothing will be confirming this one (can be used for sending movement stuff, etc.)
	SendSignal                          // Has no typed request and response, as well as no error (rarely used)
)

// SendTypeMap gives you a map of strings ot the relevant SendType, this can be useful for text/template templates.
func SendTypeMap() map[string]SendType {
	return map[string]SendType{
		"SendRequestResponse": SendRequestResponse,
		"SendOK":              SendOK,
		"SendOKNoRequest":     SendOKNoRequest,
		"SendNoRequest":       SendNoRequest,
		"SendNoResponse":      SendNoResponse,
		"SendSignal":          SendSignal,
	}
}

// GetSendType gives you the send type for a function, the 8 possibilties of the schema translated to the 6 functions we usually have for neoroute in any SDK.
func (rs RouteSchema) GetSendType() SendType {
	switch {
	case rs.HasRequest && rs.HasResponse:
		return SendRequestResponse
	case rs.HasRequest && !rs.HasResponse && rs.CanReturnError:
		return SendOK
	case !rs.HasRequest && !rs.HasResponse && rs.CanReturnError:
		return SendOKNoRequest
	case !rs.HasRequest && rs.HasResponse:
		return SendNoRequest
	case rs.HasRequest && !rs.HasResponse:
		return SendNoResponse
	default:
		return SendSignal
	}
}

func (g *Generator) Generate() (Schema, error) {
	var err error

	packedTransporters := map[string]TransporterSchema{}
	for name, transporter := range g.transporters {

		// Generate the BasicType for every event by name
		packedEvents := map[string]PackedType{}
		for _, reg := range transporter.GetRegistries() {
			events := reg.GetEvents()
			schemas := reg.GetSchemas()

			for i, event := range events {
				packedEvents[event], err = BuildPackedFor(schemas[i]())
				if err != nil {
					return Schema{}, err
				}
			}
		}

		// Generate the RouteSchema from all of definitions for the endpoints
		packedRoutes := map[string]RouteSchema{}
		for route, routeData := range transporter.GetSchema() {
			var request, response PackedType
			if routeData.HasRequest {
				request, err = BuildPackedFor(routeData.Request)
				if err != nil {
					return Schema{}, err
				}
			}
			if routeData.HasResponse {
				response, err = BuildPackedFor(routeData.Response)
				if err != nil {
					return Schema{}, err
				}
			}

			packedRoutes[route] = RouteSchema{
				HasRequest:     routeData.HasRequest,
				Request:        request,
				HasResponse:    routeData.HasResponse,
				CanReturnError: routeData.CanReturnError,
				Response:       response,
			}
		}

		packedTransporters[name] = TransporterSchema{
			Type:   transporter.Type(),
			Events: packedEvents,
			Routes: packedRoutes,
		}
	}

	return Schema{
		Version:      CurrentVersion,
		Generator:    GeneratorName,
		Transporters: packedTransporters,
	}, nil
}

// Will parse command line arguments to see if --neo-generate is set.
//
// If it is, it will print the schema and exit. If not, execution will resume like normal. This function is meant as a shorthand to easily make your program use neoroute's generation without having to introduce weird stuff yourself and having a standard across all neoroute projects.
func (g *Generator) PrintOrPanic() {

	// Check for the flag in the execution arguments
	found := false
	for _, arg := range os.Args {
		if strings.TrimSpace(arg) == "--neo-generate" {
			found = true
		}
	}
	if !found {
		return
	}

	// Actually generate the schema
	s, err := g.Generate()
	if err != nil {
		panic(fmt.Sprintf("Couldn't generate the schema: %v", err))
	}

	// Print pretty json
	marshaled, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		panic(fmt.Sprintf("Couldn't marshal the schema: %v", err))
	}

	fmt.Println(string(marshaled))
	os.Exit(0)
}
