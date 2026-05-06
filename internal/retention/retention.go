package retention

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Policy defines how long snapshots are retained.
type Policy struct {
	// MaxAge is the maximum age of a snapshot before it is pruned.
	MaxAge time.Duration
	// MaxCount is the maximum number of snapshots to keep per host+path.
	// Zero means unlimited.
	MaxCount int
}

// Manager applies a retention policy to a snapshot store directory.
type Manager struct {
	policy  Policy
	storeDir string
}

// New creates a new retention Manager.
func New(storeDir string, policy Policy) (*Manager, error) {
	if storeDir == "" {
		return nil, fmt.Errorf("retention: storeDir must not be empty")
	}
	if policy.MaxAge < 0 {
		return nil, fmt.Errorf("retention: MaxAge must be non-negative")
	}
	if policy.MaxCount < 0 {
		return nil, fmt.Errorf("retention: MaxCount must be non-negative")
	}
	return &Manager{policy: policy, storeDir: storeDir}, nil
}

// Prune removes snapshot files that violate the retention policy.
// It returns the list of removed file paths.
func (m *Manager) Prune(now time.Time) ([]string, error) {
	entries, err := os.ReadDir(m.storeDir)
	if err != nil {
		return nil, fmt.Errorf("retention: reading store dir: %w", err)
	}

	// Group files by key (host+path prefix).
	groups := map[string][]os.DirEntry{}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		key := groupKey(e.Name())
		groups[key] = append(groups[key], e)
	}

	var removed []string
	for _, files := range groups {
		// Sort oldest first by name (names contain timestamps).
		sort.Slice(files, func(i, j int) bool {
			return files[i].Name() < files[j].Name()
		})
		pruned, err := m.pruneGroup(files, now)
		if err != nil {
			return removed, err
		}
		removed = append(removed, pruned...)
	}
	return removed, nil
}

func (m *Manager) pruneGroup(files []os.DirEntry, now time.Time) ([]string, error) {
	var removed []string
	for i, f := range files {
		path := filepath.Join(m.storeDir, f.Name())
		info, err := f.Info()
		if err != nil {
			return removed, fmt.Errorf("retention: stat %s: %w", f.Name(), err)
		}
		ageExceeded := m.policy.MaxAge > 0 && now.Sub(info.ModTime()) > m.policy.MaxAge
		countExceeded := m.policy.MaxCount > 0 && i < len(files)-m.policy.MaxCount
		if ageExceeded || countExceeded {
			if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
				return removed, fmt.Errorf("retention: removing %s: %w", path, err)
			}
			removed = append(removed, path)
		}
	}
	return removed, nil
}

// groupKey derives a stable group identifier from a snapshot filename.
func groupKey(name string) string {
	// Filenames are expected to be: <host>_<escapedpath>_<timestamp>.json
	// We group by everything before the last underscore-separated timestamp.
	parts := strings.Split(strings.TrimSuffix(name, ".json"), "_")
	if len(parts) <= 1 {
		return name
	}
	return strings.Join(parts[:len(parts)-1], "_")
}
