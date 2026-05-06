package retention

import (
	"fmt"
	"time"
)

// Config holds the YAML-deserializable retention settings.
type Config struct {
	// MaxAgeDays is the maximum age in days before a snapshot is pruned.
	// Zero disables age-based pruning.
	MaxAgeDays int `yaml:"max_age_days"`
	// MaxCount is the maximum number of snapshots to keep per host+path.
	// Zero disables count-based pruning.
	MaxCount int `yaml:"max_count"`
}

// ToPolicy converts a Config to a Policy, validating values.
func (c Config) ToPolicy() (Policy, error) {
	if c.MaxAgeDays < 0 {
		return Policy{}, fmt.Errorf("retention: max_age_days must be >= 0, got %d", c.MaxAgeDays)
	}
	if c.MaxCount < 0 {
		return Policy{}, fmt.Errorf("retention: max_count must be >= 0, got %d", c.MaxCount)
	}
	return Policy{
		MaxAge:   time.Duration(c.MaxAgeDays) * 24 * time.Hour,
		MaxCount: c.MaxCount,
	}, nil
}

// DefaultConfig returns a sensible default retention configuration.
func DefaultConfig() Config {
	return Config{
		MaxAgeDays: 30,
		MaxCount:   10,
	}
}
