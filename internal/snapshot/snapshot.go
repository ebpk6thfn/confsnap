package snapshot

import (
	"crypto/sha256"
	"fmt"
	"time"
)

// Snapshot represents a captured state of a configuration file on a remote host.
type Snapshot struct {
	Host      string    `json:"host"`
	FilePath  string    `json:"file_path"`
	Content   string    `json:"content"`
	Checksum  string    `json:"checksum"`
	CapturedAt time.Time `json:"captured_at"`
}

// New creates a new Snapshot with a computed SHA-256 checksum.
func New(host, filePath, content string) *Snapshot {
	return &Snapshot{
		Host:       host,
		FilePath:   filePath,
		Content:    content,
		Checksum:   computeChecksum(content),
		CapturedAt: time.Now().UTC(),
	}
}

// computeChecksum returns the SHA-256 hex digest of the given content.
func computeChecksum(content string) string {
	sum := sha256.Sum256([]byte(content))
	return fmt.Sprintf("%x", sum)
}

// Equal reports whether two snapshots have identical content.
func (s *Snapshot) Equal(other *Snapshot) bool {
	return s.Checksum == other.Checksum
}
