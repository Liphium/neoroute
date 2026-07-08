package main

import (
	"flag"
	"fmt"
)

var (
	path    = flag.String("path", ".", "package for the golang package to generate a schema for")
	args    = flag.String("args", "--neo-generate", "arguments for generating schema")
	target  = flag.String("target", string("go"), "target language for generation")
	verbose = flag.Bool("v", false, "verbose diagnostics")
)

func main() {
	flag.Parse()

	switch *target {
	case "go":
		GenerateGo()
	default:
		fmt.Println("Unsupported target language: " + *target)
		fmt.Println(" ")
		fmt.Println("Try one of the following:")
		fmt.Println("- go")
	}
}
