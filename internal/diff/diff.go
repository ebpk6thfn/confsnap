package diff

import (
	"fmt"
	"strings"

	"github.com/user/confsnap/internal/snapshot"
)

// Result holds the comparison result between two snapshots.
type Result struct {
	Host    string
	Path    string
	Changed bool
	OldHash string
	NewHash string
	Lines   []Line
}

// Line represents a single line difference.
type Line struct {
	Op      Op
	Content string
}

// Op represents the type of diff operation.
type Op int

const (
	OpEqual  Op = iota
	OpAdd
	OpRemove
)

func (o Op) String() string {
	switch o {
	case OpAdd:
		return "+"
	case OpRemove:
		return "-"
	default:
		return " "
	}
}

// Compare computes a line-level diff between two snapshots.
func Compare(old, new *snapshot.Snapshot) Result {
	result := Result{
		Host:    new.Host,
		Path:    new.Path,
		Changed: !old.Equal(new),
		OldHash: old.Checksum,
		NewHash: new.Checksum,
	}

	if !result.Changed {
		return result
	}

	oldLines := splitLines(old.Content)
	newLines := splitLines(new.Content)
	result.Lines = computeLineDiff(oldLines, newLines)
	return result
}

// Format returns a human-readable unified diff string.
func (r Result) Format() string {
	if !r.Changed {
		return fmt.Sprintf("[no change] %s:%s", r.Host, r.Path)
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "--- %s:%s (%s)\n", r.Host, r.Path, r.OldHash[:8])
	fmt.Fprintf(&sb, "+++ %s:%s (%s)\n", r.Host, r.Path, r.NewHash[:8])
	for _, l := range r.Lines {
		fmt.Fprintf(&sb, "%s %s\n", l.Op, l.Content)
	}
	return sb.String()
}

func splitLines(content string) []string {
	if content == "" {
		return nil
	}
	return strings.Split(content, "\n")
}

func computeLineDiff(old, new []string) []Line {
	var lines []Line
	oldSet := make(map[string]bool, len(old))
	newSet := make(map[string]bool, len(new))
	for _, l := range old {
		oldSet[l] = true
	}
	for _, l := range new {
		newSet[l] = true
	}
	for _, l := range old {
		if !newSet[l] {
			lines = append(lines, Line{Op: OpRemove, Content: l})
		} else {
			lines = append(lines, Line{Op: OpEqual, Content: l})
		}
	}
	for _, l := range new {
		if !oldSet[l] {
			lines = append(lines, Line{Op: OpAdd, Content: l})
		}
	}
	return lines
}
