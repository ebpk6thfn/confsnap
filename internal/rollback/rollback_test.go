package rollback

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/confsnap/confsnap/internal/snapshot"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "rollback-test-*")
	if err != nil {
		t.Fatalf("tempDir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func makeSnapshot(host, path, content string) *snapshot.Snapshot {
	s, _ := snapshot.New(host, path, []byte(content))
	return s
}

func TestNew_CreatesDirectory(t *testing.T) {
	dir := filepath.Join(tempDir(t), "rollback")
	_, err := New(dir)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Error("expected directory to be created")
	}
}

func TestSave_ThenLoad_ReturnsEntry(t *testing.T) {
	m, _ := New(tempDir(t))
	s := makeSnapshot("web01", "/etc/nginx/nginx.conf", "worker_processes 4;")

	id, err := m.Save(s)
	if err != nil {
		t.Fatalf("Save: %v", err)
	}
	if id == "" {
		t.Fatal("expected non-empty id")
	}

	entry, err := m.Load(id)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if entry == nil {
		t.Fatal("expected entry, got nil")
	}
	if string(entry.Content) != "worker_processes 4;" {
		t.Errorf("content mismatch: got %q", entry.Content)
	}
	if entry.Host != "web01" {
		t.Errorf("host mismatch: got %q", entry.Host)
	}
}

func TestLoad_Missing_ReturnsNil(t *testing.T) {
	m, _ := New(tempDir(t))
	entry, err := m.Load("nonexistent__id__20240101T000000Z")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if entry != nil {
		t.Error("expected nil for missing entry")
	}
}

func TestDelete_RemovesEntry(t *testing.T) {
	m, _ := New(tempDir(t))
	s := makeSnapshot("db01", "/etc/mysql/my.cnf", "[mysqld]")

	id, _ := m.Save(s)
	if err := m.Delete(id); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	entry, err := m.Load(id)
	if err != nil {
		t.Fatalf("Load after delete: %v", err)
	}
	if entry != nil {
		t.Error("expected nil after deletion")
	}
}

func TestDelete_MissingEntry_NoError(t *testing.T) {
	m, _ := New(tempDir(t))
	if err := m.Delete("ghost__file__20240101T000000Z"); err != nil {
		t.Errorf("expected no error deleting missing entry, got: %v", err)
	}
}

func TestEntryID_ContainsTimestamp(t *testing.T) {
	now := time.Date(2024, 6, 15, 12, 30, 0, 0, time.UTC)
	id := entryID("host1", "/etc/hosts", now)
	if indexOf(id, "20240615T123000Z") < 0 {
		t.Errorf("expected timestamp in id, got %q", id)
	}
}
