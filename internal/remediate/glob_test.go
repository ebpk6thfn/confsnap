package remediate

import "testing"

func TestMatchGlob_ExactMatch(t *testing.T) {
	if !matchGlob("/etc/hosts", "/etc/hosts") {
		t.Error("expected exact match")
	}
}

func TestMatchGlob_WildcardMatch(t *testing.T) {
	if !matchGlob("/etc/*", "/etc/hosts") {
		t.Error("expected wildcard match")
	}
}

func TestMatchGlob_NoMatch(t *testing.T) {
	if matchGlob("/etc/nginx/*", "/etc/ssh/sshd_config") {
		t.Error("expected no match")
	}
}

func TestMatchGlob_InvalidPattern(t *testing.T) {
	if matchGlob("[", "/etc/hosts") {
		t.Error("invalid pattern should not match")
	}
}

func TestMatchGlob_StarStar_NotSupported(t *testing.T) {
	// path.Match does not support ** so this should not match nested paths
	if matchGlob("/etc/**", "/etc/ssh/sshd_config") {
		t.Error("** glob should not match nested path with path.Match")
	}
}
