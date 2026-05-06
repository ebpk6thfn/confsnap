package baseline

import (
	"fmt"

	"github.com/yourorg/confsnap/internal/diff"
	"github.com/yourorg/confsnap/internal/snapshot"
)

// DriftResult holds the diff between a baseline snapshot and a current one.
type DriftResult struct {
	Host    string
	Path    string
	Diff    *diff.Result
	Missing bool // true if the file was in the baseline but not collected now
}

// CompareToCurrent diffs each snapshot in the baseline against currentSnaps.
// currentSnaps is keyed by "host:path" matching the baseline format.
func CompareToCurrent(b *Baseline, currentSnaps map[string]*snapshot.Snapshot) []DriftResult {
	var results []DriftResult

	for key, base := range b.Snapshots {
		current, ok := currentSnaps[key]
		if !ok {
			results = append(results, DriftResult{
				Host:    base.Host,
				Path:    base.Path,
				Missing: true,
			})
			continue
		}
		d := diff.Compare(base, current)
		results = append(results, DriftResult{
			Host: base.Host,
			Path: base.Path,
			Diff: d,
		})
	}
	return results
}

// SnapshotKey returns the canonical map key for a snapshot.
func SnapshotKey(host, path string) string {
	return fmt.Sprintf("%s:%s", host, path)
}
