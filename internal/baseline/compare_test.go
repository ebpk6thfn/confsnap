package baseline_test

import (
	"testing"
	"time"

	"github.com/yourorg/confsnap/internal/baseline"
	"github.com/yourorg/confsnap/internal/snapshot"
)

func TestCompareToCurrent_NoChange(t *testing.T) {
	content := []byte("nameserver 8.8.8.8")
	snap := snapshot.New("host1", "/etc/resolv.conf", content)
	b := &baseline.Baseline{
		Name:      "test",
		CreatedAt: time.Now(),
		Snapshots: map[string]*snapshot.Snapshot{
			"host1:/etc/resolv.conf": snap,
		},
	}
	current := map[string]*snapshot.Snapshot{
		"host1:/etc/resolv.conf": snapshot.New("host1", "/etc/resolv.conf", content),
	}
	results := baseline.CompareToCurrent(b, current)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Missing {
		t.Error("expected not missing")
	}
	if results[0].Diff.Changed {
		t.Error("expected no change")
	}
}

func TestCompareToCurrent_WithChange(t *testing.T) {
	base := snapshot.New("host1", "/etc/hosts", []byte("127.0.0.1 localhost"))
	now := snapshot.New("host1", "/etc/hosts", []byte("127.0.0.1 localhost\n10.0.0.1 db"))
	b := &baseline.Baseline{
		Name:      "test",
		CreatedAt: time.Now(),
		Snapshots: map[string]*snapshot.Snapshot{"host1:/etc/hosts": base},
	}
	current := map[string]*snapshot.Snapshot{"host1:/etc/hosts": now}

	results := baseline.CompareToCurrent(b, current)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if !results[0].Diff.Changed {
		t.Error("expected change detected")
	}
}

func TestCompareToCurrent_MissingFile(t *testing.T) {
	snap := snapshot.New("host2", "/etc/ssh/sshd_config", []byte("Port 22"))
	b := &baseline.Baseline{
		Name:      "test",
		CreatedAt: time.Now(),
		Snapshots: map[string]*snapshot.Snapshot{"host2:/etc/ssh/sshd_config": snap},
	}
	results := baseline.CompareToCurrent(b, map[string]*snapshot.Snapshot{})
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if !results[0].Missing {
		t.Error("expected missing=true")
	}
}

func TestSnapshotKey_Format(t *testing.T) {
	key := baseline.SnapshotKey("web01", "/etc/nginx/nginx.conf")
	want := "web01:/etc/nginx/nginx.conf"
	if key != want {
		t.Errorf("got %q, want %q", key, want)
	}
}
