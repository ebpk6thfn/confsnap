package cache

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "cache-test-*")
	if err != nil {
		t.Fatalf("tempDir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func TestNew_CreatesDirectory(t *testing.T) {
	dir := filepath.Join(tempDir(t), "subdir")
	_, err := New(dir)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Error("expected directory to be created")
	}
}

func TestGet_MissingEntry_ReturnsNil(t *testing.T) {
	c, _ := New(tempDir(t))
	if got := c.Get("host1", "/etc/hosts"); got != nil {
		t.Errorf("expected nil, got %+v", got)
	}
}

func TestSet_ThenGet_ReturnsEntry(t *testing.T) {
	c, _ := New(tempDir(t))
	if err := c.Set("host1", "/etc/hosts", "abc123", 3600); err != nil {
		t.Fatalf("Set: %v", err)
	}
	e := c.Get("host1", "/etc/hosts")
	if e == nil {
		t.Fatal("expected entry, got nil")
	}
	if e.Checksum != "abc123" {
		t.Errorf("checksum: want abc123, got %s", e.Checksum)
	}
}

func TestGet_ExpiredEntry_ReturnsNil(t *testing.T) {
	c, _ := New(tempDir(t))
	_ = c.Set("host1", "/etc/ssh/sshd_config", "xyz", 1)
	// Manually backdate the entry.
	key := cacheKey("host1", "/etc/ssh/sshd_config")
	c.entries[key].CachedAt = time.Now().Add(-2 * time.Second)
	if got := c.Get("host1", "/etc/ssh/sshd_config"); got != nil {
		t.Errorf("expected nil for expired entry, got %+v", got)
	}
}

func TestEntry_IsExpired_ZeroTTL_NeverExpires(t *testing.T) {
	e := &Entry{CachedAt: time.Now().Add(-24 * time.Hour), TTL: 0}
	if e.IsExpired() {
		t.Error("zero TTL entry should never expire")
	}
}

func TestSet_PersistsAcrossReload(t *testing.T) {
	dir := tempDir(t)
	c1, _ := New(dir)
	_ = c1.Set("host2", "/etc/nginx/nginx.conf", "deadbeef", 0)

	c2, err := New(dir)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	e := c2.Get("host2", "/etc/nginx/nginx.conf")
	if e == nil {
		t.Fatal("expected persisted entry after reload")
	}
	if e.Checksum != "deadbeef" {
		t.Errorf("want deadbeef, got %s", e.Checksum)
	}
}
