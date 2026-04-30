package snapshot

import (
	"testing"
)

func TestNew_SetsFields(t *testing.T) {
	s := New("web01", "/etc/nginx/nginx.conf", "worker_processes 4;")

	if s.Host != "web01" {
		t.Errorf("expected host 'web01', got '%s'", s.Host)
	}
	if s.FilePath != "/etc/nginx/nginx.conf" {
		t.Errorf("expected file path '/etc/nginx/nginx.conf', got '%s'", s.FilePath)
	}
	if s.Content != "worker_processes 4;" {
		t.Errorf("unexpected content: %s", s.Content)
	}
	if s.Checksum == "" {
		t.Error("expected non-empty checksum")
	}
	if s.CapturedAt.IsZero() {
		t.Error("expected CapturedAt to be set")
	}
}

func TestComputeChecksum_Deterministic(t *testing.T) {
	c1 := computeChecksum("hello")
	c2 := computeChecksum("hello")
	if c1 != c2 {
		t.Errorf("checksums should be equal: %s != %s", c1, c2)
	}
}

func TestComputeChecksum_DifferentInputs(t *testing.T) {
	c1 := computeChecksum("hello")
	c2 := computeChecksum("world")
	if c1 == c2 {
		t.Error("different inputs should produce different checksums")
	}
}

func TestEqual_SameContent(t *testing.T) {
	s1 := New("host1", "/etc/hosts", "127.0.0.1 localhost")
	s2 := New("host2", "/etc/hosts", "127.0.0.1 localhost")
	if !s1.Equal(s2) {
		t.Error("snapshots with same content should be equal")
	}
}

func TestEqual_DifferentContent(t *testing.T) {
	s1 := New("host1", "/etc/hosts", "127.0.0.1 localhost")
	s2 := New("host1", "/etc/hosts", "192.168.1.1 myhost")
	if s1.Equal(s2) {
		t.Error("snapshots with different content should not be equal")
	}
}
