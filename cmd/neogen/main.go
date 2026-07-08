package main

import (
	"flag"

	"github.com/Liphium/neoroute/cmd/neogen/generator"
)

var (
	path    = flag.String("path", ".", "package for the golang package to generate a schema for")
	args    = flag.String("args", "--neo-generate", "arguments for generating schema")
	target  = flag.String("target", string("go"), "target language for generation")
	verbose = flag.Bool("v", false, "verbose diagnostics")
)

func main() {
	flag.Parse()

	generator.Generate(generator.GeneratorConfig{
		ServerPath:        *path,
		ArgsForGeneration: *args,
		TargetLanguage:    *target,
		Verbose:           *verbose,
	})
}
