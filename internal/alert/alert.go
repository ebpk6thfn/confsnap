package alert

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/user/confsnap/internal/diff"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelInfo  Level = "info"
	LevelWarn  Level = "warn"
	LevelCrit  Level = "crit"
)

// Alert represents a single config drift alert.
type Alert struct {
	Host      string
	Path      string
	Level     Level
	Message   string
	ChangedAt time.Time
}

// Alerter evaluates diff results and emits alerts.
type Alerter struct {
	threshold int // minimum changed lines to trigger warn/crit
	writer    io.Writer
}

// New creates an Alerter that writes alerts to w.
// threshold is the number of changed lines that triggers a warning;
// double the threshold triggers a critical alert.
func New(w io.Writer, threshold int) *Alerter {
	if threshold <= 0 {
		threshold = 1
	}
	return &Alerter{threshold: threshold, writer: w}
}

// Evaluate inspects diff results and returns alerts for any drifted configs.
func (a *Alerter) Evaluate(results []diff.Result) []Alert {
	var alerts []Alert
	for _, r := range results {
		if !r.Changed {
			continue
		}
		changed := countChangedLines(r.Diff)
		lvl := LevelInfo
		switch {
		case changed >= a.threshold*2:
			lvl = LevelCrit
		case changed >= a.threshold:
			lvl = LevelWarn
		}
		alerts = append(alerts, Alert{
			Host:      r.Host,
			Path:      r.Path,
			Level:     lvl,
			Message:   fmt.Sprintf("%d line(s) changed", changed),
			ChangedAt: time.Now().UTC(),
		})
	}
	return alerts
}

// Write formats and writes alerts to the configured writer.
func (a *Alerter) Write(alerts []Alert) error {
	for _, al := range alerts {
		line := fmt.Sprintf("[%s] %s %s — %s\n",
			strings.ToUpper(string(al.Level)),
			al.Host,
			al.Path,
			al.Message,
		)
		if _, err := fmt.Fprint(a.writer, line); err != nil {
			return err
		}
	}
	return nil
}

func countChangedLines(lines []string) int {
	count := 0
	for _, l := range lines {
		if strings.HasPrefix(l, "+") || strings.HasPrefix(l, "-") {
			count++
		}
	}
	return count
}
