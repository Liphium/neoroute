package generator_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/Liphium/neoroute/v2/cmd/neogen/generator"
	"github.com/Liphium/neoroute/v2/cmd/neogen/languages"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const tupleServerSource = `package main

import "github.com/Liphium/neoroute/v2/neoschema"

type Coordinates struct {
	neoschema.IsTuple ` + "`msg:\"-\"`" + `
	X float64 ` + "`msg:\"x\"`" + `
	Y float64 ` + "`msg:\"y\"`" + `
}

type Location struct {
	Name string ` + "`msg:\"name\"`" + `
	Coordinates Coordinates ` + "`msg:\"coordinates\"`" + `
}

func main() {
	g := neoschema.NewGenerator()
	g.Add(Location{})
	g.PrintIfDesired()
}
`

func TestTuplesEndToEnd(t *testing.T) {
	projectDir, err := filepath.Abs("../../../")
	require.NoError(t, err)
	run := func(dir string, args ...string) {
		t.Helper()
		cmd := exec.Command("go", args...)
		cmd.Dir = dir
		out, err := cmd.CombinedOutput()
		require.NoError(t, err, "go %v\n%s", args, out)
	}

	serverDir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(serverDir, "main.go"), []byte(tupleServerSource), 0644))
	run(serverDir, "mod", "init", "tuples-server")
	run(serverDir, "mod", "edit", "-replace", "github.com/Liphium/neoroute/v2="+projectDir)
	run(serverDir, "mod", "tidy")

	schema, err := generator.GetSchema(serverDir, "go run . --neo-generate")
	require.NoError(t, err)

	goDir := t.TempDir()
	t.Setenv("GOPACKAGE", "generated")
	require.NoError(t, os.WriteFile(filepath.Join(goDir, "go.mod"), []byte("module generated\n\ngo 1.27\n"), 0644))
	t.Chdir(goDir)
	generator.GenerateGo(schema)

	assert.Equal(t, map[string]map[string]string{
		"Coordinates": {
			"X": "float64",
			"Y": "float64",
		},
		"Location": {
			"Name":        "string",
			"Coordinates": "Coordinates",
		},
	}, parseGoStructs(t, filepath.Join(goDir, "models.go")))

	models, err := os.ReadFile(filepath.Join(goDir, "models.go"))
	require.NoError(t, err)
	assert.Contains(t, string(models), "//msgp:tuple Coordinates\ntype Coordinates struct")
	assert.NotContains(t, string(models), "//msgp:tuple Location")

	run(goDir, "mod", "tidy")
	run(goDir, "vet", "./...")

	tsDir := t.TempDir()
	t.Chdir(tsDir)
	generator.GenerateWithConfig(schema, languages.NewTSConfig())

	tsModels, err := os.ReadFile(filepath.Join(tsDir, "models.ts"))
	require.NoError(t, err)
	assert.Equal(t, map[string][]string{
		"Coordinates": {"x: number", "y: number"},
	}, parseTSTuples(t, string(tsModels)), "tuple elements must retain their declaration order")
	assert.Equal(t, map[string]map[string]string{
		"Location": {
			"name":        "string",
			"coordinates": "Coordinates",
		},
	}, parseTSInterfaces(t, string(tsModels)), "only non-tuple objects should generate interfaces")
}

// parseTSTuples returns tuple names and their ordered, labeled element types.
// It handles the one-element-per-line format emitted by the generator.
func parseTSTuples(t *testing.T, source string) map[string][]string {
	t.Helper()
	tupleRe := regexp.MustCompile(`(?ms)^export type (\w+) = \[\n(.*?)^\];`)
	tuples := map[string][]string{}
	for _, match := range tupleRe.FindAllStringSubmatch(source, -1) {
		elements := []string{}
		for line := range strings.SplitSeq(strings.TrimSpace(match[2]), "\n") {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			require.True(t, strings.HasSuffix(line, ","), "tuple element must end with a comma: %s", line)
			elements = append(elements, strings.TrimSuffix(line, ","))
		}
		tuples[match[1]] = elements
	}
	return tuples
}
