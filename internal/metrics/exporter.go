package metrics

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
)

// Format controls the output format of an Exporter.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Exporter writes a Registry Snapshot to an io.Writer.
type Exporter struct {
	reg    *Registry
	format Format
}

// NewExporter returns an Exporter for the given registry and format.
// Returns an error for unknown formats.
func NewExporter(reg *Registry, format Format) (*Exporter, error) {
	switch format {
	case FormatText, FormatJSON:
	default:
		return nil, fmt.Errorf("metrics: unknown format %q", format)
	}
	return &Exporter{reg: reg, format: format}, nil
}

// Write serialises the current snapshot to w.
func (e *Exporter) Write(w io.Writer) error {
	s := e.reg.Snapshot()
	switch e.format {
	case FormatJSON:
		return e.writeJSON(w, s)
	default:
		return e.writeText(w, s)
	}
}

func (e *Exporter) writeText(w io.Writer, s Snapshot) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# snapshot at %s\n", s.At.Format("2006-01-02T15:04:05Z")))

	keys := make([]string, 0, len(s.Counters))
	for k := range s.Counters {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		sb.WriteString(fmt.Sprintf("counter %s %d\n", k, s.Counters[k]))
	}

	keys = keys[:0]
	for k := range s.Gauges {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		sb.WriteString(fmt.Sprintf("gauge   %s %g\n", k, s.Gauges[k]))
	}

	_, err := io.WriteString(w, sb.String())
	return err
}

func (e *Exporter) writeJSON(w io.Writer, s Snapshot) error {
	payload := map[string]any{
		"at":       s.At,
		"counters": s.Counters,
		"gauges":   s.Gauges,
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(payload)
}
