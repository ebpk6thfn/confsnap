package ratelimit

import (
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	c := DefaultConfig()
	if c.MaxPerHost <= 0 {
		t.Errorf("expected MaxPerHost > 0, got %d", c.MaxPerHost)
	}
	if c.IntervalMs < 0 {
		t.Errorf("expected IntervalMs >= 0, got %d", c.IntervalMs)
	}
}

func TestToPolicy_Valid(t *testing.T) {
	c := Config{MaxPerHost: 4, IntervalMs: 200}
	p, err := c.ToPolicy()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.MaxPerHost != 4 {
		t.Errorf("expected MaxPerHost=4, got %d", p.MaxPerHost)
	}
	if p.Interval.Milliseconds() != 200 {
		t.Errorf("expected Interval=200ms, got %v", p.Interval)
	}
}

func TestToPolicy_ZeroMaxPerHost(t *testing.T) {
	c := Config{MaxPerHost: 0, IntervalMs: 100}
	_, err := c.ToPolicy()
	if err == nil {
		t.Fatal("expected error for MaxPerHost=0")
	}
}

func TestToPolicy_NegativeInterval(t *testing.T) {
	c := Config{MaxPerHost: 1, IntervalMs: -1}
	_, err := c.ToPolicy()
	if err == nil {
		t.Fatal("expected error for negative IntervalMs")
	}
}

func TestToPolicy_ZeroInterval_Allowed(t *testing.T) {
	c := Config{MaxPerHost: 2, IntervalMs: 0}
	p, err := c.ToPolicy()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Interval != 0 {
		t.Errorf("expected zero Interval, got %v", p.Interval)
	}
}

func TestToPolicy_RoundTrip_DefaultConfig(t *testing.T) {
	c := DefaultConfig()
	p, err := c.ToPolicy()
	if err != nil {
		t.Fatalf("unexpected error converting DefaultConfig: %v", err)
	}
	default_ := DefaultPolicy()
	if p.MaxPerHost != default_.MaxPerHost {
		t.Errorf("MaxPerHost mismatch: got %d, want %d", p.MaxPerHost, default_.MaxPerHost)
	}
	if p.Interval != default_.Interval {
		t.Errorf("Interval mismatch: got %v, want %v", p.Interval, default_.Interval)
	}
}
