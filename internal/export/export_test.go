package export

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/confsnap/internal/diff"
)

func makeResult(host, path string, changed bool, added, removed int) diff.Result {
	return diff.Result{
		Host:         host,
		Path:         path,
		Changed:      changed,
		LinesAdded:   added,
		LinesRemoved: removed,
	}
}

func TestNew_ValidFormat(t *testing.T) {
	for _, f := range []Format{FormatCSV, FormatJSON} {
		_, err := New(f)
		if err != nil {
			t.Errorf("New(%q) unexpected error: %v", f, err)
		}
	}
}

func TestNew_InvalidFormat(t *testing.T) {
	_, err := New("xml")
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestWrite_CSVContainsHeader(t *testing.T) {
	ex, _ := New(FormatCSV)
	var buf bytes.Buffer
	if err := ex.Write(&buf, nil); err != nil {
		t.Fatalf("Write error: %v", err)
	}
	if !strings.Contains(buf.String(), "timestamp,host,path,status") {
		t.Errorf("CSV missing header, got: %s", buf.String())
	}
}

func TestWrite_CSVContainsRow(t *testing.T) {
	ex, _ := New(FormatCSV)
	ex.now = func() time.Time { return time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC) }
	results := []diff.Result{makeResult("web-01", "/etc/nginx.conf", true, 3, 1)}
	var buf bytes.Buffer
	if err := ex.Write(&buf, results); err != nil {
		t.Fatalf("Write error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"web-01", "/etc/nginx.conf", "changed", "3", "1"} {
		if !strings.Contains(out, want) {
			t.Errorf("CSV missing %q, got:\n%s", want, out)
		}
	}
}

func TestWrite_JSONContainsFields(t *testing.T) {
	ex, _ := New(FormatJSON)
	ex.now = func() time.Time { return time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC) }
	results := []diff.Result{makeResult("db-01", "/etc/my.cnf", false, 0, 0)}
	var buf bytes.Buffer
	if err := ex.Write(&buf, results); err != nil {
		t.Fatalf("Write error: %v", err)
	}
	var records []Record
	if err := json.Unmarshal(buf.Bytes(), &records); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}
	r := records[0]
	if r.Host != "db-01" || r.Path != "/etc/my.cnf" || r.Status != "unchanged" {
		t.Errorf("unexpected record: %+v", r)
	}
}

func TestWrite_JSON_EmptyResults(t *testing.T) {
	ex, _ := New(FormatJSON)
	var buf bytes.Buffer
	if err := ex.Write(&buf, []diff.Result{}); err != nil {
		t.Fatalf("Write error: %v", err)
	}
	if !strings.Contains(buf.String(), "[") {
		t.Errorf("expected JSON array, got: %s", buf.String())
	}
}
