package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Liphium/neoroute/cmd/neogen/engine"
	"github.com/Liphium/neoroute/cmd/neogen/languages/go_new"
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
	goFiles(schema)

	// Run the go formatter to make sure it's all nice and clean
	cmd = exec.Command("go", "fmt", ".")
	if err := cmd.Run(); err != nil {
		panic(fmt.Errorf("couldn't fmt: %v", err))
	}
}

func goFiles(schema neoschema.Schema) {
	fileName := os.Getenv("GOFILE")
	if fileName == "" {
		fileName = "models.go"
	} else {
		fileName = strings.TrimSuffix(fileName, ".go") + "_models.go"
	}

	engine := engine.NewGenerationEngine(go_new.NewGoConfig())

	files, err := engine.Generate(schema)
	if err != nil {
		panic(fmt.Errorf("couldn't generate files: %v", err))
	}

	for fileName, content := range files {
		if err := os.WriteFile(fileName, []byte(content), os.ModePerm); err != nil {
			panic(fmt.Errorf("couldn't write transporter file: %v", err))
		}
	}

	// Run message pack generation
	cmd := exec.Command("msgp")
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "GOFILE="+fileName)
	if err := cmd.Run(); err != nil {
		panic(fmt.Errorf("couldn't run message pack: %v", err))
	}
}
