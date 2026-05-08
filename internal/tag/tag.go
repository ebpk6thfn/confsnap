// Package tag provides host and file tagging for grouping and filtering snapshots.
package tag

import (
	"fmt"
	"sort"
	"strings"
)

// Registry holds tags associated with hosts and files.
type Registry struct {
	hostTags map[string][]string
	fileTags map[string][]string
}

// New returns an empty Registry.
func New() *Registry {
	return &Registry{
		hostTags: make(map[string][]string),
		fileTags: make(map[string][]string),
	}
}

// TagHost associates one or more tags with a host.
func (r *Registry) TagHost(host string, tags ...string) {
	r.hostTags[host] = mergeTags(r.hostTags[host], tags)
}

// TagFile associates one or more tags with a file path.
func (r *Registry) TagFile(path string, tags ...string) {
	r.fileTags[path] = mergeTags(r.fileTags[path], tags)
}

// HostTags returns sorted tags for the given host.
func (r *Registry) HostTags(host string) []string {
	return sorted(r.hostTags[host])
}

// FileTags returns sorted tags for the given file path.
func (r *Registry) FileTags(path string) []string {
	return sorted(r.fileTags[path])
}

// HostsWithTag returns all hosts that carry the given tag.
func (r *Registry) HostsWithTag(tag string) []string {
	var out []string
	for host, tags := range r.hostTags {
		if containsTag(tags, tag) {
			out = append(out, host)
		}
	}
	sort.Strings(out)
	return out
}

// FilesWithTag returns all file paths that carry the given tag.
func (r *Registry) FilesWithTag(tag string) []string {
	var out []string
	for path, tags := range r.fileTags {
		if containsTag(tags, tag) {
			out = append(out, path)
		}
	}
	sort.Strings(out)
	return out
}

// Validate checks that all tags are non-empty and contain no whitespace.
func Validate(tags []string) error {
	for _, t := range tags {
		if strings.TrimSpace(t) == "" {
			return fmt.Errorf("tag must not be empty or whitespace")
		}
		if strings.ContainsAny(t, " \t\n") {
			return fmt.Errorf("tag %q must not contain whitespace", t)
		}
	}
	return nil
}

func mergeTags(existing, incoming []string) []string {
	seen := make(map[string]struct{}, len(existing))
	for _, t := range existing {
		seen[t] = struct{}{}
	}
	result := append([]string{}, existing...)
	for _, t := range incoming {
		if _, ok := seen[t]; !ok {
			result = append(result, t)
			seen[t] = struct{}{}
		}
	}
	return result
}

func sorted(tags []string) []string {
	copy := append([]string{}, tags...)
	sort.Strings(copy)
	return copy
}

func containsTag(tags []string, target string) bool {
	for _, t := range tags {
		if t == target {
			return true
		}
	}
	return false
}
