package generator_test

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/Liphium/neoroute/cmd/neogen/generator"
)

//go:embed test/*
var testFS embed.FS

func TestGenerate(t *testing.T) {

	// IF YOU MOVE THIS TEST UPDATE THIS
	projectDir, _ := filepath.Abs("../../../")
	testDir := filepath.Join(projectDir, ".testing", "neogen-env")
	_ = os.RemoveAll(testDir)
	if err := os.MkdirAll(testDir, os.ModePerm); err != nil {
		panic(err)
	}

	// 1. Setup go module
	cmd := exec.Command("go", "mod", "init", "test-gen")
	cmd.Dir = testDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("mod init fail: %v\n%s", err, out)
	}
	actualFiles, _ := fs.Sub(testFS, "test")
	_ = os.CopyFS(testDir, actualFiles)

	// Add replaces
	replaces := []string{
		fmt.Sprintf("github.com/Liphium/neoroute=%s", projectDir),
		fmt.Sprintf("github.com/Liphium/neoroute/client=%s/client", projectDir),
		fmt.Sprintf("github.com/Liphium/neoroute/transporter/http=%s/transporter/http", projectDir),
		fmt.Sprintf("github.com/Liphium/neoroute/transporter/websocket=%s/transporter/websocket", projectDir),
	}

	for _, r := range replaces {
		cmd = exec.Command("go", "mod", "edit", "-replace", r)
		cmd.Dir = testDir
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("mod edit replace %s fail: %v\n%s", r, err, out)
		}
	}

	// 2. Install dependencies
	cmd = exec.Command("go", "mod", "tidy")
	cmd.Dir = testDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("mod tidy fail: %v\n%s", err, out)
	}

	// 3. Generate
	os.Setenv("GOPACKAGE", "definitions")
	definitionsDir := filepath.Join(testDir, "definitions")
	_ = os.Mkdir(definitionsDir, os.ModePerm)
	_ = os.Chdir(definitionsDir)
	generator.Generate(generator.GeneratorConfig{
		ServerPath:     "..",
		Command:        "go run .",
		TargetLanguage: "go",
		Verbose:        true,
	})

	// 4. Install more dependencies
	cmd = exec.Command("go", "mod", "tidy")
	cmd.Dir = testDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("mod tidy fail: %v\n%s", err, out)
	}

	// 5. Make sure everything builds
	cmd = exec.Command("go", "vet", "./...")
	cmd.Dir = testDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Errorf("go vet fail: %v\n%s", err, out)
	}

	cmd = exec.Command("go", "build", "-o", "/dev/null", "./...")
	cmd.Dir = testDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Errorf("go build fail: %v\n%s", err, out)
	}
}
