package main

import (
	"flag"

	"github.com/Liphium/neoroute/cmd/neogen/generator"
)

var (
	path    = flag.String("path", ".", "package for the golang package to generate a schema for")
	command = flag.String("command", "go run . --neo-generate", "command for generating the schema")
	target  = flag.String("target", string("go"), "target language for generation")
	verbose = flag.Bool("v", false, "verbose diagnostics")
)

func main() {
	flag.Parse()

	generator.Generate(generator.GeneratorConfig{
		ServerPath:     *path,
		Command:        *command,
		TargetLanguage: *target,
		Verbose:        *verbose,
	})
}
