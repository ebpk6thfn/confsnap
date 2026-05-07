package export

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/yourorg/confsnap/internal/diff"
)

// Format represents the output format for exported data.
type Format string

const (
	FormatCSV  Format = "csv"
	FormatJSON Format = "json"
)

// Record is a flat representation of a diff result suitable for export.
type Record struct {
	Timestamp string `json:"timestamp"`
	Host      string `json:"host"`
	Path      string `json:"path"`
	Status    string `json:"status"`
	Added     int    `json:"lines_added"`
	Removed   int    `json:"lines_removed"`
}

// Exporter writes diff results to an io.Writer in a chosen format.
type Exporter struct {
	format Format
	now    func() time.Time
}

// New creates an Exporter for the given format.
func New(format Format) (*Exporter, error) {
	switch format {
	case FormatCSV, FormatJSON:
		return &Exporter{format: format, now: time.Now}, nil
	default:
		return nil, fmt.Errorf("export: unsupported format %q", format)
	}
}

// Write serialises results to w in the configured format.
func (e *Exporter) Write(w io.Writer, results []diff.Result) error {
	records := e.toRecords(results)
	switch e.format {
	case FormatCSV:
		return writeCSV(w, records)
	case FormatJSON:
		return writeJSON(w, records)
	}
	return nil
}

func (e *Exporter) toRecords(results []diff.Result) []Record {
	ts := e.now().UTC().Format(time.RFC3339)
	records := make([]Record, 0, len(results))
	for _, r := range results {
		status := "unchanged"
		if r.Changed {
			status = "changed"
		}
		records = append(records, Record{
			Timestamp: ts,
			Host:      r.Host,
			Path:      r.Path,
			Status:    status,
			Added:     r.LinesAdded,
			Removed:   r.LinesRemoved,
		})
	}
	return records
}

func writeCSV(w io.Writer, records []Record) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{"timestamp", "host", "path", "status", "lines_added", "lines_removed"}); err != nil {
		return err
	}
	for _, r := range records {
		row := []string{
			r.Timestamp, r.Host, r.Path, r.Status,
			fmt.Sprintf("%d", r.Added),
			fmt.Sprintf("%d", r.Removed),
		}
		if err := cw.Write(row); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}

func writeJSON(w io.Writer, records []Record) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(records)
}
