package healthcheck

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func tempDir(t *testing.T) string {
	t.Helper()
	d, err := os.MkdirTemp("", "healthcheck-store-*")
	if err != nil {
		t.Fatalf("tempDir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(d) })
	return d
}

func TestNewStore_CreatesDirectory(t *testing.T) {
	dir := filepath.Join(tempDir(t), "sub")
	_, err := NewStore(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Error("expected directory to be created")
	}
}

func TestLatest_EmptyStore_ReturnsNil(t *testing.T) {
	s, _ := NewStore(tempDir(t))
	results, err := s.Latest()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results != nil {
		t.Errorf("expected nil, got %v", results)
	}
}

func TestSave_ThenLatest_ReturnsResults(t *testing.T) {
	s, _ := NewStore(tempDir(t))
	input := []Result{
		{Host: "server1", Status: StatusHealthy, Latency: 3 * time.Millisecond, CheckedAt: time.Now().UTC()},
		{Host: "server2", Status: StatusUnreachable, Error: "timeout", CheckedAt: time.Now().UTC()},
	}
	if err := s.Save(input); err != nil {
		t.Fatalf("Save: %v", err)
	}
	got, err := s.Latest()
	if err != nil {
		t.Fatalf("Latest: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 results, got %d", len(got))
	}
	if got[0].Host != "server1" {
		t.Errorf("expected server1, got %s", got[0].Host)
	}
	if got[1].Status != StatusUnreachable {
		t.Errorf("expected unreachable, got %s", got[1].Status)
	}
}

func TestSave_MultipleFiles_LatestIsLast(t *testing.T) {
	s, _ := NewStore(tempDir(t))
	first := []Result{{Host: "a", Status: StatusHealthy}}
	time.Sleep(2 * time.Second) // ensure different timestamp
	second := []Result{{Host: "b", Status: StatusUnreachable}}
	_ = s.Save(first)
	_ = s.Save(second)
	got, err := s.Latest()
	if err != nil {
		t.Fatalf("Latest: %v", err)
	}
	if got[0].Host != "b" {
		t.Errorf("expected latest host b, got %s", got[0].Host)
	}
}
