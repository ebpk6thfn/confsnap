package retention

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func tempDir(t *testing.T) string {
	t.Helper()
	d, err := os.MkdirTemp("", "retention-test-*")
	if err != nil {
		t.Fatalf("tempDir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(d) })
	return d
}

func touch(t *testing.T, dir, name string, modTime time.Time) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(`{}`), 0o644); err != nil {
		t.Fatalf("touch: %v", err)
	}
	if err := os.Chtimes(p, modTime, modTime); err != nil {
		t.Fatalf("chtimes: %v", err)
	}
	return p
}

func TestNew_InvalidPolicy(t *testing.T) {
	_, err := New("", Policy{})
	if err == nil {
		t.Fatal("expected error for empty storeDir")
	}
	_, err = New("/tmp", Policy{MaxAge: -1})
	if err == nil {
		t.Fatal("expected error for negative MaxAge")
	}
}

func TestPrune_ByAge(t *testing.T) {
	dir := tempDir(t)
	now := time.Now()
	touch(t, dir, "host1_etc_sshd_config_001.json", now.Add(-48*time.Hour))
	touch(t, dir, "host1_etc_sshd_config_002.json", now.Add(-1*time.Hour))

	m, _ := New(dir, Policy{MaxAge: 24 * time.Hour})
	removed, err := m.Prune(now)
	if err != nil {
		t.Fatalf("Prune: %v", err)
	}
	if len(removed) != 1 {
		t.Fatalf("expected 1 removed, got %d", len(removed))
	}
}

func TestPrune_ByCount(t *testing.T) {
	dir := tempDir(t)
	now := time.Now()
	for _, name := range []string{
		"host1_etc_nginx_001.json",
		"host1_etc_nginx_002.json",
		"host1_etc_nginx_003.json",
	} {
		touch(t, dir, name, now)
	}

	m, _ := New(dir, Policy{MaxCount: 2})
	removed, err := m.Prune(now)
	if err != nil {
		t.Fatalf("Prune: %v", err)
	}
	if len(removed) != 1 {
		t.Fatalf("expected 1 removed, got %d", len(removed))
	}
}

func TestPrune_EmptyDir(t *testing.T) {
	dir := tempDir(t)
	m, _ := New(dir, Policy{MaxCount: 5})
	removed, err := m.Prune(time.Now())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(removed) != 0 {
		t.Fatalf("expected 0 removed, got %d", len(removed))
	}
}

func TestGroupKey(t *testing.T) {
	cases := []struct {
		name string
		want string
	}{
		{"host1_etc_sshd_config_20240101.json", "host1_etc_sshd_config"},
		{"single.json", "single"},
	}
	for _, c := range cases {
		got := groupKey(c.name)
		if got != c.want {
			t.Errorf("groupKey(%q) = %q, want %q", c.name, got, c.want)
		}
	}
}
