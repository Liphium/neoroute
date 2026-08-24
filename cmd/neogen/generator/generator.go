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
	schema, err := GetSchema(config.ServerPath, config.Command)
	if err != nil {
		panic(err)
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
		fmt.Println("- typescript")
	}
}

// GetSchema gets the schema of any neoroute server supporting generation of the schema (or a custom command).
func GetSchema(serverPath string, command string) (neoschema.Schema, error) {
	var schema neoschema.Schema

	// Find the server and run it
	cmd := exec.Command(strings.Split(command, " ")[0], strings.Split(command, " ")[1:]...)
	var err error
	cmd.Dir, err = filepath.Abs(serverPath)
	if err != nil {
		return schema, fmt.Errorf("couldn't get absolute path of server: %v", err)
	}

	bytes, err := cmd.Output()
	if err != nil {
		return schema, fmt.Errorf("couldn't run app: %v", err)
	}

	if err := json.Unmarshal(bytes, &schema); err != nil {
		return schema, fmt.Errorf("invalid schema: %v", err)
	}

	return schema, nil
}
