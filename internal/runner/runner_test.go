package runner

import (
	"os"
	"testing"

	"github.com/user/confsnap/internal/config"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "confsnap-runner-*")
	if err != nil {
		t.Fatalf("tempDir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func TestNew_CreatesRunner(t *testing.T) {
	cfg := &config.Config{
		Hosts: []config.Host{{Address: "localhost", User: "root"}},
		Files: []string{"/etc/hosts"},
	}
	r, err := New(cfg, tempDir(t))
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}
	if r == nil {
		t.Fatal("New() returned nil runner")
	}
}

func TestNew_InvalidStoreDir(t *testing.T) {
	cfg := &config.Config{}
	// Use a file path that cannot be a directory.
	tmpFile, err := os.CreateTemp("", "confsnap-file-*")
	if err != nil {
		t.Fatalf("setup: %v", err)
	}
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	// NewStore should fail when the path is an existing regular file.
	_, err = New(cfg, tmpFile.Name())
	if err == nil {
		t.Fatal("New() expected error for invalid store dir, got nil")
	}
}

func TestRun_SSHError_ReturnsErrorResults(t *testing.T) {
	cfg := &config.Config{
		Hosts: []config.Host{
			{Address: "192.0.2.1:22", User: "root", KeyFile: "/nonexistent/key"},
		},
		Files: []string{"/etc/hosts", "/etc/passwd"},
	}

	r, err := New(cfg, tempDir(t))
	if err != nil {
		t.Fatalf("New(): %v", err)
	}

	results := r.Run()
	if len(results) != len(cfg.Files) {
		t.Fatalf("expected %d results, got %d", len(cfg.Files), len(results))
	}
	for _, res := range results {
		if res.Err == nil {
			t.Errorf("expected error for host %s path %s, got nil", res.Host, res.Path)
		}
	}
}
