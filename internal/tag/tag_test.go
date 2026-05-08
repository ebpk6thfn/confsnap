package tag

import (
	"testing"
)

func TestTagHost_And_HostTags(t *testing.T) {
	r := New()
	r.TagHost("web-01", "production", "web")
	tags := r.HostTags("web-01")
	if len(tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(tags))
	}
	if tags[0] != "production" || tags[1] != "web" {
		t.Errorf("unexpected tags: %v", tags)
	}
}

func TestTagHost_Deduplicates(t *testing.T) {
	r := New()
	r.TagHost("db-01", "db", "db")
	if got := len(r.HostTags("db-01")); got != 1 {
		t.Errorf("expected 1 unique tag, got %d", got)
	}
}

func TestTagFile_And_FileTags(t *testing.T) {
	r := New()
	r.TagFile("/etc/nginx/nginx.conf", "nginx", "web")
	tags := r.FileTags("/etc/nginx/nginx.conf")
	if len(tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(tags))
	}
}

func TestHostsWithTag(t *testing.T) {
	r := New()
	r.TagHost("web-01", "production")
	r.TagHost("web-02", "production", "canary")
	r.TagHost("db-01", "staging")

	hosts := r.HostsWithTag("production")
	if len(hosts) != 2 {
		t.Fatalf("expected 2 hosts, got %d", len(hosts))
	}
	if hosts[0] != "web-01" || hosts[1] != "web-02" {
		t.Errorf("unexpected hosts: %v", hosts)
	}
}

func TestHostsWithTag_NoMatch(t *testing.T) {
	r := New()
	r.TagHost("web-01", "production")
	if got := r.HostsWithTag("staging"); len(got) != 0 {
		t.Errorf("expected empty, got %v", got)
	}
}

func TestFilesWithTag(t *testing.T) {
	r := New()
	r.TagFile("/etc/nginx/nginx.conf", "nginx")
	r.TagFile("/etc/ssh/sshd_config", "security")
	r.TagFile("/etc/hosts", "nginx")

	files := r.FilesWithTag("nginx")
	if len(files) != 2 {
		t.Fatalf("expected 2 files, got %d", len(files))
	}
}

func TestHostTags_UnknownHost_ReturnsEmpty(t *testing.T) {
	r := New()
	if tags := r.HostTags("unknown"); len(tags) != 0 {
		t.Errorf("expected empty tags, got %v", tags)
	}
}

func TestValidate_Valid(t *testing.T) {
	if err := Validate([]string{"production", "web", "v2"}); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidate_EmptyTag(t *testing.T) {
	if err := Validate([]string{"ok", ""}); err == nil {
		t.Error("expected error for empty tag")
	}
}

func TestValidate_WhitespaceTag(t *testing.T) {
	if err := Validate([]string{"bad tag"}); err == nil {
		t.Error("expected error for tag with whitespace")
	}
}

func TestTagHost_MergesAcrossCalls(t *testing.T) {
	r := New()
	r.TagHost("web-01", "production")
	r.TagHost("web-01", "web")
	if got := len(r.HostTags("web-01")); got != 2 {
		t.Errorf("expected 2 tags after two calls, got %d", got)
	}
}
