package notify

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func makeEvent(level Level, host, file, msg string) Event {
	return Event{
		Timestamp: time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
		Level:     level,
		Host:      host,
		File:      file,
		Message:   msg,
	}
}

func TestNewLogger_NotNil(t *testing.T) {
	var buf bytes.Buffer
	l := NewLogger(&buf)
	if l == nil {
		t.Fatal("expected non-nil Logger")
	}
}

func TestSend_WritesFormattedLine(t *testing.T) {
	var buf bytes.Buffer
	l := NewLogger(&buf)
	e := makeEvent(LevelWarn, "web-01", "/etc/nginx/nginx.conf", "config changed")

	if err := l.Send([]Event{e}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "[WARN]") {
		t.Errorf("expected [WARN] in output, got: %s", out)
	}
	if !strings.Contains(out, "host=web-01") {
		t.Errorf("expected host=web-01 in output, got: %s", out)
	}
	if !strings.Contains(out, "file=/etc/nginx/nginx.conf") {
		t.Errorf("expected file path in output, got: %s", out)
	}
}

func TestSend_EmptyEvents_NoWrite(t *testing.T) {
	var buf bytes.Buffer
	l := NewLogger(&buf)
	if err := l.Send(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected no output for empty events, got %d bytes", buf.Len())
	}
}

func TestSend_MultipleEvents_MultipleLines(t *testing.T) {
	var buf bytes.Buffer
	l := NewLogger(&buf)
	events := []Event{
		makeEvent(LevelInfo, "db-01", "/etc/mysql/my.cnf", "no change"),
		makeEvent(LevelCrit, "db-02", "/etc/mysql/my.cnf", "critical drift"),
	}
	if err := l.Send(events); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 lines, got %d", len(lines))
	}
}

func TestFormatEvent_ContainsAllFields(t *testing.T) {
	e := makeEvent(LevelCrit, "app-01", "/etc/app/config.yaml", "drift detected")
	formatted := FormatEvent(e)

	for _, want := range []string{"[CRIT]", "host=app-01", "file=/etc/app/config.yaml", "drift detected", "2024-06-01"} {
		if !strings.Contains(formatted, want) {
			t.Errorf("expected %q in formatted output: %s", want, formatted)
		}
	}
}
