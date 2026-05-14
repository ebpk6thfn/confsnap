package metrics

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestNewExporter_ValidFormats(t *testing.T) {
	r := New()
	for _, f := range []Format{FormatText, FormatJSON} {
		ex, err := NewExporter(r, f)
		if err != nil {
			t.Fatalf("format %q: unexpected error: %v", f, err)
		}
		if ex == nil {
			t.Fatalf("format %q: expected non-nil exporter", f)
		}
	}
}

func TestNewExporter_InvalidFormat(t *testing.T) {
	r := New()
	_, err := NewExporter(r, "xml")
	if err == nil {
		t.Fatal("expected error for unknown format")
	}
}

func TestWrite_TextFormat_ContainsCounterAndGauge(t *testing.T) {
	r := New()
	r.Counter("runs").Add(5)
	r.Gauge("drift").Set(2.5)

	ex, _ := NewExporter(r, FormatText)
	var buf bytes.Buffer
	if err := ex.Write(&buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "counter runs 5") {
		t.Errorf("expected counter line, got:\n%s", out)
	}
	if !strings.Contains(out, "gauge   drift 2.5") {
		t.Errorf("expected gauge line, got:\n%s", out)
	}
	if !strings.Contains(out, "# snapshot at ") {
		t.Errorf("expected header line, got:\n%s", out)
	}
}

func TestWrite_JSONFormat_ContainsFields(t *testing.T) {
	r := New()
	r.Counter("errors").Inc()
	r.Gauge("hosts").Set(3)

	ex, _ := NewExporter(r, FormatJSON)
	var buf bytes.Buffer
	if err := ex.Write(&buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var payload map[string]any
	if err := json.Unmarshal(buf.Bytes(), &payload); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if _, ok := payload["at"]; !ok {
		t.Error("expected 'at' field in JSON")
	}
	if _, ok := payload["counters"]; !ok {
		t.Error("expected 'counters' field in JSON")
	}
	if _, ok := payload["gauges"]; !ok {
		t.Error("expected 'gauges' field in JSON")
	}
}

func TestWrite_TextFormat_SortedKeys(t *testing.T) {
	r := New()
	r.Counter("z_last").Inc()
	r.Counter("a_first").Add(2)

	ex, _ := NewExporter(r, FormatText)
	var buf bytes.Buffer
	_ = ex.Write(&buf)
	out := buf.String()

	idxA := strings.Index(out, "a_first")
	idxZ := strings.Index(out, "z_last")
	if idxA > idxZ {
		t.Error("expected counters to be sorted alphabetically")
	}
}
