package generator

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/Liphium/neoroute/cmd/neogen/engine"
	"github.com/Liphium/neoroute/cmd/neogen/languages"
	"github.com/Liphium/neoroute/neoschema"
)

func GenerateGo(schema neoschema.Schema) {

	// Generate object and transporters
	goFiles(schema)

	// Run the go formatter to make sure it's all nice and clean
	cmd := exec.Command("go", "fmt", ".")
	if err := cmd.Run(); err != nil {
		panic(fmt.Errorf("couldn't fmt: %v", err))
	}
}

func goFiles(schema neoschema.Schema) {
	modelFileName, config := languages.NewGoConfig()
	engine := engine.NewGenerationEngine(config)

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
	cmd.Env = append(cmd.Env, "GOFILE="+modelFileName)
	if err := cmd.Run(); err != nil {
		panic(fmt.Errorf("couldn't run message pack: %v", err))
	}
}
