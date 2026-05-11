package template_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/confsnap/internal/diff"
	"github.com/yourorg/confsnap/internal/template"
)

func makeResults(changed bool) []diff.Result {
	return []diff.Result{
		{Host: "web-01", Path: "/etc/nginx.conf", Changed: changed},
		{Host: "web-02", Path: "/etc/hosts", Changed: false},
	}
}

func TestNew_EmptySource_ReturnsError(t *testing.T) {
	_, err := template.New("")
	if err == nil {
		t.Fatal("expected error for empty source")
	}
}

func TestNew_InvalidTemplate_ReturnsError(t *testing.T) {
	_, err := template.New("{{ .Unclosed")
	if err == nil {
		t.Fatal("expected parse error")
	}
}

func TestNew_ValidTemplate_NotNil(t *testing.T) {
	r, err := template.New("hosts: {{ .TotalHosts }}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r == nil {
		t.Fatal("expected non-nil renderer")
	}
}

func TestRender_ContainsExpectedFields(t *testing.T) {
	src := "hosts={{.TotalHosts}} files={{.TotalFiles}} changed={{.Changed}}"
	r, err := template.New(src)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	d := template.BuildData(makeResults(true))
	var buf bytes.Buffer
	if err := r.Render(&buf, d); err != nil {
		t.Fatalf("Render: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "hosts=2") {
		t.Errorf("expected hosts=2, got %q", out)
	}
	if !strings.Contains(out, "files=2") {
		t.Errorf("expected files=2, got %q", out)
	}
	if !strings.Contains(out, "changed=1") {
		t.Errorf("expected changed=1, got %q", out)
	}
}

func TestBuildData_CountsHosts(t *testing.T) {
	d := template.BuildData(makeResults(false))
	if d.TotalHosts != 2 {
		t.Errorf("expected 2 hosts, got %d", d.TotalHosts)
	}
	if d.Unchanged != 2 {
		t.Errorf("expected 2 unchanged, got %d", d.Unchanged)
	}
	if d.Changed != 0 {
		t.Errorf("expected 0 changed, got %d", d.Changed)
	}
}

func TestRender_FuncMap_FormatTime(t *testing.T) {
	r, err := template.New("at={{ .GeneratedAt | formatTime }}")
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	d := template.BuildData(nil)
	var buf bytes.Buffer
	if err := r.Render(&buf, d); err != nil {
		t.Fatalf("Render: %v", err)
	}
	if !strings.HasPrefix(buf.String(), "at=") {
		t.Errorf("unexpected output: %q", buf.String())
	}
}
