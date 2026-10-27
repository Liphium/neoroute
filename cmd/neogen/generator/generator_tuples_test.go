package generator_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/Liphium/neoroute/v2/cmd/neogen/generator"
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
}
