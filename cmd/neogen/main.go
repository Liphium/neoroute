package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/Liphium/neoroute/cmd/neogen/generator"
)

var (
	version = flag.Bool("version", false, "show version")
	path    = flag.String("path", ".", "package for the golang package to generate a schema for")
	command = flag.String("command", "go run . --neo-generate", "command for generating the schema")
	target  = flag.String("target", string("go"), "target language for generation")
	verbose = flag.Bool("v", false, "verbose diagnostics")
)

func main() {
	flag.Parse()

	// Cause I want --version and not -version
	if version != nil && *version {
		info, ok := debug.ReadBuildInfo()
		if !ok {
			fmt.Println("neogen dev build")
		} else {
			fmt.Printf("neogen %s  \n", info.Main.Version)
		}
		os.Exit(0)
	}

	generator.Generate(generator.GeneratorConfig{
		ServerPath:     *path,
		Command:        *command,
		TargetLanguage: *target,
		Verbose:        *verbose,
	})
}
