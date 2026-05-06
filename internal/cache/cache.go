package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Entry holds a cached snapshot result with metadata.
type Entry struct {
	Host      string    `json:"host"`
	Path      string    `json:"path"`
	Checksum  string    `json:"checksum"`
	CachedAt  time.Time `json:"cached_at"`
	TTL       int64     `json:"ttl_seconds"`
}

// IsExpired reports whether the cache entry has passed its TTL.
func (e *Entry) IsExpired() bool {
	if e.TTL <= 0 {
		return false
	}
	return time.Since(e.CachedAt) > time.Duration(e.TTL)*time.Second
}

// Cache is a simple file-backed checksum cache to avoid redundant SSH reads.
type Cache struct {
	mu      sync.RWMutex
	dir     string
	entries map[string]*Entry
}

// New creates a Cache backed by dir, loading any existing entries from disk.
func New(dir string) (*Cache, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("cache: create dir: %w", err)
	}
	c := &Cache{dir: dir, entries: make(map[string]*Entry)}
	_ = c.load() // best-effort load; ignore missing file
	return c, nil
}

func cacheKey(host, path string) string {
	return host + "::" + path
}

// Get returns the cached entry for host+path, or nil if absent or expired.
func (c *Cache) Get(host, path string) *Entry {
	c.mu.RLock()
	defer c.mu.RUnlock()
	e, ok := c.entries[cacheKey(host, path)]
	if !ok || e.IsExpired() {
		return nil
	}
	return e
}

// Set stores an entry for host+path and persists it to disk.
func (c *Cache) Set(host, path, checksum string, ttl int64) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[cacheKey(host, path)] = &Entry{
		Host:     host,
		Path:     path,
		Checksum: checksum,
		CachedAt: time.Now(),
		TTL:      ttl,
	}
	return c.persist()
}

func (c *Cache) cacheFile() string {
	return filepath.Join(c.dir, "cache.json")
}

func (c *Cache) load() error {
	data, err := os.ReadFile(c.cacheFile())
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &c.entries)
}

func (c *Cache) persist() error {
	data, err := json.MarshalIndent(c.entries, "", "  ")
	if err != nil {
		return fmt.Errorf("cache: marshal: %w", err)
	}
	return os.WriteFile(c.cacheFile(), data, 0o644)
}
