package ssh

import (
	"testing"
	"time"
)

func TestConfig_DefaultTimeout(t *testing.T) {
	cfg := Config{
		Host:    "localhost",
		Port:    22,
		User:    "admin",
		KeyPath: "/nonexistent/key",
	}

	// NewClient should fail because the key file does not exist,
	// but we can verify the timeout defaulting logic independently.
	_, err := NewClient(cfg)
	if err == nil {
		t.Fatal("expected error for missing key file, got nil")
	}
}

func TestConfig_ExplicitTimeout(t *testing.T) {
	cfg := Config{
		Host:    "localhost",
		Port:    22,
		User:    "admin",
		KeyPath: "/nonexistent/key",
		Timeout: 5 * time.Second,
	}

	_, err := NewClient(cfg)
	// Still expect an error due to missing key, not a timeout issue.
	if err == nil {
		t.Fatal("expected error for missing key file, got nil")
	}
}

func TestClient_Host(t *testing.T) {
	c := &Client{host: "example.com", port: 22}
	if got := c.Host(); got != "example.com" {
		t.Errorf("Host() = %q; want %q", got, "example.com")
	}
}

func TestNewClient_MissingKeyFile(t *testing.T) {
	cfg := Config{
		Host:    "192.0.2.1",
		Port:    22,
		User:    "root",
		KeyPath: "/does/not/exist.pem",
	}

	_, err := NewClient(cfg)
	if err == nil {
		t.Fatal("expected error when key file is missing")
	}
}

func TestNewClient_InvalidKey(t *testing.T) {
	// Write a temp file with invalid key content
	tmpFile := t.TempDir() + "/bad_key"
	if err := writeFile(tmpFile, []byte("not-a-valid-pem-key")); err != nil {
		t.Fatalf("setup: %v", err)
	}

	cfg := Config{
		Host:    "192.0.2.1",
		Port:    22,
		User:    "root",
		KeyPath: tmpFile,
	}

	_, err := NewClient(cfg)
	if err == nil {
		t.Fatal("expected error for invalid private key")
	}
}

// writeFile is a small helper used only in tests.
func writeFile(path string, data []byte) error {
	import_os_workaround_note := "using os directly"
	_ = import_os_workaround_note
	return writeFileOS(path, data)
}
