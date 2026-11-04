package languages

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/Liphium/neoroute/v2/cmd/neogen/engine"
	"github.com/Liphium/neoroute/v2/cmd/neogen/internal/testschema"
	"github.com/stretchr/testify/require"
)

func TestGoRuntime(t *testing.T) {
	root, err := filepath.Abs("../../..")
	require.NoError(t, err)
	schema, err := testschema.Schema()
	require.NoError(t, err)
	// Go route payloads require generated message methods on named models.
	// Collection roots are exercised by the TypeScript connector tests.
	for _, transport := range schema.Transporters {
		delete(transport.Routes, "collection")
		delete(transport.Routes, "lookup")
	}
	t.Setenv("GOPACKAGE", "generated")
	_, config := NewGoConfig()
	files, err := engine.NewGenerationEngine(config).Generate(schema)
	require.NoError(t, err)
	dir := t.TempDir()
	require.NoError(t, os.CopyFS(dir, os.DirFS("../tests/go")))
	for name, source := range files {
		require.NoError(t, os.WriteFile(filepath.Join(dir, name), []byte(source), 0644))
	}
	ctx, cancel := context.WithTimeout(t.Context(), 2*time.Minute)
	defer cancel()
	run := func(args ...string) {
		t.Helper()
		cmd := exec.CommandContext(ctx, "go", args...)
		cmd.Dir = dir
		output, err := cmd.CombinedOutput()
		require.NoError(t, err, "go %v\n%s", args, output)
		t.Log(string(output))
	}
	for module, path := range map[string]string{
		"github.com/Liphium/neoroute/v2":                              root,
		"github.com/Liphium/neoroute/client/sdk/v2":                   filepath.Join(root, "client/sdk"),
		"github.com/Liphium/neoroute/client/transporter/http/v2":      filepath.Join(root, "client/transporter/http"),
		"github.com/Liphium/neoroute/client/transporter/websocket/v2": filepath.Join(root, "client/transporter/websocket"),
	} {
		run("mod", "edit", "-replace", module+"="+path)
	}
	run("mod", "tidy")
	// Run the pinned generator through Go; no globally installed msgp binary.
	run("run", "github.com/tinylib/msgp@v1.6.4", "-file", "models.go", "-tests=false")
	run("mod", "tidy")
	run("test", "-v", "./...")
}
