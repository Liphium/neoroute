package neoschema

import "reflect"

type Generator struct {
	transporters  map[string]Transporter
	customObjects []reflect.Type
}

// Create a new generator.
//
// Can generate schemas for your transporters to make your life easier.
//
// You'll of course need the other part of the generator as well, this just generates a json file that can be parsed and used to generate bindings using neogen or other tools.
func NewGenerator() *Generator {
	return &Generator{
		transporters: map[string]Transporter{},
	}
}

// Add adds custom objects that will be generated alongside the normal models for transporters, this is useful for custom types that are not directly linked to a transporter, but you still use in other shared contexts.
func (g *Generator) Add(objects ...any) *Generator {
	for _, o := range objects {
		g.customObjects = append(g.customObjects, reflect.TypeOf(o))
	}
	return g
}

// Add a new transporter, needs to implement the interface for schema generation of course...
func (g *Generator) Transporter(name string, schema Transporter) {
	g.transporters[name] = schema
}
