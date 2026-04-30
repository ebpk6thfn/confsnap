package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Store persists and retrieves snapshots from a local directory.
type Store struct {
	BaseDir string
}

// NewStore creates a Store rooted at baseDir, creating the directory if needed.
func NewStore(baseDir string) (*Store, error) {
	if err := os.MkdirAll(baseDir, 0o755); err != nil {
		return nil, fmt.Errorf("creating store directory: %w", err)
	}
	return &Store{BaseDir: baseDir}, nil
}

// Save writes a snapshot to disk as a JSON file.
func (st *Store) Save(s *Snapshot) error {
	fileName := snapshotFileName(s.Host, s.FilePath, s.CapturedAt)
	path := filepath.Join(st.BaseDir, fileName)

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("marshalling snapshot: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("writing snapshot file: %w", err)
	}
	return nil
}

// Load reads a snapshot from a file path.
func (st *Store) Load(path string) (*Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading snapshot file: %w", err)
	}
	var s Snapshot
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("unmarshalling snapshot: %w", err)
	}
	return &s, nil
}

// snapshotFileName generates a safe file name for a snapshot.
func snapshotFileName(host, filePath string, t time.Time) string {
	safePath := strings.NewReplacer("/", "_", ".", "_").Replace(filePath)
	timestamp := t.Format("20060102T150405Z")
	return fmt.Sprintf("%s%s_%s.json", host, safePath, timestamp)
}
