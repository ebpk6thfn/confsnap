package ratelimit

import (
	"fmt"
	"time"
)

// Config is the user-facing configuration for rate limiting, suitable for
// unmarshalling from YAML/JSON.
type Config struct {
	// MaxPerHost is the maximum concurrent SSH connections per host.
	MaxPerHost int `yaml:"max_per_host" json:"max_per_host"`
	// IntervalMs is the minimum milliseconds between connections to a host.
	IntervalMs int `yaml:"interval_ms" json:"interval_ms"`
}

// DefaultConfig returns a Config that mirrors DefaultPolicy.
func DefaultConfig() Config {
	p := DefaultPolicy()
	return Config{
		MaxPerHost: p.MaxPerHost,
		IntervalMs: int(p.Interval.Milliseconds()),
	}
}

// ToPolicy converts Config to a Policy, validating values.
func (c Config) ToPolicy() (Policy, error) {
	if c.MaxPerHost <= 0 {
		return Policy{}, fmt.Errorf("ratelimit: max_per_host must be > 0, got %d", c.MaxPerHost)
	}
	if c.IntervalMs < 0 {
		return Policy{}, fmt.Errorf("ratelimit: interval_ms must be >= 0, got %d", c.IntervalMs)
	}
	return Policy{
		MaxPerHost: c.MaxPerHost,
		Interval:   time.Duration(c.IntervalMs) * time.Millisecond,
	}, nil
}
