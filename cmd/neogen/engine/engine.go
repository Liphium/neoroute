package engine

import (
	"bytes"
	"fmt"
	"maps"
	"text/template"

	"github.com/Liphium/neoroute/v2/cmd/neogen/util"
	"github.com/Liphium/neoroute/v2/neoschema"
)

// LanguageConfig defines how a language should be generated
type LanguageConfig struct {
	Name        string
	GetType     func(packed neoschema.PackedType) (string, error)
	GetFileMap  func(schema neoschema.Schema) (map[string]TemplateData, error)
	Funcs       template.FuncMap
	ProcessFile func(filename string, source []byte) ([]byte, error)
}

// TemplateData holds the template and the data object for execution
type TemplateData struct {
	Template string
	Object   any
}

// GenerationEngine runs the templates
type GenerationEngine struct {
	config LanguageConfig
}

func NewGenerationEngine(config LanguageConfig) *GenerationEngine {
	return &GenerationEngine{config: config}
}

func (e *GenerationEngine) Generate(schema neoschema.Schema) (map[string]string, error) {
	fileMap, err := e.config.GetFileMap(schema)
	if err != nil {
		return nil, err
	}

	// Default funcs
	baseFuncs := template.FuncMap{
		"camel":   util.ToCamelCase,
		"getType": e.config.GetType,
		"toStruct": func(p neoschema.PackedType) *neoschema.StructType {
			if s, ok := p.(*neoschema.StructType); ok {
				return s
			}
			return nil
		},
	}

	// Merge with language specific funcs
	maps.Copy(baseFuncs, e.config.Funcs)

	results := make(map[string]string)
	for fileName, data := range fileMap {
		tmpl, err := template.New(fileName).Funcs(baseFuncs).Parse(data.Template)
		if err != nil {
			return nil, fmt.Errorf("failed to parse template for %s: %v", fileName, err)
		}

		var out bytes.Buffer
		err = tmpl.Execute(&out, map[string]any{
			"Object": data.Object,
			"Schema": schema, // Global access if needed
		})
		if err != nil {
			return nil, fmt.Errorf("failed to execute template for %s: %v", fileName, err)
		}

		source := out.Bytes()
		if e.config.ProcessFile != nil {
			source, err = e.config.ProcessFile(fileName, source)
			if err != nil {
				return nil, fmt.Errorf("failed to process %s: %w", fileName, err)
			}
		}
		results[fileName] = string(source)
	}

	return results, nil
}
