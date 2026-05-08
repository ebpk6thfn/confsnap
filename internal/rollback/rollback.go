package rollback

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/confsnap/confsnap/internal/snapshot"
)

// Manager handles saving and restoring configuration snapshots for rollback.
type Manager struct {
	dir string
}

// Entry represents a saved rollback point.
type Entry struct {
	ID        string
	Host      string
	Path      string
	SavedAt   time.Time
	Checksum  string
	Content   []byte
}

// New creates a Manager that persists rollback entries under dir.
func New(dir string) (*Manager, error) {
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, fmt.Errorf("rollback: create dir: %w", err)
	}
	return &Manager{dir: dir}, nil
}

// Save records a snapshot as a rollback point and returns its ID.
func (m *Manager) Save(s *snapshot.Snapshot) (string, error) {
	id := entryID(s.Host, s.Path, time.Now())
	path := filepath.Join(m.dir, id)
	if err := os.WriteFile(path, s.Content, 0o600); err != nil {
		return "", fmt.Errorf("rollback: save entry: %w", err)
	}
	return id, nil
}

// Load retrieves a previously saved rollback entry by ID.
func (m *Manager) Load(id string) (*Entry, error) {
	path := filepath.Join(m.dir, id)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("rollback: load entry: %w", err)
	}
	host, filePath, savedAt, err := parseEntryID(id)
	if err != nil {
		return nil, fmt.Errorf("rollback: parse id: %w", err)
	}
	return &Entry{
		ID:      id,
		Host:    host,
		Path:    filePath,
		SavedAt: savedAt,
		Content: data,
	}, nil
}

// Delete removes a rollback entry by ID.
func (m *Manager) Delete(id string) error {
	path := filepath.Join(m.dir, id)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("rollback: delete entry: %w", err)
	}
	return nil
}

func entryID(host, path string, t time.Time) string {
	safe := func(s string) string {
		out := make([]byte, len(s))
		for i := range s {
			if s[i] == '/' || s[i] == ':' || s[i] == ' ' {
				out[i] = '_'
			} else {
				out[i] = s[i]
			}
		}
		return string(out)
	}
	return fmt.Sprintf("%s__%s__%s", safe(host), safe(path), t.UTC().Format("20060102T150405Z"))
}

func parseEntryID(id string) (host, path string, t time.Time, err error) {
	var ts string
	_, err = fmt.Sscanf(id, "%s", &id) // no-op, just reuse id
	// manual split on "__"
	parts := splitN(id, "__", 3)
	if len(parts) != 3 {
		return "", "", time.Time{}, fmt.Errorf("unexpected id format")
	}
	host = parts[0]
	path = parts[1]
	ts = parts[2]
	t, err = time.Parse("20060102T150405Z", ts)
	return host, path, t, err
}

func splitN(s, sep string, n int) []string {
	var parts []string
	for i := 0; i < n-1; i++ {
		idx := indexOf(s, sep)
		if idx < 0 {
			break
		}
		parts = append(parts, s[:idx])
		s = s[idx+len(sep):]
	}
	parts = append(parts, s)
	return parts
}

func indexOf(s, sub string) int {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
