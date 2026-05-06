package schedule

import (
	"testing"
	"time"
)

var base = time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

func TestNew_ValidInterval(t *testing.T) {
	for _, interval := range []Interval{IntervalHourly, IntervalDaily, IntervalWeekly, IntervalManual} {
		s, err := New(interval, base)
		if err != nil {
			t.Fatalf("New(%q): unexpected error: %v", interval, err)
		}
		if s.Interval != interval {
			t.Errorf("expected interval %q, got %q", interval, s.Interval)
		}
	}
}

func TestNew_InvalidInterval(t *testing.T) {
	_, err := New("monthly", base)
	if err == nil {
		t.Fatal("expected error for unknown interval, got nil")
	}
}

func TestNew_NextRunAt_Hourly(t *testing.T) {
	s, _ := New(IntervalHourly, base)
	expected := base.Add(time.Hour)
	if !s.NextRunAt.Equal(expected) {
		t.Errorf("expected NextRunAt %v, got %v", expected, s.NextRunAt)
	}
}

func TestNew_NextRunAt_Manual_IsZero(t *testing.T) {
	s, _ := New(IntervalManual, base)
	if !s.NextRunAt.IsZero() {
		t.Errorf("expected zero NextRunAt for manual, got %v", s.NextRunAt)
	}
}

func TestDue_BeforeNextRun(t *testing.T) {
	s, _ := New(IntervalDaily, base)
	if s.Due(base) {
		t.Error("expected not due immediately after creation")
	}
}

func TestDue_AfterNextRun(t *testing.T) {
	s, _ := New(IntervalDaily, base)
	future := base.Add(25 * time.Hour)
	if !s.Due(future) {
		t.Error("expected due after interval elapsed")
	}
}

func TestAdvance_UpdatesLastAndNext(t *testing.T) {
	s, _ := New(IntervalHourly, base)
	now := base.Add(time.Hour)
	if err := s.Advance(now); err != nil {
		t.Fatalf("Advance: unexpected error: %v", err)
	}
	if !s.LastRunAt.Equal(now) {
		t.Errorf("expected LastRunAt %v, got %v", now, s.LastRunAt)
	}
	expectedNext := now.Add(time.Hour)
	if !s.NextRunAt.Equal(expectedNext) {
		t.Errorf("expected NextRunAt %v, got %v", expectedNext, s.NextRunAt)
	}
}
