// Package testschema provides shared fixtures for generated-language tests.
package testschema

import (
	"github.com/Liphium/neoroute/v2/neoschema"
	"reflect"
	"time"
)

type Point struct {
	neoschema.IsTuple `msg:"-"`
	X                 float64 `msg:"x"`
	Label             string  `msg:"label"`
}
type Envelope struct {
	Point  Point             `msg:"point"`
	Items  []*Point          `msg:"items"`
	Matrix [][]Point         `msg:"matrix"`
	Lookup map[string]*Point `msg:"lookup"`
	Next   *Envelope         `msg:"next"`
	Bytes  []byte            `msg:"bytes"`
	Date   time.Time         `msg:"date"`
	Opaque any               `msg:"opaque"`
}

type Address struct {
	City string `msg:"city"`
	Zip  *int   `msg:"zip"`
}

type Inner struct {
	Note string `msg:"note"`
}

type User struct {
	Name    string   `msg:"name"`
	Addr    *Address `msg:"addr"`
	Tags    []string `msg:"tags"`
	Inner   Inner    `msg:"inner"`
	Friends []User   `msg:"friends"`
}

// StandaloneInner is not referenced by any route or event, it must still be generated because it is inside of Standalone.
type StandaloneInner struct {
	Info string `msg:"info"`
}

// Standalone is not referenced by any route or event, it must still be generated.
type Standalone struct {
	Code      int               `msg:"code"`
	Label     *string           `msg:"label"`
	Items     []StandaloneInner `msg:"items"`
	SomeStamp time.Time         `msg:"some_stamp"`

	// These fields should not be generated
	secret string
	hello  int64 `msg:"hello"`
	sth    time.Time
}

type Coordinates struct {
	neoschema.IsTuple `msg:"-"`
	X                 float64 `msg:"x"`
	Y                 float64 `msg:"y"`
}
type Location struct {
	Name        string      `msg:"name"`
	Coordinates Coordinates `msg:"coordinates"`
}

func Schema() (neoschema.Schema, error) {
	g := neoschema.NewGenerator()
	g.Add(Standalone{})
	g.Add(Location{})
	schema, err := g.Generate()
	if err != nil {
		return schema, err
	}
	packed, err := neoschema.BuildPackedFor(reflect.TypeFor[Envelope]())
	if err != nil {
		return schema, err
	}
	user, err := neoschema.BuildPackedFor(reflect.TypeFor[User]())
	if err != nil {
		return schema, err
	}
	routes := map[string]neoschema.RouteSchema{
		"getUser":  {HasRequest: true, Request: user, HasResponse: true, Response: user},
		"exchange": {HasRequest: true, Request: packed, HasResponse: true, Response: packed},
		"get":      {HasResponse: true, Response: packed},
		"ok":       {HasRequest: true, Request: packed, CanReturnError: true},
		"signal":   {CanReturnError: true},
		"push":     {HasRequest: true, Request: packed},
		"ping":     {},
	}
	schema.Transporters = map[string]neoschema.TransporterSchema{
		"http": {Type: neoschema.TransporterHTTP, Routes: routes, Events: map[string]neoschema.PackedType{"update": packed}},
		"ws":   {Type: neoschema.TransporterWebSocket, Routes: routes},
	}
	return schema, nil
}
