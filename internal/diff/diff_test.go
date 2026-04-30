package diff

import (
	"strings"
	"testing"
	"time"

	"github.com/user/confsnap/internal/snapshot"
)

func makeSnapshot(host, path, content string) *snapshot.Snapshot {
	return snapshot.New(host, path, content, time.Now())
}

func TestCompare_NoChange(t *testing.T) {
	old := makeSnapshot("web01", "/etc/nginx.conf", "worker_processes 4;")
	new := makeSnapshot("web01", "/etc/nginx.conf", "worker_processes 4;")

	result := Compare(old, new)

	if result.Changed {
		t.Error("expected no change")
	}
	if len(result.Lines) != 0 {
		t.Errorf("expected 0 diff lines, got %d", len(result.Lines))
	}
}

func TestCompare_WithChange(t *testing.T) {
	old := makeSnapshot("web01", "/etc/nginx.conf", "worker_processes 4;")
	new := makeSnapshot("web01", "/etc/nginx.conf", "worker_processes 8;")

	result := Compare(old, new)

	if !result.Changed {
		t.Error("expected change to be detected")
	}
	if result.OldHash == result.NewHash {
		t.Error("expected hashes to differ")
	}
	if len(result.Lines) == 0 {
		t.Error("expected diff lines")
	}
}

func TestCompare_HostAndPathSet(t *testing.T) {
	old := makeSnapshot("db01", "/etc/my.cnf", "max_connections=100")
	new := makeSnapshot("db01", "/etc/my.cnf", "max_connections=200")

	result := Compare(old, new)

	if result.Host != "db01" {
		t.Errorf("expected host db01, got %s", result.Host)
	}
	if result.Path != "/etc/my.cnf" {
		t.Errorf("expected path /etc/my.cnf, got %s", result.Path)
	}
}

func TestResult_Format_NoChange(t *testing.T) {
	old := makeSnapshot("web01", "/etc/hosts", "127.0.0.1 localhost")
	new := makeSnapshot("web01", "/etc/hosts", "127.0.0.1 localhost")

	result := Compare(old, new)
	output := result.Format()

	if !strings.Contains(output, "no change") {
		t.Errorf("expected 'no change' in output, got: %s", output)
	}
}

func TestResult_Format_WithChange(t *testing.T) {
	old := makeSnapshot("web01", "/etc/hosts", "127.0.0.1 localhost")
	new := makeSnapshot("web01", "/etc/hosts", "127.0.0.1 newhost")

	result := Compare(old, new)
	output := result.Format()

	if !strings.Contains(output, "---") || !strings.Contains(output, "++") {
		t.Errorf("expected unified diff header in output, got: %s", output)
	}
}

func TestOp_String(t *testing.T) {
	if OpAdd.String() != "+" {
		t.Errorf("expected '+', got %s", OpAdd.String())
	}
	if OpRemove.String() != "-" {
		t.Errorf("expected '-', got %s", OpRemove.String())
	}
	if OpEqual.String() != " " {
		t.Errorf("expected ' ', got %s", OpEqual.String())
	}
}
