// Package compress provides gzip compression and decompression utilities
// for snapshot data stored on disk.
package compress

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
)

// Level represents a compression level.
type Level int

const (
	LevelDefault Level = Level(gzip.DefaultCompression)
	LevelBest    Level = Level(gzip.BestCompression)
	LevelFast    Level = Level(gzip.BestSpeed)
	LevelNone    Level = Level(gzip.NoCompression)
)

// Compressor compresses and decompresses data using gzip.
type Compressor struct {
	level Level
}

// New returns a new Compressor with the given compression level.
func New(level Level) *Compressor {
	return &Compressor{level: level}
}

// Compress compresses src using gzip and returns the compressed bytes.
func (c *Compressor) Compress(src []byte) ([]byte, error) {
	var buf bytes.Buffer
	w, err := gzip.NewWriterLevel(&buf, int(c.level))
	if err != nil {
		return nil, fmt.Errorf("compress: create writer: %w", err)
	}
	if _, err := w.Write(src); err != nil {
		return nil, fmt.Errorf("compress: write: %w", err)
	}
	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("compress: close writer: %w", err)
	}
	return buf.Bytes(), nil
}

// Decompress decompresses gzip-compressed src and returns the original bytes.
func (c *Compressor) Decompress(src []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(src))
	if err != nil {
		return nil, fmt.Errorf("compress: create reader: %w", err)
	}
	defer r.Close()

	out, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("compress: read: %w", err)
	}
	return out, nil
}

// Ratio returns the compression ratio (compressed/original) for the given data.
// A value less than 1.0 means the data was compressed successfully.
func (c *Compressor) Ratio(src []byte) (float64, error) {
	if len(src) == 0 {
		return 1.0, nil
	}
	compressed, err := c.Compress(src)
	if err != nil {
		return 0, err
	}
	return float64(len(compressed)) / float64(len(src)), nil
}
