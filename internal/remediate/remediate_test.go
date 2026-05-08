package remediate_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/confsnap/internal/diff"
	"github.com/yourorg/confsnap/internal/remediate"
	"github.com/yourorg/confsnap/internal/rollback"
	"github.com/yourorg/confsnap/internal/snapshot"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "remediate-test-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func makeResult(host, path string, changed bool) diff.Result {
	return diff.Result{Host: host, Path: path, Changed: changed}
}

func TestEvaluate_NoChanges_ReturnsEmpty(t *testing.T) {
	dir := tempDir(t)
	store, _ := rollback.New(filepath.Join(dir, "rb"))
	rm := remediate.New([]remediate.Rule{{Pattern: "*", Action: remediate.ActionNotify}}, store)
	results := rm.Evaluate([]diff.Result{makeResult("host1", "/etc/foo", false)})
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestEvaluate_NoMatchingRule_ReturnsEmpty(t *testing.T) {
	dir := tempDir(t)
	store, _ := rollback.New(filepath.Join(dir, "rb"))
	rm := remediate.New([]remediate.Rule{{Pattern: "/etc/nginx/*", Action: remediate.ActionNotify}}, store)
	results := rm.Evaluate([]diff.Result{makeResult("host1", "/etc/ssh/sshd_config", true)})
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestEvaluate_NotifyAction_Applied(t *testing.T) {
	dir := tempDir(t)
	store, _ := rollback.New(filepath.Join(dir, "rb"))
	rm := remediate.New([]remediate.Rule{{Pattern: "/etc/*", Action: remediate.ActionNotify}}, store)
	results := rm.Evaluate([]diff.Result{makeResult("host1", "/etc/hosts", true)})
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if !results[0].Applied {
		t.Error("expected Applied=true for notify action")
	}
	if results[0].Action != remediate.ActionNotify {
		t.Errorf("unexpected action: %s", results[0].Action)
	}
}

func TestEvaluate_RollbackAction_NoEntry(t *testing.T) {
	dir := tempDir(t)
	store, _ := rollback.New(filepath.Join(dir, "rb"))
	rm := remediate.New([]remediate.Rule{{Pattern: "/etc/*", Action: remediate.ActionRollback}}, store)
	results := rm.Evaluate([]diff.Result{makeResult("host1", "/etc/hosts", true)})
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Applied {
		t.Error("expected Applied=false when no rollback entry exists")
	}
}

func TestEvaluate_RollbackAction_WithEntry(t *testing.T) {
	dir := tempDir(t)
	store, _ := rollback.New(filepath.Join(dir, "rb"))
	snap := snapshot.New("host1", "/etc/hosts", []byte("127.0.0.1 localhost"))
	_ = store.Save("host1", "/etc/hosts", snap)
	rm := remediate.New([]remediate.Rule{{Pattern: "/etc/*", Action: remediate.ActionRollback}}, store)
	results := rm.Evaluate([]diff.Result{makeResult("host1", "/etc/hosts", true)})
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if !results[0].Applied {
		t.Errorf("expected Applied=true, message: %s", results[0].Message)
	}
}
