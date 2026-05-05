package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/confsnap/internal/config"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "confsnap-*.yaml")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

const validYAML = `
defaults:
  user: admin
  port: 22
  key_file: ~/.ssh/id_rsa
hosts:
  - name: web-01
    address: 192.168.1.10
  - name: db-01
    address: 192.168.1.20
    user: dbuser
    port: 2222
files:
  - /etc/nginx/nginx.conf
  - /etc/hosts
`

func TestLoad_ValidConfig(t *testing.T) {
	path := writeTemp(t, validYAML)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Hosts) != 2 {
		t.Errorf("expected 2 hosts, got %d", len(cfg.Hosts))
	}
	if len(cfg.Files) != 2 {
		t.Errorf("expected 2 files, got %d", len(cfg.Files))
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := config.Load(filepath.Join(t.TempDir(), "nonexistent.yaml"))
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoad_NoHosts(t *testing.T) {
	yaml := "files:\n  - /etc/hosts\n"
	_, err := config.Load(writeTemp(t, yaml))
	if err == nil {
		t.Fatal("expected validation error for missing hosts")
	}
}

func TestLoad_NoFiles(t *testing.T) {
	yaml := "hosts:\n  - name: h1\n    address: 1.2.3.4\n"
	_, err := config.Load(writeTemp(t, yaml))
	if err == nil {
		t.Fatal("expected validation error for missing files")
	}
}

func TestHost_EffectivePort_UsesDefault(t *testing.T) {
	h := config.Host{Address: "1.2.3.4"}
	d := config.Defaults{Port: 2222}
	if got := h.EffectivePort(d); got != 2222 {
		t.Errorf("expected 2222, got %d", got)
	}
}

func TestHost_EffectivePort_FallbackTo22(t *testing.T) {
	h := config.Host{Address: "1.2.3.4"}
	d := config.Defaults{}
	if got := h.EffectivePort(d); got != 22 {
		t.Errorf("expected 22, got %d", got)
	}
}

func TestHost_EffectiveUser_OverridesDefault(t *testing.T) {
	h := config.Host{Address: "1.2.3.4", User: "override"}
	d := config.Defaults{User: "default"}
	if got := h.EffectiveUser(d); got != "override" {
		t.Errorf("expected override, got %s", got)
	}
}
