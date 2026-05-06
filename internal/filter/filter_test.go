package filter_test

import (
	"testing"

	"github.com/yourorg/confsnap/internal/filter"
)

func TestAllow_NoRules_AllowsAll(t *testing.T) {
	f := filter.New(nil, nil)
	paths := []string{"/etc/hosts", "/etc/nginx/nginx.conf", "/var/log/syslog"}
	for _, p := range paths {
		if !f.Allow(p) {
			t.Errorf("expected %q to be allowed with no rules", p)
		}
	}
}

func TestAllow_IncludePattern_FiltersOthers(t *testing.T) {
	f := filter.New([]string{"/etc/*"}, nil)

	if !f.Allow("/etc/hosts") {
		t.Error("expected /etc/hosts to be allowed")
	}
	if f.Allow("/var/log/syslog") {
		t.Error("expected /var/log/syslog to be excluded")
	}
}

func TestAllow_ExcludePattern_BlocksMatch(t *testing.T) {
	f := filter.New(nil, []string{"*.bak"})

	if f.Allow("/etc/hosts.bak") {
		t.Error("expected /etc/hosts.bak to be blocked by exclude rule")
	}
	if !f.Allow("/etc/hosts") {
		t.Error("expected /etc/hosts to be allowed")
	}
}

func TestAllow_ExcludeTakesPrecedence(t *testing.T) {
	f := filter.New([]string{"/etc/*"}, []string{"/etc/shadow"})

	if !f.Allow("/etc/hosts") {
		t.Error("expected /etc/hosts to be allowed")
	}
	if f.Allow("/etc/shadow") {
		t.Error("expected /etc/shadow to be excluded despite include rule")
	}
}

func TestAllow_DirectoryPrefix(t *testing.T) {
	f := filter.New(nil, []string{"/etc/ssl/"})

	if f.Allow("/etc/ssl/private/key.pem") {
		t.Error("expected /etc/ssl/private/key.pem to be excluded by directory prefix")
	}
	if !f.Allow("/etc/hosts") {
		t.Error("expected /etc/hosts to be allowed")
	}
}

func TestAllow_BasenameMatch(t *testing.T) {
	f := filter.New(nil, []string{"*.log"})

	if f.Allow("/var/log/app.log") {
		t.Error("expected /var/log/app.log to be excluded via basename match")
	}
	if !f.Allow("/etc/nginx/nginx.conf") {
		t.Error("expected /etc/nginx/nginx.conf to be allowed")
	}
}
