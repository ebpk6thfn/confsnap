package audit

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "audit-test-*")
	if err != nil {
		t.Fatalf("tempDir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func TestNew_CreatesDirectory(t *testing.T) {
	dir := filepath.Join(tempDir(t), "logs")
	_, err := New(dir)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Error("expected directory to be created")
	}
}

func TestLog_WritesEntry(t *testing.T) {
	l, _ := New(tempDir(t))
	if err := l.Log("host1", "/etc/ssh/sshd_config", "snapshot", ""); err != nil {
		t.Fatalf("Log: %v", err)
	}

	entries, err := l.ReadDay(time.Now().UTC())
	if err != nil {
		t.Fatalf("ReadDay: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Host != "host1" {
		t.Errorf("host: got %q, want %q", entries[0].Host, "host1")
	}
	if entries[0].Event != "snapshot" {
		t.Errorf("event: got %q, want %q", entries[0].Event, "snapshot")
	}
}

func TestLog_MultipleEntries(t *testing.T) {
	l, _ := New(tempDir(t))
	for i := 0; i < 3; i++ {
		if err := l.Log("host", "/etc/hosts", "diff", "changed"); err != nil {
			t.Fatalf("Log[%d]: %v", i, err)
		}
	}
	entries, err := l.ReadDay(time.Now().UTC())
	if err != nil {
		t.Fatalf("ReadDay: %v", err)
	}
	if len(entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(entries))
	}
}

func TestReadDay_NoFile_ReturnsNil(t *testing.T) {
	l, _ := New(tempDir(t))
	entries, err := l.ReadDay(time.Now().UTC())
	if err != nil {
		t.Fatalf("ReadDay: %v", err)
	}
	if entries != nil {
		t.Errorf("expected nil entries, got %v", entries)
	}
}

func TestLog_EntryFields(t *testing.T) {
	l, _ := New(tempDir(t))
	_ = l.Log("srv01", "/etc/nginx/nginx.conf", "alert", "threshold exceeded")

	entries, _ := l.ReadDay(time.Now().UTC())
	if len(entries) == 0 {
		t.Fatal("no entries")
	}
	e := entries[0]
	if e.File != "/etc/nginx/nginx.conf" {
		t.Errorf("file: got %q", e.File)
	}
	if e.Details != "threshold exceeded" {
		t.Errorf("details: got %q", e.Details)
	}
	if e.Timestamp.IsZero() {
		t.Error("timestamp should not be zero")
	}
}
