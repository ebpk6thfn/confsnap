package notify

import (
	"fmt"
	"io"
	"strings"
	"time"
)

// Level represents the severity of a notification.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelCrit  Level = "CRIT"
)

// Event represents a single notification event.
type Event struct {
	Timestamp time.Time
	Level     Level
	Host      string
	File      string
	Message   string
}

// Sender is the interface for delivering notifications.
type Sender interface {
	Send(events []Event) error
}

// Logger writes notification events to an io.Writer in a structured text format.
type Logger struct {
	w io.Writer
}

// NewLogger creates a Logger that writes to w.
func NewLogger(w io.Writer) *Logger {
	return &Logger{w: w}
}

// Send writes each event as a formatted line to the underlying writer.
func (l *Logger) Send(events []Event) error {
	for _, e := range events {
		ts := e.Timestamp.UTC().Format(time.RFC3339)
		line := fmt.Sprintf("%s [%s] host=%s file=%s msg=%s\n",
			ts, e.Level, e.Host, e.File, e.Message)
		if _, err := io.WriteString(l.w, line); err != nil {
			return fmt.Errorf("notify: write failed: %w", err)
		}
	}
	return nil
}

// FormatEvent returns a human-readable single-line representation of an event.
func FormatEvent(e Event) string {
	parts := []string{
		e.Timestamp.UTC().Format(time.RFC3339),
		fmt.Sprintf("[%s]", e.Level),
		"host=" + e.Host,
		"file=" + e.File,
		"msg=" + e.Message,
	}
	return strings.Join(parts, " ")
}
