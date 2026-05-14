package metrics

import (
	"testing"
)

func TestNew_NotNil(t *testing.T) {
	r := New()
	if r == nil {
		t.Fatal("expected non-nil registry")
	}
}

func TestCounter_IncAndValue(t *testing.T) {
	r := New()
	c := r.Counter("ssh.connections")
	c.Inc()
	c.Inc()
	if c.Value() != 2 {
		t.Fatalf("expected 2, got %d", c.Value())
	}
}

func TestCounter_Add(t *testing.T) {
	r := New()
	c := r.Counter("files.scanned")
	c.Add(10)
	c.Add(5)
	if c.Value() != 15 {
		t.Fatalf("expected 15, got %d", c.Value())
	}
}

func TestCounter_SameName_ReturnsSameInstance(t *testing.T) {
	r := New()
	a := r.Counter("errors")
	b := r.Counter("errors")
	a.Inc()
	if b.Value() != 1 {
		t.Fatal("expected same counter instance")
	}
}

func TestGauge_SetAndValue(t *testing.T) {
	r := New()
	g := r.Gauge("drift.score")
	g.Set(3.14)
	if g.Value() != 3.14 {
		t.Fatalf("expected 3.14, got %f", g.Value())
	}
}

func TestGauge_SameName_ReturnsSameInstance(t *testing.T) {
	r := New()
	a := r.Gauge("hosts.reachable")
	b := r.Gauge("hosts.reachable")
	a.Set(7)
	if b.Value() != 7 {
		t.Fatal("expected same gauge instance")
	}
}

func TestSnapshot_ContainsValues(t *testing.T) {
	r := New()
	r.Counter("runs").Add(3)
	r.Gauge("latency").Set(1.5)

	s := r.Snapshot()

	if s.Counters["runs"] != 3 {
		t.Fatalf("expected runs=3, got %d", s.Counters["runs"])
	}
	if s.Gauges["latency"] != 1.5 {
		t.Fatalf("expected latency=1.5, got %f", s.Gauges["latency"])
	}
	if s.At.IsZero() {
		t.Fatal("expected non-zero snapshot time")
	}
}

func TestSnapshot_IsCopy(t *testing.T) {
	r := New()
	c := r.Counter("x")
	c.Add(1)
	s := r.Snapshot()
	c.Add(99)
	if s.Counters["x"] != 1 {
		t.Fatal("snapshot should not reflect mutations after capture")
	}
}
