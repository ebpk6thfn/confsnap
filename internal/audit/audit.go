package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Entry represents a single audit log record.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Host      string    `json:"host"`
	File      string    `json:"file"`
	Event     string    `json:"event"`
	Details   string    `json:"details,omitempty"`
}

// Logger writes audit entries to a JSON-lines file.
type Logger struct {
	dir string
}

// New creates a Logger that stores audit logs under dir.
func New(dir string) (*Logger, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("audit: create log dir: %w", err)
	}
	return &Logger{dir: dir}, nil
}

// Log appends an entry to today's audit log file.
func (l *Logger) Log(host, file, event, details string) error {
	entry := Entry{
		Timestamp: time.Now().UTC(),
		Host:      host,
		File:      file,
		Event:     event,
		Details:   details,
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("audit: marshal entry: %w", err)
	}

	path := l.logPath(entry.Timestamp)
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("audit: open log file: %w", err)
	}
	defer f.Close()

	_, err = fmt.Fprintf(f, "%s\n", data)
	return err
}

// ReadDay returns all entries logged on the given date (UTC).
func (l *Logger) ReadDay(t time.Time) ([]Entry, error) {
	path := l.logPath(t)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("audit: read log file: %w", err)
	}

	var entries []Entry
	for _, line := range splitLines(data) {
		if len(line) == 0 {
			continue
		}
		var e Entry
		if err := json.Unmarshal(line, &e); err != nil {
			return nil, fmt.Errorf("audit: parse entry: %w", err)
		}
		entries = append(entries, e)
	}
	return entries, nil
}

func (l *Logger) logPath(t time.Time) string {
	name := fmt.Sprintf("audit-%s.jsonl", t.UTC().Format("2006-01-02"))
	return filepath.Join(l.dir, name)
}

func splitLines(data []byte) [][]byte {
	var lines [][]byte
	start := 0
	for i, b := range data {
		if b == '\n' {
			lines = append(lines, data[start:i])
			start = i + 1
		}
	}
	if start < len(data) {
		lines = append(lines, data[start:])
	}
	return lines
}
