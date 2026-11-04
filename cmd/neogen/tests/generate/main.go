// Generate the shared fixtures for inspecting and running language tests locally.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Liphium/neoroute/v2/cmd/neogen/engine"
	"github.com/Liphium/neoroute/v2/cmd/neogen/internal/testschema"
	"github.com/Liphium/neoroute/v2/cmd/neogen/languages"
)

func main() {
	language := flag.String("language", "", "go or typescript")
	output := flag.String("output", "", "directory for generated files")
	flag.Parse()
	if err := generate(*language, *output); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func generate(language, output string) error {
	if output == "" {
		return fmt.Errorf("-output is required")
	}
	schema, err := testschema.Schema()
	if err != nil {
		return err
	}
	var config engine.LanguageConfig
	switch language {
	case "typescript":
		config = languages.NewTSConfig()
	case "go":
		for _, transport := range schema.Transporters {
			delete(transport.Routes, "collection")
			delete(transport.Routes, "lookup")
		}
		if err := os.Setenv("GOPACKAGE", "generated"); err != nil {
			return err
		}
		_, config = languages.NewGoConfig()
	default:
		return fmt.Errorf("unsupported language %q", language)
	}
	files, err := engine.NewGenerationEngine(config).Generate(schema)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(output, 0755); err != nil {
		return err
	}
	for name, source := range files {
		if err := os.WriteFile(filepath.Join(output, name), []byte(source), 0644); err != nil {
			return err
		}
	}
	return nil
}
