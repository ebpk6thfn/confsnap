package filter

import (
	"path/filepath"
	"strings"
)

// Rule defines a single include or exclude pattern for file paths.
type Rule struct {
	Pattern string
	Exclude bool
}

// Filter holds a set of rules used to determine whether a file path
// should be included in a snapshot run.
type Filter struct {
	rules []Rule
}

// New creates a Filter from lists of include and exclude glob patterns.
// Exclude patterns take precedence over include patterns.
func New(include, exclude []string) *Filter {
	f := &Filter{}
	for _, p := range include {
		f.rules = append(f.rules, Rule{Pattern: p, Exclude: false})
	}
	for _, p := range exclude {
		f.rules = append(f.rules, Rule{Pattern: p, Exclude: true})
	}
	return f
}

// Allow returns true if the given path should be processed.
// If no include rules are defined, all paths are allowed unless excluded.
// Exclude rules always win over include rules.
func (f *Filter) Allow(path string) bool {
	excluded := false
	included := false
	hasInclude := false

	for _, rule := range f.rules {
		matched := matchPattern(rule.Pattern, path)
		if rule.Exclude {
			if matched {
				excluded = true
			}
		} else {
			hasInclude = true
			if matched {
				included = true
			}
		}
	}

	if excluded {
		return false
	}
	if hasInclude {
		return included
	}
	return true
}

// matchPattern checks whether path matches a glob pattern, also
// supporting simple prefix matching for directory patterns ending in /.
func matchPattern(pattern, path string) bool {
	if strings.HasSuffix(pattern, "/") {
		return strings.HasPrefix(path, pattern) ||
			strings.HasPrefix(path, strings.TrimSuffix(pattern, "/"))
	}
	matched, err := filepath.Match(pattern, path)
	if err != nil {
		return false
	}
	// Also match basename for patterns without path separators.
	if !matched && !strings.Contains(pattern, "/") {
		matched, _ = filepath.Match(pattern, filepath.Base(path))
	}
	return matched
}
