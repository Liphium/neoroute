package neoschema

import (
	"reflect"

	"github.com/Liphium/neoroute"
)

// The definition for one request + response endpoint.
//
// Contains everything neoroute can currently return + needs.
type RequestResponse struct {
	HasRequest bool
	Request    reflect.Type

	HasResponse    bool
	CanReturnError bool
	Response       reflect.Type
}

type Transporter interface {
	// Should return the transporter type (defined in the neogen package).
	Type() int

	// Should return the request response schemas for this transporter.
	GetSchema() map[string]RequestResponse

	// Should return the event registries for this transporter (can also be none).
	GetRegistries() []*neoroute.EventRegistry
}

// Convert a route map returned by a router into a request response schema as used by neogen.
func ToRouteSchema[D any](routeMap map[string]neoroute.RouteData[D]) map[string]RequestResponse {
	schema := map[string]RequestResponse{}
	for route, routeData := range routeMap {
		hasRequest, request := routeData.RequestType()
		canReturnError, hasResponse, response := routeData.ResponseType()

		schema[route] = RequestResponse{
			HasRequest:     hasRequest,
			Request:        request,
			HasResponse:    hasResponse,
			CanReturnError: canReturnError,
			Response:       response,
		}
	}

	return schema
}
