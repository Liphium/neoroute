package neoschema

type Generator struct {
	transporters map[string]Transporter
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

// Add a new transporter, needs to implement the interface for schema generation of course...
func (g *Generator) Transporter(name string, schema Transporter) {
	g.transporters[name] = schema
}
