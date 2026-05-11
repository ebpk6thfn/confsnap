// Package template provides snapshot report templating using Go's text/template engine.
package template

import (
	"bytes"
	"fmt"
	"io"
	"text/template"
	"time"

	"github.com/yourorg/confsnap/internal/diff"
)

// Data holds the values passed to a report template.
type Data struct {
	GeneratedAt time.Time
	TotalHosts  int
	TotalFiles  int
	Changed     int
	Unchanged   int
	Results     []diff.Result
}

// Renderer renders diff results using a named Go template.
type Renderer struct {
	tmpl *template.Template
}

var funcMap = template.FuncMap{
	"formatTime": func(t time.Time) string {
		return t.Format(time.RFC3339)
	},
	"inc": func(i int) int { return i + 1 },
}

// New creates a Renderer from the given template source string.
func New(src string) (*Renderer, error) {
	if src == "" {
		return nil, fmt.Errorf("template: source must not be empty")
	}
	t, err := template.New("report").Funcs(funcMap).Parse(src)
	if err != nil {
		return nil, fmt.Errorf("template: parse error: %w", err)
	}
	return &Renderer{tmpl: t}, nil
}

// Render executes the template with the provided Data and writes output to w.
func (r *Renderer) Render(w io.Writer, d Data) error {
	var buf bytes.Buffer
	if err := r.tmpl.Execute(&buf, d); err != nil {
		return fmt.Errorf("template: render error: %w", err)
	}
	_, err := w.Write(buf.Bytes())
	return err
}

// BuildData constructs a Data value from a slice of diff results.
func BuildData(results []diff.Result) Data {
	d := Data{
		GeneratedAt: time.Now().UTC(),
		Results:     results,
	}
	hosts := make(map[string]struct{})
	for _, r := range results {
		hosts[r.Host] = struct{}{}
		d.TotalFiles++
		if r.Changed {
			d.Changed++
		} else {
			d.Unchanged++
		}
	}
	d.TotalHosts = len(hosts)
	return d
}
