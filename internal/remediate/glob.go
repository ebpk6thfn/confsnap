package remediate

import "path"

// matchGlob reports whether name matches the shell pattern.
// It delegates to path.Match and treats any error as no match.
func matchGlob(pattern, name string) bool {
	matched, err := path.Match(pattern, name)
	if err != nil {
		return false
	}
	return matched
}
