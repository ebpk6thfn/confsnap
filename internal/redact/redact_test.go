package redact

import (
	"strings"
	"testing"
)

func TestNew_ValidRules(t *testing.T) {
	r, err := New([]Rule{{Pattern: "password", Mask: ""}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r == nil {
		t.Fatal("expected non-nil Redactor")
	}
}

func TestNew_InvalidPattern(t *testing.T) {
	_, err := New([]Rule{{Pattern: "[", Mask: ""}})
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestApply_NoRules_ReturnsUnchanged(t *testing.T) {
	r, _ := New(nil)
	lines := []string{"password=secret", "host=localhost"}
	out := r.Apply(lines)
	for i, l := range out {
		if l != lines[i] {
			t.Errorf("line %d changed unexpectedly: got %q", i, l)
		}
	}
}

func TestApply_KeyValueEquals_Redacted(t *testing.T) {
	r, _ := New([]Rule{{Pattern: "(?i)password"}})
	out := r.Apply([]string{"password=supersecret"})
	if strings.Contains(out[0], "supersecret") {
		t.Errorf("secret not redacted: %q", out[0])
	}
	if !strings.Contains(out[0], defaultMask) {
		t.Errorf("expected mask in output: %q", out[0])
	}
}

func TestApply_KeyValueColon_Redacted(t *testing.T) {
	r, _ := New([]Rule{{Pattern: "api_key"}})
	out := r.Apply([]string{"api_key: abc123"})
	if strings.Contains(out[0], "abc123") {
		t.Errorf("secret not redacted: %q", out[0])
	}
}

func TestApply_CustomMask(t *testing.T) {
	r, _ := New([]Rule{{Pattern: "token", Mask: "<hidden>"}})
	out := r.Apply([]string{"token=mytoken"})
	if !strings.Contains(out[0], "<hidden>") {
		t.Errorf("expected custom mask, got: %q", out[0])
	}
}

func TestApply_NonMatchingKey_Unchanged(t *testing.T) {
	r, _ := New([]Rule{{Pattern: "password"}})
	line := "hostname=myserver"
	out := r.Apply([]string{line})
	if out[0] != line {
		t.Errorf("expected unchanged line, got: %q", out[0])
	}
}

func TestApplyString_RoundTrip(t *testing.T) {
	r, _ := New([]Rule{{Pattern: "secret"}})
	input := "host=localhost\nsecret=topsecret\nport=22"
	out := r.ApplyString(input)
	if strings.Contains(out, "topsecret") {
		t.Errorf("secret found in output: %q", out)
	}
	if !strings.Contains(out, "host=localhost") {
		t.Errorf("unrelated line altered: %q", out)
	}
}

func TestApply_LineWithoutSeparator_Unchanged(t *testing.T) {
	r, _ := New([]Rule{{Pattern: "password"}})
	line := "# this is a comment"
	out := r.Apply([]string{line})
	if out[0] != line {
		t.Errorf("comment line changed: %q", out[0])
	}
}
