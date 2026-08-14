package neodebug

import (
	"fmt"
	"os"

	"github.com/Liphium/neoroute/cmd/neogen/generator"
	"github.com/Liphium/neoroute/neoschema"
)

type DebugConfig struct {
	// The name of the transpoter you want to launch (in the schema).
	TransporterName string

	// The URL of the transporter endpoint (HTTP).
	TransporterURL string

	// The command you have to use to get the schema of the server.
	GenerateCommand string
}

var Config DebugConfig

func Run(config DebugConfig) {
	Config = config

	// Make sure we're in the project directory
	_, err := os.ReadFile("go.mod")
	if err != nil {
		panic("please run neodebug in the project directory")
	}

	// Generate the schema
	cmd := "go run . --neo-generate"
	if config.GenerateCommand != "" {
		cmd = config.GenerateCommand
	}
	schema, err := generator.GetSchema(".", cmd)
	if err != nil {
		panic(fmt.Errorf("schema generation failed: %w", err))
	}

	// Find the transporter we want to connect to
	transporter, ok := schema.Transporters[config.TransporterName]
	if !ok {
		panic(fmt.Errorf("couldn't find transporter %s", config.TransporterName))
	}

	run(transporter)
}

// run opens the debug UI.
func run(transporter neoschema.TransporterSchema) {
	// TODO: Open the neodebug UI
}
