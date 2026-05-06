package retention

import (
	"testing"
	"time"
)

func TestToPolicy_Valid(t *testing.T) {
	c := Config{MaxAgeDays: 7, MaxCount: 5}
	p, err := c.ToPolicy()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.MaxAge != 7*24*time.Hour {
		t.Errorf("MaxAge = %v, want %v", p.MaxAge, 7*24*time.Hour)
	}
	if p.MaxCount != 5 {
		t.Errorf("MaxCount = %d, want 5", p.MaxCount)
	}
}

func TestToPolicy_ZeroDisablesLimits(t *testing.T) {
	c := Config{MaxAgeDays: 0, MaxCount: 0}
	p, err := c.ToPolicy()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.MaxAge != 0 {
		t.Errorf("expected MaxAge 0, got %v", p.MaxAge)
	}
	if p.MaxCount != 0 {
		t.Errorf("expected MaxCount 0, got %d", p.MaxCount)
	}
}

func TestToPolicy_NegativeMaxAgeDays(t *testing.T) {
	c := Config{MaxAgeDays: -1}
	_, err := c.ToPolicy()
	if err == nil {
		t.Fatal("expected error for negative MaxAgeDays")
	}
}

func TestToPolicy_NegativeMaxCount(t *testing.T) {
	c := Config{MaxCount: -3}
	_, err := c.ToPolicy()
	if err == nil {
		t.Fatal("expected error for negative MaxCount")
	}
}

func TestDefaultConfig(t *testing.T) {
	c := DefaultConfig()
	if c.MaxAgeDays <= 0 {
		t.Errorf("DefaultConfig MaxAgeDays should be positive, got %d", c.MaxAgeDays)
	}
	if c.MaxCount <= 0 {
		t.Errorf("DefaultConfig MaxCount should be positive, got %d", c.MaxCount)
	}
}
