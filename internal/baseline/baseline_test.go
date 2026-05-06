package baseline_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourorg/confsnap/internal/baseline"
	"github.com/yourorg/confsnap/internal/snapshot"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "baseline-test-*")
	if err != nil {
		t.Fatalf("tempDir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func makeBaseline(name string) *baseline.Baseline {
	snap := snapshot.New("host1", "/etc/hosts", []byte("127.0.0.1 localhost"))
	return &baseline.Baseline{
		Name:      name,
		CreatedAt: time.Now(),
		Snapshots: map[string]*snapshot.Snapshot{
			"host1:/etc/hosts": snap,
		},
	}
}

func TestNewManager_CreatesDirectory(t *testing.T) {
	dir := filepath.Join(tempDir(t), "baselines")
	_, err := baseline.NewManager(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Error("expected directory to be created")
	}
}

func TestSave_ThenLoad_ReturnsBaseline(t *testing.T) {
	m, _ := baseline.NewManager(tempDir(t))
	b := makeBaseline("prod-2024")

	if err := m.Save(b); err != nil {
		t.Fatalf("Save: %v", err)
	}
	got, err := m.Load("prod-2024")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got == nil {
		t.Fatal("expected non-nil baseline")
	}
	if got.Name != b.Name {
		t.Errorf("name: got %q, want %q", got.Name, b.Name)
	}
	if len(got.Snapshots) != 1 {
		t.Errorf("snapshots: got %d, want 1", len(got.Snapshots))
	}
}

func TestLoad_Missing_ReturnsNil(t *testing.T) {
	m, _ := baseline.NewManager(tempDir(t))
	got, err := m.Load("nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != nil {
		t.Error("expected nil for missing baseline")
	}
}

func TestList_ReturnsNames(t *testing.T) {
	m, _ := baseline.NewManager(tempDir(t))
	for _, name := range []string{"alpha", "beta", "gamma"} {
		if err := m.Save(makeBaseline(name)); err != nil {
			t.Fatalf("Save %s: %v", name, err)
		}
	}
	names, err := m.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(names) != 3 {
		t.Errorf("got %d names, want 3", len(names))
	}
}
