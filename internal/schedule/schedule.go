package schedule

import (
	"fmt"
	"time"
)

// Interval represents how often snapshots should be taken.
type Interval string

const (
	IntervalHourly  Interval = "hourly"
	IntervalDaily   Interval = "daily"
	IntervalWeekly  Interval = "weekly"
	IntervalManual  Interval = "manual"
)

// Schedule holds the configuration for automated snapshot scheduling.
type Schedule struct {
	Interval  Interval
	NextRunAt time.Time
	LastRunAt time.Time
}

// New creates a Schedule with the given interval, computing the first NextRunAt
// relative to the provided base time.
func New(interval Interval, base time.Time) (*Schedule, error) {
	if err := validateInterval(interval); err != nil {
		return nil, err
	}
	next, err := nextRun(interval, base)
	if err != nil {
		return nil, err
	}
	return &Schedule{
		Interval:  interval,
		NextRunAt: next,
	}, nil
}

// Due reports whether the schedule's next run time is at or before now.
func (s *Schedule) Due(now time.Time) bool {
	return !s.NextRunAt.After(now)
}

// Advance updates LastRunAt to now and recalculates NextRunAt.
func (s *Schedule) Advance(now time.Time) error {
	s.LastRunAt = now
	next, err := nextRun(s.Interval, now)
	if err != nil {
		return err
	}
	s.NextRunAt = next
	return nil
}

func validateInterval(i Interval) error {
	switch i {
	case IntervalHourly, IntervalDaily, IntervalWeekly, IntervalManual:
		return nil
	}
	return fmt.Errorf("schedule: unknown interval %q", i)
}

func nextRun(i Interval, from time.Time) (time.Time, error) {
	switch i {
	case IntervalHourly:
		return from.Add(time.Hour), nil
	case IntervalDaily:
		return from.Add(24 * time.Hour), nil
	case IntervalWeekly:
		return from.Add(7 * 24 * time.Hour), nil
	case IntervalManual:
		return time.Time{}, nil
	}
	return time.Time{}, fmt.Errorf("schedule: unknown interval %q", i)
}
