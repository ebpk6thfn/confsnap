package alert

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/confsnap/internal/diff"
)

func makeResult(host, path string, changed bool, diffLines []string) diff.Result {
	return diff.Result{
		Host:    host,
		Path:    path,
		Changed: changed,
		Diff:    diffLines,
	}
}

func TestEvaluate_NoChanges_ReturnsEmpty(t *testing.T) {
	a := New(&bytes.Buffer{}, 3)
	results := []diff.Result{makeResult("host1", "/etc/hosts", false, nil)}
	alerts := a.Evaluate(results)
	if len(alerts) != 0 {
		t.Fatalf("expected 0 alerts, got %d", len(alerts))
	}
}

func TestEvaluate_BelowThreshold_InfoLevel(t *testing.T) {
	a := New(&bytes.Buffer{}, 5)
	results := []diff.Result{
		makeResult("host1", "/etc/hosts", true, []string{"+newline"}),
	}
	alerts := a.Evaluate(results)
	if len(alerts) != 1 {
		t.Fatalf("expected 1 alert, got %d", len(alerts))
	}
	if alerts[0].Level != LevelInfo {
		t.Errorf("expected info, got %s", alerts[0].Level)
	}
}

func TestEvaluate_AtThreshold_WarnLevel(t *testing.T) {
	a := New(&bytes.Buffer{}, 2)
	diffLines := []string{"+a", "-b"}
	results := []diff.Result{makeResult("host2", "/etc/ssh/sshd_config", true, diffLines)}
	alerts := a.Evaluate(results)
	if alerts[0].Level != LevelWarn {
		t.Errorf("expected warn, got %s", alerts[0].Level)
	}
}

func TestEvaluate_DoubleThreshold_CritLevel(t *testing.T) {
	a := New(&bytes.Buffer{}, 2)
	diffLines := []string{"+a", "-b", "+c", "-d"}
	results := []diff.Result{makeResult("host3", "/etc/passwd", true, diffLines)}
	alerts := a.Evaluate(results)
	if alerts[0].Level != LevelCrit {
		t.Errorf("expected crit, got %s", alerts[0].Level)
	}
}

func TestWrite_FormatsCorrectly(t *testing.T) {
	var buf bytes.Buffer
	a := New(&buf, 1)
	alerts := []Alert{
		{Host: "web01", Path: "/etc/nginx/nginx.conf", Level: LevelWarn, Message: "2 line(s) changed"},
	}
	if err := a.Write(alerts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "[WARN]") {
		t.Errorf("expected [WARN] in output, got: %s", out)
	}
	if !strings.Contains(out, "web01") {
		t.Errorf("expected host in output, got: %s", out)
	}
}

func TestNew_DefaultThreshold(t *testing.T) {
	a := New(&bytes.Buffer{}, 0)
	if a.threshold != 1 {
		t.Errorf("expected threshold 1, got %d", a.threshold)
	}
}
