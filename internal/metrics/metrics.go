package metrics

import (
	"sync"
	"time"
)

// Counter tracks a named integer count.
type Counter struct {
	mu    sync.Mutex
	value int64
}

func (c *Counter) Inc() { c.mu.Lock(); c.value++; c.mu.Unlock() }
func (c *Counter) Add(n int64) { c.mu.Lock(); c.value += n; c.mu.Unlock() }
func (c *Counter) Value() int64 { c.mu.Lock(); defer c.mu.Unlock(); return c.value }

// Gauge tracks a named float value.
type Gauge struct {
	mu    sync.Mutex
	value float64
}

func (g *Gauge) Set(v float64) { g.mu.Lock(); g.value = v; g.mu.Unlock() }
func (g *Gauge) Value() float64 { g.mu.Lock(); defer g.mu.Unlock(); return g.value }

// Snapshot is a point-in-time view of all metrics.
type Snapshot struct {
	Counters map[string]int64
	Gauges   map[string]float64
	At       time.Time
}

// Registry holds named counters and gauges.
type Registry struct {
	mu       sync.Mutex
	counters map[string]*Counter
	gauges   map[string]*Gauge
}

// New returns an empty Registry.
func New() *Registry {
	return &Registry{
		counters: make(map[string]*Counter),
		gauges:   make(map[string]*Gauge),
	}
}

// Counter returns (creating if needed) the named counter.
func (r *Registry) Counter(name string) *Counter {
	r.mu.Lock()
	defer r.mu.Unlock()
	if c, ok := r.counters[name]; ok {
		return c
	}
	c := &Counter{}
	r.counters[name] = c
	return c
}

// Gauge returns (creating if needed) the named gauge.
func (r *Registry) Gauge(name string) *Gauge {
	r.mu.Lock()
	defer r.mu.Unlock()
	if g, ok := r.gauges[name]; ok {
		return g
	}
	g := &Gauge{}
	r.gauges[name] = g
	return g
}

// Snapshot returns a copy of all current metric values.
func (r *Registry) Snapshot() Snapshot {
	r.mu.Lock()
	defer r.mu.Unlock()
	s := Snapshot{
		Counters: make(map[string]int64, len(r.counters)),
		Gauges:   make(map[string]float64, len(r.gauges)),
		At:       time.Now().UTC(),
	}
	for k, c := range r.counters {
		s.Counters[k] = c.Value()
	}
	for k, g := range r.gauges {
		s.Gauges[k] = g.Value()
	}
	return s
}
