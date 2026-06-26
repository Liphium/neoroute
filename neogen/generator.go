package neogen

import (
	"github.com/Liphium/neoroute"
)

type Generator struct {
	transporters map[string]RequestResponseSchema
	registries   map[string][]neoroute.EventRegistry
}

// Create a new generator.
//
// Can generate schemas for your transporters to make your life easier.
//
// You'll of course need the other part of the generator is well, this just generates a json file that can be parsed and used to generate bindings that can be used by client libraries's generators.
func NewGenerator() *Generator {
	return &Generator{
		transporters: map[string]RequestResponseSchema{},
		registries:   map[string][]neoroute.EventRegistry{},
	}
}

// Add a new transporter, any transporter probably implements a schema.
func (g *Generator) Transporter(name string, schema RequestResponseSchema) {
	g.transporters[name] = schema
}

// Add a new event registry for a transporter. Will make sure the events are generated into different transporters.
func (g *Generator) RegistryFor(name string, registries ...neoroute.EventRegistry) {
	g.registries[name] = append(g.registries[name], registries...)
}
