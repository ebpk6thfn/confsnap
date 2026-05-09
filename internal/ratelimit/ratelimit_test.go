package ratelimit

import (
	"testing"
	"time"
)

func TestNew_ValidPolicy(t *testing.T) {
	l, err := New(DefaultPolicy())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l == nil {
		t.Fatal("expected non-nil Limiter")
	}
}

func TestNew_InvalidMaxPerHost(t *testing.T) {
	_, err := New(Policy{MaxPerHost: 0, Interval: time.Second})
	if err == nil {
		t.Fatal("expected error for MaxPerHost=0")
	}
}

func TestNew_NegativeInterval(t *testing.T) {
	_, err := New(Policy{MaxPerHost: 1, Interval: -time.Second})
	if err == nil {
		t.Fatal("expected error for negative Interval")
	}
}

func TestAcquire_Release_Basic(t *testing.T) {
	l, _ := New(Policy{MaxPerHost: 2, Interval: 0})

	if err := l.Acquire("host1"); err != nil {
		t.Fatalf("first acquire failed: %v", err)
	}
	if l.Active("host1") != 1 {
		t.Errorf("expected active=1, got %d", l.Active("host1"))
	}
	l.Release("host1")
	if l.Active("host1") != 0 {
		t.Errorf("expected active=0 after release, got %d", l.Active("host1"))
	}
}

func TestAcquire_ExceedsMaxPerHost(t *testing.T) {
	l, _ := New(Policy{MaxPerHost: 1, Interval: 0})

	if err := l.Acquire("host1"); err != nil {
		t.Fatalf("first acquire failed: %v", err)
	}
	if err := l.Acquire("host1"); err == nil {
		t.Fatal("expected error when exceeding MaxPerHost")
	}
}

func TestAcquire_IntervalNotElapsed(t *testing.T) {
	l, _ := New(Policy{MaxPerHost: 5, Interval: 10 * time.Second})

	if err := l.Acquire("host2"); err != nil {
		t.Fatalf("first acquire failed: %v", err)
	}
	l.Release("host2")
	// Second acquire immediately — interval not elapsed.
	if err := l.Acquire("host2"); err == nil {
		t.Fatal("expected rate limit error before interval elapsed")
	}
}

func TestAcquire_DifferentHosts_Independent(t *testing.T) {
	l, _ := New(Policy{MaxPerHost: 1, Interval: 0})

	if err := l.Acquire("hostA"); err != nil {
		t.Fatalf("acquire hostA failed: %v", err)
	}
	// hostB should be independent of hostA.
	if err := l.Acquire("hostB"); err != nil {
		t.Fatalf("acquire hostB failed: %v", err)
	}
}

func TestRelease_BelowZero_Noop(t *testing.T) {
	l, _ := New(DefaultPolicy())
	// Release without prior acquire should not go negative.
	l.Release("ghost")
	if l.Active("ghost") != 0 {
		t.Errorf("expected active=0, got %d", l.Active("ghost"))
	}
}
