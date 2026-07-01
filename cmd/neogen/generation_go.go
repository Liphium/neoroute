package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Liphium/neoroute/cmd/neogen/languages"
	go_gen "github.com/Liphium/neoroute/cmd/neogen/languages/go"
	"github.com/Liphium/neoroute/neoschema"
)

func GenerateGo() {

	// Run the other thingy
	cmd := exec.Command("go", append([]string{"run", "."}, strings.Split(*args, " ")...)...)
	wd, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("couldn't get working directory: %v", err))
	}
	cmd.Dir = filepath.Clean(filepath.Join(wd, *path))

	bytes, err := cmd.Output()
	if err != nil {
		panic(fmt.Errorf("couldn't run app: %v", err))
	}

	var schema neoschema.Schema
	if err := json.Unmarshal(bytes, &schema); err != nil {
		panic(fmt.Errorf("invalid schema: %v", err))
	}

	// Generate object and transporters
	goObjectFile(schema)
	goTransporters(schema)

	// Run the go formatter to make sure it's all nice and clean
	cmd = exec.Command("go", "fmt", ".")
	if err := cmd.Run(); err != nil {
		panic(fmt.Errorf("couldn't fmt: %v", err))
	}
}

func goObjectFile(schema neoschema.Schema) {
	fileName := os.Getenv("GOFILE")
	if fileName == "" {
		fileName = "models.go"
	} else {
		fileName = strings.TrimSuffix(fileName, ".go") + "_models.go"
	}

	code := go_gen.GenerationLine(schema) + `
package ` + os.Getenv("GOPACKAGE") + `

`

	// For all transporters, generate the models by collecting all types
	schemas := []neoschema.PackedType{}
	for _, transporter := range schema.Transporters {
		for _, event := range transporter.Events {
			schemas = append(schemas, event)
		}
		for _, route := range transporter.Routes {
			if route.HasRequest {
				schemas = append(schemas, route.Request)
			}
			if route.HasResponse {
				schemas = append(schemas, route.Response)
			}
		}
	}

	generated, err := languages.GenerateObjects(languages.LanguageGo, schemas)
	if err != nil {
		panic(fmt.Errorf("couldn't generate objects: %v", err))
	}
	code += generated

	if err := os.WriteFile(fileName, []byte(code), os.ModePerm); err != nil {
		panic(fmt.Errorf("couldn't write objects file: %v", err))
	}

	// Run message pack generation
	cmd := exec.Command("msgp")
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "GOFILE="+fileName)
	if err := cmd.Run(); err != nil {
		panic(fmt.Errorf("couldn't run message pack: %v", err))
	}
}

func goTransporters(schema neoschema.Schema) {
	files, err := languages.GenerateTransporters(languages.LanguageGo, schema)
	if err != nil {
		panic(fmt.Errorf("couldn't generate transporters: %v", err))
	}

	for fileName, code := range files {
		if err := os.WriteFile(fileName, []byte(code), os.ModePerm); err != nil {
			panic(fmt.Errorf("couldn't write transporter file: %v", err))
		}
	}
}
