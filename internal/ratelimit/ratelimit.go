package ratelimit

import (
	"fmt"
	"sync"
	"time"
)

// Policy defines rate limiting behaviour for SSH connections.
type Policy struct {
	// MaxPerHost is the maximum number of concurrent connections per host.
	MaxPerHost int
	// Interval is the minimum time between successive connections to the same host.
	Interval time.Duration
}

// DefaultPolicy returns a sensible default rate limit policy.
func DefaultPolicy() Policy {
	return Policy{
		MaxPerHost: 3,
		Interval:   500 * time.Millisecond,
	}
}

// Limiter enforces a rate limit policy across hosts.
type Limiter struct {
	policy  Policy
	mu      sync.Mutex
	lastHit map[string]time.Time
	active  map[string]int
}

// New creates a Limiter with the given policy.
func New(p Policy) (*Limiter, error) {
	if p.MaxPerHost <= 0 {
		return nil, fmt.Errorf("ratelimit: MaxPerHost must be > 0, got %d", p.MaxPerHost)
	}
	if p.Interval < 0 {
		return nil, fmt.Errorf("ratelimit: Interval must be >= 0, got %v", p.Interval)
	}
	return &Limiter{
		policy:  p,
		lastHit: make(map[string]time.Time),
		active:  make(map[string]int),
	}, nil
}

// Acquire attempts to acquire a slot for the given host.
// It returns an error if the concurrent limit is reached or the interval has
// not elapsed since the last connection.
func (l *Limiter) Acquire(host string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.active[host] >= l.policy.MaxPerHost {
		return fmt.Errorf("ratelimit: max concurrent connections (%d) reached for host %s",
			l.policy.MaxPerHost, host)
	}
	if last, ok := l.lastHit[host]; ok {
		if elapsed := time.Since(last); elapsed < l.policy.Interval {
			return fmt.Errorf("ratelimit: too soon to connect to %s (wait %v)",
				host, l.policy.Interval-elapsed)
		}
	}
	l.active[host]++
	l.lastHit[host] = time.Now()
	return nil
}

// Release decrements the active connection count for the given host.
func (l *Limiter) Release(host string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.active[host] > 0 {
		l.active[host]--
	}
}

// Active returns the current number of active connections for a host.
func (l *Limiter) Active(host string) int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.active[host]
}
