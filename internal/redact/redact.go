// Package redact provides utilities for masking sensitive values in
// configuration snapshots before storing or displaying them.
package redact

import (
	"regexp"
	"strings"
)

// Rule describes a single redaction rule: a pattern to match a config key
// and the mask string to replace its value with.
type Rule struct {
	Pattern string
	Mask    string
}

// Redactor applies a set of rules to config file lines.
type Redactor struct {
	rules    []Rule
	compiled []*regexp.Regexp
}

const defaultMask = "***REDACTED***"

// New creates a Redactor from the given rules. Returns an error if any
// pattern fails to compile.
func New(rules []Rule) (*Redactor, error) {
	compiled := make([]*regexp.Regexp, 0, len(rules))
	for _, r := range rules {
		re, err := regexp.Compile(r.Pattern)
		if err != nil {
			return nil, err
		}
		compiled = append(compiled, re)
	}
	return &Redactor{rules: rules, compiled: compiled}, nil
}

// Apply scans each line and replaces the value portion of any key=value (or
// key: value) pair whose key matches a rule pattern.
func (r *Redactor) Apply(lines []string) []string {
	out := make([]string, len(lines))
	for i, line := range lines {
		out[i] = r.redactLine(line)
	}
	return out
}

// ApplyString is a convenience wrapper around Apply for a full file body.
func (r *Redactor) ApplyString(content string) string {
	lines := strings.Split(content, "\n")
	return strings.Join(r.Apply(lines), "\n")
}

func (r *Redactor) redactLine(line string) string {
	// Support both "key=value" and "key: value" formats.
	var sep string
	var keyPart, valPart string

	if idx := strings.Index(line, "="); idx != -1 {
		sep = "="
		keyPart = strings.TrimSpace(line[:idx])
		valPart = line[idx+1:]
	} else if idx := strings.Index(line, ":"); idx != -1 {
		sep = ":"
		keyPart = strings.TrimSpace(line[:idx])
		valPart = line[idx+1:]
	} else {
		return line
	}

	for i, re := range r.compiled {
		if re.MatchString(keyPart) {
			mask := r.rules[i].Mask
			if mask == "" {
				mask = defaultMask
			}
			// Preserve leading whitespace in value portion.
			leading := valPart[:len(valPart)-len(strings.TrimLeft(valPart, " \t"))]
			return line[:strings.Index(line, sep)+1] + leading + mask
		}
	}
	return line
}
