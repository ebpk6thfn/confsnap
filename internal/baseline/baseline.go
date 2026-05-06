package baseline

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/yourorg/confsnap/internal/snapshot"
)

// Baseline represents a saved reference point for configuration files.
type Baseline struct {
	Name      string                        `json:"name"`
	CreatedAt time.Time                     `json:"created_at"`
	Snapshots map[string]*snapshot.Snapshot `json:"snapshots"` // key: "host:path"
}

// Manager handles storing and loading baselines.
type Manager struct {
	dir string
}

// NewManager creates a Manager that persists baselines under dir.
func NewManager(dir string) (*Manager, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("baseline: create dir: %w", err)
	}
	return &Manager{dir: dir}, nil
}

// Save persists a named baseline to disk.
func (m *Manager) Save(b *Baseline) error {
	data, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return fmt.Errorf("baseline: marshal: %w", err)
	}
	path := m.filePath(b.Name)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("baseline: write %s: %w", path, err)
	}
	return nil
}

// Load retrieves a baseline by name. Returns nil, nil if not found.
func (m *Manager) Load(name string) (*Baseline, error) {
	path := m.filePath(name)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("baseline: read %s: %w", path, err)
	}
	var b Baseline
	if err := json.Unmarshal(data, &b); err != nil {
		return nil, fmt.Errorf("baseline: unmarshal: %w", err)
	}
	return &b, nil
}

// List returns the names of all saved baselines.
func (m *Manager) List() ([]string, error) {
	entries, err := os.ReadDir(m.dir)
	if err != nil {
		return nil, fmt.Errorf("baseline: list: %w", err)
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".json" {
			names = append(names, e.Name()[:len(e.Name())-5])
		}
	}
	return names, nil
}

func (m *Manager) filePath(name string) string {
	return filepath.Join(m.dir, name+".json")
}
