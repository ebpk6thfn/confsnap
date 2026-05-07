package plugin_test

import (
	"errors"
	"testing"

	"github.com/yourorg/confsnap/internal/plugin"
)

func TestNew_NotNil(t *testing.T) {
	r := plugin.New()
	if r == nil {
		t.Fatal("expected non-nil registry")
	}
}

func TestRegister_ThenNames(t *testing.T) {
	r := plugin.New()
	if err := r.Register(plugin.EventSnapshotTaken, "logger", func(e plugin.Event) error { return nil }); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	names := r.Names(plugin.EventSnapshotTaken)
	if len(names) != 1 || names[0] != "logger" {
		t.Fatalf("expected [logger], got %v", names)
	}
}

func TestRegister_DuplicateName_ReturnsError(t *testing.T) {
	r := plugin.New()
	h := func(e plugin.Event) error { return nil }
	_ = r.Register(plugin.EventSnapshotTaken, "dup", h)
	err := r.Register(plugin.EventSnapshotTaken, "dup", h)
	if err == nil {
		t.Fatal("expected error for duplicate name")
	}
}

func TestEmit_CallsHandlers(t *testing.T) {
	r := plugin.New()
	called := 0
	_ = r.Register(plugin.EventDiffDetected, "counter", func(e plugin.Event) error {
		called++
		return nil
	})
	e := plugin.Event{Type: plugin.EventDiffDetected, Host: "host1", Path: "/etc/hosts"}
	if err := r.Emit(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called != 1 {
		t.Fatalf("expected handler called once, got %d", called)
	}
}

func TestEmit_NoHandlers_ReturnsNil(t *testing.T) {
	r := plugin.New()
	if err := r.Emit(plugin.Event{Type: plugin.EventRunComplete}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestEmit_HandlerError_ReturnsError(t *testing.T) {
	r := plugin.New()
	_ = r.Register(plugin.EventSnapshotTaken, "failing", func(e plugin.Event) error {
		return errors.New("boom")
	})
	err := r.Emit(plugin.Event{Type: plugin.EventSnapshotTaken})
	if err == nil {
		t.Fatal("expected error from failing handler")
	}
}

func TestEmit_MultipleHandlers_AllCalled(t *testing.T) {
	r := plugin.New()
	var calls []string
	for _, name := range []string{"a", "b", "c"} {
		n := name
		_ = r.Register(plugin.EventRunComplete, n, func(e plugin.Event) error {
			calls = append(calls, n)
			return nil
		})
	}
	_ = r.Emit(plugin.Event{Type: plugin.EventRunComplete})
	if len(calls) != 3 {
		t.Fatalf("expected 3 calls, got %d", len(calls))
	}
}

func TestEmit_MetaPropagated(t *testing.T) {
	r := plugin.New()
	var received plugin.Event
	_ = r.Register(plugin.EventDiffDetected, "spy", func(e plugin.Event) error {
		received = e
		return nil
	})
	e := plugin.Event{
		Type: plugin.EventDiffDetected,
		Host: "web01",
		Path: "/etc/nginx/nginx.conf",
		Meta: map[string]string{"lines_changed": "4"},
	}
	_ = r.Emit(e)
	if received.Meta["lines_changed"] != "4" {
		t.Fatalf("expected meta propagated, got %v", received.Meta)
	}
}
