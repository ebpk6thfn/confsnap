package report

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/confsnap/internal/diff"
)

func makeResults() []diff.Result {
	return []diff.Result{
		{Host: "web-01", Path: "/etc/nginx/nginx.conf", Changed: true, Diff: "- old\n+ new"},
		{Host: "web-01", Path: "/etc/hosts", Changed: false},
		{Host: "db-01", Path: "/etc/mysql/my.cnf", Changed: true, Diff: "- a\n+ b"},
	}
}

func TestNew_SetsResults(t *testing.T) {
	results := makeResults()
	rep := New(results)
	if len(rep.Results) != len(results) {
		t.Fatalf("expected %d results, got %d", len(results), len(rep.Results))
	}
	if rep.GeneratedAt.IsZero() {
		t.Error("expected GeneratedAt to be set")
	}
}

func TestSummary_Counts(t *testing.T) {
	rep := New(makeResults())
	changed, unchanged := rep.Summary()
	if changed != 2 {
		t.Errorf("expected 2 changed, got %d", changed)
	}
	if unchanged != 1 {
		t.Errorf("expected 1 unchanged, got %d", unchanged)
	}
}

func TestWrite_TextFormat_ContainsHosts(t *testing.T) {
	rep := New(makeResults())
	var buf bytes.Buffer
	if err := rep.Write(&buf, FormatText); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "web-01") {
		t.Error("expected output to contain 'web-01'")
	}
	if !strings.Contains(out, "Changed: 2") {
		t.Error("expected summary line with Changed: 2")
	}
}

func TestWrite_JSONFormat_ContainsFields(t *testing.T) {
	rep := New(makeResults())
	var buf bytes.Buffer
	if err := rep.Write(&buf, FormatJSON); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"changed":2`) {
		t.Errorf("expected JSON to contain changed count, got: %s", out)
	}
	if !strings.Contains(out, `"host":"db-01"`) {
		t.Error("expected JSON to contain db-01 host")
	}
}

func TestWrite_DefaultFormat_IsText(t *testing.T) {
	rep := New(makeResults())
	var buf bytes.Buffer
	if err := rep.Write(&buf, ""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "confsnap report") {
		t.Error("expected text report header")
	}
}
