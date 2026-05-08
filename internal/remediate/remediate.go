package remediate

import (
	"fmt"
	"time"

	"github.com/yourorg/confsnap/internal/diff"
	"github.com/yourorg/confsnap/internal/rollback"
)

// Action describes what should happen when drift is detected.
type Action string

const (
	ActionRollback Action = "rollback"
	ActionNotify   Action = "notify"
)

// Rule maps a file path pattern to a remediation action.
type Rule struct {
	Pattern string
	Action  Action
}

// Result records the outcome of a remediation attempt.
type Result struct {
	Host      string
	Path      string
	Action    Action
	Applied   bool
	Message   string
	Timestamp time.Time
}

// Remediator evaluates diff results and applies configured rules.
type Remediator struct {
	rules   []Rule
	store   *rollback.Store
}

// New creates a Remediator with the given rules and rollback store.
func New(rules []Rule, store *rollback.Store) *Remediator {
	return &Remediator{rules: rules, store: store}
}

// Evaluate checks each diff result against rules and returns remediation results.
func (r *Remediator) Evaluate(results []diff.Result) []Result {
	var out []Result
	for _, res := range results {
		if !res.Changed {
			continue
		}
		action := r.matchRule(res.Path)
		if action == "" {
			continue
		}
		out = append(out, r.apply(res, action))
	}
	return out
}

func (r *Remediator) matchRule(path string) Action {
	for _, rule := range r.rules {
		if matchGlob(rule.Pattern, path) {
			return rule.Action
		}
	}
	return ""
}

func (r *Remediator) apply(res diff.Result, action Action) Result {
	out := Result{
		Host:      res.Host,
		Path:      res.Path,
		Action:    action,
		Timestamp: time.Now().UTC(),
	}
	switch action {
	case ActionRollback:
		entry, err := r.store.Load(res.Host, res.Path)
		if err != nil || entry == nil {
			out.Applied = false
			out.Message = fmt.Sprintf("no rollback entry found: %v", err)
		} else {
			out.Applied = true
			out.Message = fmt.Sprintf("rolled back to snapshot %s", entry.SnapshotID)
		}
	case ActionNotify:
		out.Applied = true
		out.Message = fmt.Sprintf("drift detected on %s:%s", res.Host, res.Path)
	}
	return out
}
