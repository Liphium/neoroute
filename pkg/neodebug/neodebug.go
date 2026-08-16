package neodebug

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/Liphium/neoroute/cmd/neogen/generator"
	"github.com/Liphium/neoroute/pkg/neodebug/config"
	"github.com/Liphium/neoroute/pkg/neodebug/tui"
)

func Run(cfg config.DebugConfig) {
	config.Config = cfg

	// Make sure we're in the project directory
	_, err := os.ReadFile("go.mod")
	if err != nil {
		panic("please run neodebug in the project directory")
	}

	// Generate the schema
	cmd := "go run . --neo-generate"
	if cfg.GenerateCommand != "" {
		cmd = cfg.GenerateCommand
	}
	schema, err := generator.GetSchema(".", cmd)
	if err != nil {
		panic(fmt.Errorf("schema generation failed: %w", err))
	}

	// Find the transporter we want to connect to
	transporter, ok := schema.Transporters[cfg.TransporterName]
	if !ok {
		panic(fmt.Errorf("couldn't find transporter %s", cfg.TransporterName))
	}

	// Start the actual TUI with the correct transporter
	p := tea.NewProgram(tui.Run(transporter))
	if _, err := p.Run(); err != nil {
		fmt.Println("neodebug crashed:", err)
		os.Exit(1)
	}
}
