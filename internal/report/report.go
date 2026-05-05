package report

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/user/confsnap/internal/diff"
)

// Format represents the output format for a report.
type Format string

const (
	FormatText Format = "text"
	FormatJSON  Format = "json"
)

// Report holds the results of a snapshot comparison run.
type Report struct {
	GeneratedAt time.Time
	Results     []diff.Result
}

// New creates a new Report from a slice of diff results.
func New(results []diff.Result) *Report {
	return &Report{
		GeneratedAt: time.Now(),
		Results:     results,
	}
}

// Write writes the report in the given format to w.
func (r *Report) Write(w io.Writer, format Format) error {
	switch format {
	case FormatJSON:
		return r.writeJSON(w)
	default:
		return r.writeText(w)
	}
}

// Summary returns counts of changed and unchanged files.
func (r *Report) Summary() (changed, unchanged int) {
	for _, res := range r.Results {
		if res.Changed {
			changed++
		} else {
			unchanged++
		}
	}
	return
}

func (r *Report) writeText(w io.Writer) error {
	changed, unchanged := r.Summary()
	fmt.Fprintf(w, "confsnap report — %s\n", r.GeneratedAt.Format(time.RFC3339))
	fmt.Fprintf(w, "Changed: %d  Unchanged: %d\n", changed, unchanged)
	fmt.Fprintln(w, strings.Repeat("-", 60))
	for _, res := range r.Results {
		fmt.Fprintln(w, res.Format())
	}
	return nil
}

func (r *Report) writeJSON(w io.Writer) error {
	changed, unchanged := r.Summary()
	fmt.Fprintf(w, `{"generated_at":%q,"summary":{"changed":%d,"unchanged":%d},"results":[`,
		r.GeneratedAt.Format(time.RFC3339), changed, unchanged)
	for i, res := range r.Results {
		if i > 0 {
			fmt.Fprint(w, ",")
		}
		fmt.Fprintf(w, `{"host":%q,"path":%q,"changed":%v}`,
			res.Host, res.Path, res.Changed)
	}
	fmt.Fprintln(w, "]}")
	return nil
}
