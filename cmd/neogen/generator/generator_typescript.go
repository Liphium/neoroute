package generator

import (
	"fmt"
	"os"

	"github.com/Liphium/neoroute/cmd/neogen/engine"
	"github.com/Liphium/neoroute/neoschema"
)

func GenerateWithConfig(schema neoschema.Schema, config engine.LanguageConfig) {
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
}
