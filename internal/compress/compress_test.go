package compress_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/confsnap/internal/compress"
)

func TestNew_NotNil(t *testing.T) {
	c := compress.New(compress.LevelDefault)
	if c == nil {
		t.Fatal("expected non-nil Compressor")
	}
}

func TestCompress_ThenDecompress_RoundTrip(t *testing.T) {
	c := compress.New(compress.LevelDefault)
	orig := []byte(strings.Repeat("config line\n", 100))

	compressed, err := c.Compress(orig)
	if err != nil {
		t.Fatalf("Compress: %v", err)
	}

	got, err := c.Decompress(compressed)
	if err != nil {
		t.Fatalf("Decompress: %v", err)
	}

	if !bytes.Equal(orig, got) {
		t.Errorf("round-trip mismatch: got %d bytes, want %d bytes", len(got), len(orig))
	}
}

func TestCompress_ReducesSize(t *testing.T) {
	c := compress.New(compress.LevelBest)
	src := []byte(strings.Repeat("aaaa", 512))

	out, err := c.Compress(src)
	if err != nil {
		t.Fatalf("Compress: %v", err)
	}
	if len(out) >= len(src) {
		t.Errorf("expected compressed size < original: %d >= %d", len(out), len(src))
	}
}

func TestDecompress_InvalidData_ReturnsError(t *testing.T) {
	c := compress.New(compress.LevelDefault)
	_, err := c.Decompress([]byte("not gzip data"))
	if err == nil {
		t.Fatal("expected error for invalid gzip data")
	}
}

func TestCompress_EmptyInput(t *testing.T) {
	c := compress.New(compress.LevelDefault)
	out, err := c.Compress([]byte{})
	if err != nil {
		t.Fatalf("Compress empty: %v", err)
	}
	got, err := c.Decompress(out)
	if err != nil {
		t.Fatalf("Decompress empty: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty output, got %d bytes", len(got))
	}
}

func TestRatio_RepetitiveData_LessThanOne(t *testing.T) {
	c := compress.New(compress.LevelDefault)
	src := []byte(strings.Repeat("x", 1000))
	ratio, err := c.Ratio(src)
	if err != nil {
		t.Fatalf("Ratio: %v", err)
	}
	if ratio >= 1.0 {
		t.Errorf("expected ratio < 1.0 for repetitive data, got %.3f", ratio)
	}
}

func TestRatio_EmptyInput_ReturnsOne(t *testing.T) {
	c := compress.New(compress.LevelDefault)
	ratio, err := c.Ratio([]byte{})
	if err != nil {
		t.Fatalf("Ratio: %v", err)
	}
	if ratio != 1.0 {
		t.Errorf("expected ratio 1.0 for empty input, got %.3f", ratio)
	}
}
