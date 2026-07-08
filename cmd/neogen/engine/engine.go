package engine

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/Liphium/neoroute/cmd/neogen/util"
	"github.com/Liphium/neoroute/neoschema"
)

// LanguageConfig defines how a language should be generated
type LanguageConfig struct {
	Name       string
	GetType    func(packed neoschema.PackedType) (string, error)
	GetFileMap func(schema neoschema.Schema) (map[string]TemplateData, error)
	Funcs      template.FuncMap
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
	for k, v := range e.config.Funcs {
		baseFuncs[k] = v
	}

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

		results[fileName] = out.String()
	}

	return results, nil
}
