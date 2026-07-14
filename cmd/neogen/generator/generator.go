package generator

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Liphium/neoroute/cmd/neogen/languages"
	"github.com/Liphium/neoroute/neoschema"
)

type GeneratorConfig struct {
	ServerPath     string
	Command        string
	TargetLanguage string
	Verbose        bool
}

var Config GeneratorConfig

func Generate(config GeneratorConfig) {
	Config = config

	// Find the server and run it
	cmd := exec.Command(strings.Split(Config.Command, " ")[0], strings.Split(Config.Command, " ")[1:]...)
	var err error
	cmd.Dir, err = filepath.Abs(Config.ServerPath)
	if err != nil {
		panic(fmt.Errorf("couldn't get absolute path of server: %v", err))
	}

	bytes, err := cmd.Output()
	if err != nil {
		panic(fmt.Errorf("couldn't run app: %v", err))
	}

	var schema neoschema.Schema
	if err := json.Unmarshal(bytes, &schema); err != nil {
		panic(fmt.Errorf("invalid schema: %v", err))
	}

	switch config.TargetLanguage {
	case "go":
		GenerateGo(schema)
	case "typescript":
		GenerateWithConfig(schema, languages.NewTSConfig())
	default:
		fmt.Println("Unsupported target language: " + config.TargetLanguage)
		fmt.Println(" ")
		fmt.Println("Try one of the following:")
		fmt.Println("- go")
	}
}
