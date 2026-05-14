package healthcheck

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Store persists health-check results to disk as JSON files.
type Store struct {
	dir string
}

// NewStore creates a Store rooted at dir, creating the directory if needed.
func NewStore(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("healthcheck store: mkdir %s: %w", dir, err)
	}
	return &Store{dir: dir}, nil
}

// Save writes results to a timestamped JSON file.
func (s *Store) Save(results []Result) error {
	name := fmt.Sprintf("%s.json", time.Now().UTC().Format("20060102T150405Z"))
	path := filepath.Join(s.dir, name)
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("healthcheck store: create %s: %w", path, err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(results); err != nil {
		return fmt.Errorf("healthcheck store: encode: %w", err)
	}
	return nil
}

// Latest returns the most recently saved results, or nil if none exist.
func (s *Store) Latest() ([]Result, error) {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		return nil, fmt.Errorf("healthcheck store: readdir: %w", err)
	}
	if len(entries) == 0 {
		return nil, nil
	}
	latest := entries[len(entries)-1]
	path := filepath.Join(s.dir, latest.Name())
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("healthcheck store: read %s: %w", path, err)
	}
	var results []Result
	if err := json.Unmarshal(data, &results); err != nil {
		return nil, fmt.Errorf("healthcheck store: unmarshal: %w", err)
	}
	return results, nil
}
