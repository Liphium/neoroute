package generator_test

import (
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/Liphium/neoroute/cmd/neogen/generator"
	"github.com/Liphium/neoroute/cmd/neogen/languages"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// serverSource is a minimal neoroute server registering custom objects (nested, pointer, slice, recursive)
// and using one of them as a route request/response.
const serverSource = `package main

import (
	"reflect"

	"github.com/Liphium/neoroute"
	"github.com/Liphium/neoroute/neoschema"
)

type Address struct {
	City string ` + "`msg:\"city\"`" + `
	Zip  *int   ` + "`msg:\"zip\"`" + `
}

type Inner struct {
	Note string ` + "`msg:\"note\"`" + `
}

type User struct {
	Name    string   ` + "`msg:\"name\"`" + `
	Addr    *Address ` + "`msg:\"addr\"`" + `
	Tags    []string ` + "`msg:\"tags\"`" + `
	Inner   Inner    ` + "`msg:\"inner\"`" + `
	Friends []User   ` + "`msg:\"friends\"`" + `
}

// StandaloneInner is not referenced by any route or event, it must still be generated because it is inside of Standalone.
type StandaloneInner struct {
	Info string ` + "`msg:\"info\"`" + `
}

// Standalone is not referenced by any route or event, it must still be generated.
type Standalone struct {
	Code  int      ` + "`msg:\"code\"`" + `
	Label *string  ` + "`msg:\"label\"`" + `
	Items []Inner  ` + "`msg:\"items\"`" + `
}

type mockTransporter struct{ routes map[string]neoschema.RequestResponse }

func (m *mockTransporter) Type() neoschema.TransporterType                 { return neoschema.TransporterHTTP }
func (m *mockTransporter) GetSchema() map[string]neoschema.RequestResponse { return m.routes }
func (m *mockTransporter) GetRegistries() []neoroute.IEventRegistry        { return nil }

func main() {
	g := neoschema.NewGenerator()
	g.Add(Standalone{})
	g.Transporter("http", &mockTransporter{routes: map[string]neoschema.RequestResponse{
		"getUser": {
			HasRequest:     true,
			Request:        reflect.TypeFor[User](),
			HasResponse:    true,
			Response:       reflect.TypeFor[User](),
			CanReturnError: true,
		},
	}})
	g.PrintIfDesired()
}
`

func TestCustomObjectsEndToEnd(t *testing.T) {
	projectDir, _ := filepath.Abs("../../../")

	// 1. Create a temporary server module that depends on the local neoroute checkout
	serverDir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(serverDir, "main.go"), []byte(serverSource), 0644))
	run := func(dir string, args ...string) {
		t.Helper()
		cmd := exec.Command("go", args...)
		cmd.Dir = dir
		out, err := cmd.CombinedOutput()
		require.NoError(t, err, "go %v\n%s", args, out)
	}
	run(serverDir, "mod", "init", "custom-objects-server")
	run(serverDir, "mod", "edit", "-replace", "github.com/Liphium/neoroute="+projectDir)
	run(serverDir, "mod", "tidy")

	// 2. Extract the schema from the server (same path neogen uses)
	schema, err := generator.GetSchema(serverDir, "go run . --neo-generate")
	require.NoError(t, err)
	require.Len(t, schema.CustomObjects, 1, "only Standalone is registered via Add")

	// 3. Generate Go models
	goDir := t.TempDir()
	t.Setenv("GOPACKAGE", "generated")
	require.NoError(t, os.WriteFile(filepath.Join(goDir, "go.mod"), []byte("module generated\n\ngo 1.27\n"), 0644))
	cwd, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(goDir))
	generator.GenerateGo(schema)
	require.NoError(t, os.Chdir(cwd))

	goStructs := parseGoStructs(t, filepath.Join(goDir, "models.go"))
	assert.Equal(t, map[string]map[string]string{
		"User": {
			"Name":    "string",
			"Addr":    "*Address",
			"Tags":    "[]string",
			"Inner":   "Inner",
			"Friends": "[]User",
		},
		"Address": {
			"City": "string",
			"Zip":  "*int64",
		},
		"Inner": {
			"Note": "string",
		},
		"Standalone": {
			"Code":  "int64",
			"Label": "*string",
			"Items": "[]Inner",
		},
	}, goStructs, "generated Go structs must match exactly")

	// The generated Go package must compile
	run(goDir, "mod", "edit", "-replace", "github.com/Liphium/neoroute/client="+filepath.Join(projectDir, "client"))
	run(goDir, "mod", "edit", "-replace", "github.com/Liphium/neoroute/client/transporter/http="+filepath.Join(projectDir, "client", "transporter", "http"))
	run(goDir, "mod", "edit", "-replace", "github.com/Liphium/neoroute/client/transporter/websocket="+filepath.Join(projectDir, "client", "transporter", "websocket"))
	run(goDir, "mod", "tidy")
	run(goDir, "vet", "./...")

	// 4. Generate TypeScript models + connector
	tsDir := t.TempDir()
	require.NoError(t, os.Chdir(tsDir))
	generator.GenerateWithConfig(schema, languages.NewTSConfig())
	require.NoError(t, os.Chdir(cwd))

	tsModels, err := os.ReadFile(filepath.Join(tsDir, "models.ts"))
	require.NoError(t, err)
	assert.Equal(t, map[string]map[string]string{
		"User": {
			"name":    "string",
			"addr":    "Address | null",
			"tags":    "string[]",
			"inner":   "Inner",
			"friends": "User[]",
		},
		"Address": {
			"city": "string",
			"zip":  "bigint | null",
		},
		"Inner": {
			"note": "string",
		},
		"Standalone": {
			"code":  "bigint",
			"label": "string | null",
			"items": "Inner[]",
		},
	}, parseTSInterfaces(t, string(tsModels)), "generated TS interfaces must match exactly")

	connector, err := os.ReadFile(filepath.Join(tsDir, "connector_http.ts"))
	require.NoError(t, err)
	assert.Contains(t, string(connector), "type User,", "connector must import custom object User")
	assert.Contains(t, string(connector), "getUser(payload: User): Promise<User | UserError>", "connector method must use User")
}

// parseGoStructs returns struct name -> field name -> field type (as source text) for all top-level structs.
func parseGoStructs(t *testing.T, path string) map[string]map[string]string {
	t.Helper()
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	require.NoError(t, err)

	res := map[string]map[string]string{}
	for _, decl := range file.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.TYPE {
			continue
		}
		for _, spec := range gd.Specs {
			ts := spec.(*ast.TypeSpec)
			st, ok := ts.Type.(*ast.StructType)
			if !ok {
				continue
			}
			fields := map[string]string{}
			for _, f := range st.Fields.List {
				var typ strings.Builder
				require.NoError(t, printer.Fprint(&typ, fset, f.Type))
				for _, name := range f.Names {
					fields[name.Name] = typ.String()
				}
			}
			res[ts.Name.Name] = fields
		}
	}
	return res
}

var (
	tsInterfaceRe = regexp.MustCompile(`(?s)export interface (\w+) \{(.*?)\}`)
	tsFieldRe     = regexp.MustCompile(`(?m)^\s*(\w+): (.+?);\s*$`)
)

// parseTSInterfaces returns interface name -> field name -> field type for all exported interfaces.
func parseTSInterfaces(t *testing.T, src string) map[string]map[string]string {
	t.Helper()
	res := map[string]map[string]string{}
	for _, m := range tsInterfaceRe.FindAllStringSubmatch(src, -1) {
		fields := map[string]string{}
		for _, f := range tsFieldRe.FindAllStringSubmatch(m[2], -1) {
			fields[f[1]] = f[2]
		}
		res[m[1]] = fields
	}
	return res
}
