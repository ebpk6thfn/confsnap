package healthcheck

import (
	"fmt"
	"time"
)

// Status represents the health status of a host.
type Status int

const (
	StatusUnknown Status = iota
	StatusHealthy
	StatusUnreachable
	StatusDegraded
)

func (s Status) String() string {
	switch s {
	case StatusHealthy:
		return "healthy"
	case StatusUnreachable:
		return "unreachable"
	case StatusDegraded:
		return "degraded"
	default:
		return "unknown"
	}
}

// Result holds the outcome of a health check for a single host.
type Result struct {
	Host      string
	Status    Status
	Latency   time.Duration
	CheckedAt time.Time
	Error     string
}

// Checker performs SSH reachability checks against a list of hosts.
type Checker struct {
	hosts   []string
	timeout time.Duration
	probe   func(host string, timeout time.Duration) (time.Duration, error)
}

// New creates a Checker for the given hosts.
// timeout controls how long each probe may take.
func New(hosts []string, timeout time.Duration) *Checker {
	return &Checker{
		hosts:   hosts,
		timeout: timeout,
		probe:   defaultProbe,
	}
}

// Check runs health probes against all configured hosts and returns results.
func (c *Checker) Check() []Result {
	results := make([]Result, 0, len(c.hosts))
	for _, h := range c.hosts {
		r := Result{Host: h, CheckedAt: time.Now()}
		latency, err := c.probe(h, c.timeout)
		if err != nil {
			r.Status = StatusUnreachable
			r.Error = err.Error()
		} else {
			r.Status = StatusHealthy
			r.Latency = latency
		}
		results = append(results, r)
	}
	return results
}

// Healthy returns true only when all results indicate a healthy status.
func Healthy(results []Result) bool {
	for _, r := range results {
		if r.Status != StatusHealthy {
			return false
		}
	}
	return true
}

// defaultProbe dials the SSH port of host and measures round-trip time.
func defaultProbe(host string, timeout time.Duration) (time.Duration, error) {
	import_net := func() (time.Duration, error) {
		start := time.Now()
		conn, err := dialTCP(host+":22", timeout)
		if err != nil {
			return 0, fmt.Errorf("dial %s: %w", host, err)
		}
		conn.Close()
		return time.Since(start), nil
	}
	return import_net()
}
