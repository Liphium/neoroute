package config

import (
	"fmt"
	"reflect"

	"github.com/Liphium/neoroute/v2/neoschema"
)

// This is just here so we don't get into import cycles.

type DebugConfig struct {
	// The name of the transpoter you want to launch (in the schema).
	TransporterName string

	// The URL of the transporter endpoint (HTTP).
	TransporterURL string

	// Which HTTP method to use to connect to the transporter (e.g. POST, GET). Some transporters ignore this, but the HTTP transporter for example needs it.
	//
	// Default for HTTP: POST
	TransporterMethod string

	// The command you have to use to get the schema of the server.
	GenerateCommand string

	// Presets are a list of pre-filled requests you can send to the server. You will be able to select and edit them in the route selector.
	Presets []Preset
}

type Preset struct {
	// Name of the preset shown in neodebug.
	Name string

	// Actual route the data gets sent to.
	Route string

	// Data to be sent. For routes with no request, just leave this empty.
	Data any
}

// Check checks if the preset is valid and returns an error message if it is not. If it is valid, it returns an empty string.
func (p Preset) Check(routes map[string]neoschema.RouteSchema) string {
	if p.Route == "" {
		return "No route specified"
	}

	if p.Data == nil {
		return "No data specified"
	}

	route, exists := routes[p.Route]
	if exists {
		if !route.HasRequest {
			return ""
		}

		if route.Request == nil {
			return "Route does not accept a request body"
		}

		t, err := neoschema.BuildPackedFor(reflect.TypeOf(p.Data))
		if err != nil {
			return fmt.Sprintf("Schema build error: %v", err)
		}

		if t.Type() != route.Request.Type() {
			return fmt.Sprintf("Route accepts %v, but sending %v", route.Request.Type(), t.Type())
		}

		st, ok := t.(*neoschema.ReferenceType)
		st2, ok2 := route.Request.(*neoschema.ReferenceType)
		if ok && ok2 && st.Object != st2.Object {
			return fmt.Sprintf("Route accepts struct %v, but sending %v", st2.Object, st.Object)
		}
	} else {
		return "Route for this preset does not exist"
	}

	return ""
}

var Config DebugConfig
