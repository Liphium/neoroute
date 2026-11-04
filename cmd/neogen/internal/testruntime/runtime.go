// Package testruntime runs generated-code test projects in isolated language containers. The host only needs Go and access to a Docker daemon.
//
// This package is written and maintained by AI. Has been validated by Unbreathable on 04.11.2026.
package testruntime

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	tcexec "github.com/testcontainers/testcontainers-go/exec"
)

// Config is stored in each language project's runtime.json. Commands are argv
// lists executed sequentially in /workspace. Use ["sh", "-c", "..."] if needed.
type Config struct {
	Image          string            `json:"image"`
	Commands       [][]string        `json:"commands"`
	Env            map[string]string `json:"env,omitempty"`
	TimeoutSeconds int               `json:"timeout_seconds,omitempty"`
}

// Run overlays generated files under generated/ in a copy of projectDir, then
// runs the project's configured commands. No host directories are mounted.
// Docker is optional locally; NEOGEN_REQUIRE_DOCKER=1 or CI=true makes it required.
func Run(t *testing.T, projectDir string, generated map[string]string) {
	t.Helper()
	// SELinux enforcing hosts block Ryuk's Docker socket access unless the
	// cleanup helper is privileged. Set this before Testcontainers caches its
	// configuration; leave explicit environment settings alone.
	if _, configured := os.LookupEnv("TESTCONTAINERS_RYUK_CONTAINER_PRIVILEGED"); !configured {
		if enforcing, err := os.ReadFile("/sys/fs/selinux/enforce"); err == nil && strings.TrimSpace(string(enforcing)) == "1" {
			t.Setenv("TESTCONTAINERS_RYUK_CONTAINER_PRIVILEGED", "true")
		}
	}
	configBytes, err := os.ReadFile(filepath.Join(projectDir, "runtime.json"))
	if err != nil {
		t.Fatal(err)
	}
	var config Config
	if err := json.Unmarshal(configBytes, &config); err != nil {
		t.Fatal(err)
	}
	if config.Image == "" || len(config.Commands) == 0 {
		t.Fatal("runtime.json requires image and commands")
	}
	for _, command := range config.Commands {
		if len(command) == 0 {
			t.Fatal("runtime command cannot be empty")
		}
	}
	if config.TimeoutSeconds < 0 {
		t.Fatal("timeout_seconds cannot be negative")
	}
	timeout := 5 * time.Minute
	if config.TimeoutSeconds > 0 {
		timeout = time.Duration(config.TimeoutSeconds) * time.Second
	}
	ctx, cancel := context.WithTimeout(t.Context(), timeout)
	defer cancel()
	if err := dockerAvailable(ctx); err != nil {
		if os.Getenv("NEOGEN_REQUIRE_DOCKER") == "1" || strings.EqualFold(os.Getenv("CI"), "true") || os.Getenv("CI") == "1" {
			t.Fatalf("Docker required: %v", err)
		}
		t.Skipf("Docker unavailable: %v (set NEOGEN_REQUIRE_DOCKER=1 to require it)", err)
	}
	workspace := t.TempDir()
	err = fs.WalkDir(os.DirFS(projectDir), ".", func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			switch entry.Name() {
			case "node_modules", "dist", "generated", ".git":
				return fs.SkipDir
			}
			return os.MkdirAll(filepath.Join(workspace, path), 0755)
		}
		data, err := os.ReadFile(filepath.Join(projectDir, path))
		if err != nil {
			return err
		}
		return os.WriteFile(filepath.Join(workspace, path), data, 0644)
	})
	if err != nil {
		t.Fatal(err)
	}
	for name, source := range generated {
		if !filepath.IsLocal(name) {
			t.Fatalf("generated path must be relative: %q", name)
		}
		target := filepath.Join(workspace, "generated", name)
		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(target, []byte(source), 0644); err != nil {
			t.Fatal(err)
		}
	}

	var inputs []testcontainers.ContainerFile
	err = filepath.WalkDir(workspace, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		relative, err := filepath.Rel(workspace, path)
		if err != nil {
			return err
		}
		inputs = append(inputs, testcontainers.ContainerFile{HostFilePath: path, ContainerFilePath: "/workspace/" + filepath.ToSlash(relative), FileMode: 0644})
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		Image: config.Image, WorkingDir: "/workspace", Env: config.Env,
		Cmd:   []string{"sleep", "infinity"},
		Files: inputs, Started: true,
	})
	if container != nil {
		defer func() {
			cleanup, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			if err := container.Terminate(cleanup); err != nil {
				t.Errorf("terminate runtime container: %v", err)
			}
		}()
	}
	if err != nil {
		t.Fatalf("start %s: %v", config.Image, err)
	}
	for _, command := range config.Commands {
		t.Logf("running %q in %s", command, config.Image)
		code, output, err := container.Exec(ctx, command, tcexec.Multiplexed())
		var logs []byte
		if output != nil {
			var readErr error
			logs, readErr = io.ReadAll(output)
			if readErr != nil && err == nil {
				err = readErr
			}
		}
		if err != nil || code != 0 {
			t.Fatalf("command %q failed (exit %d): %v\n%s", command, code, err, logs)
		}
		t.Log(string(logs))
	}
}

// Testcontainers' Docker discovery can panic when no provider is available.
// Limit recovery to discovery: errors starting images or running tests must fail.
func dockerAvailable(ctx context.Context) (err error) {
	defer func() {
		if value := recover(); value != nil {
			err = fmt.Errorf("Docker discovery: %v", value)
		}
	}()
	probe, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	client, err := testcontainers.NewDockerClientWithOpts(probe)
	if err != nil {
		return err
	}
	defer client.Close()
	_, err = client.Ping(probe)
	return err
}
