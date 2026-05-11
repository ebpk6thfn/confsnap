package encrypt_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/confsnap/internal/encrypt"
)

func tempDir(t *testing.T) string {
	t.Helper()
	d, err := os.MkdirTemp("", "encrypt-store-*")
	if err != nil {
		t.Fatalf("MkdirTemp: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(d) })
	return d
}

func TestNewStore_CreatesDirectory(t *testing.T) {
	dir := filepath.Join(tempDir(t), "sub", "dir")
	_, err := encrypt.NewStore(dir, "pass")
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	if _, err := os.Stat(dir); err != nil {
		t.Errorf("expected directory to exist: %v", err)
	}
}

func TestStore_Write_ThenRead_RoundTrip(t *testing.T) {
	s, err := encrypt.NewStore(tempDir(t), "passphrase")
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	data := []byte("[server]\nhost = db01")
	if err := s.Write("web01:/etc/app.conf", data); err != nil {
		t.Fatalf("Write: %v", err)
	}
	got, err := s.Read("web01:/etc/app.conf")
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if string(got) != string(data) {
		t.Errorf("got %q, want %q", got, data)
	}
}

func TestStore_Read_MissingKey_ReturnsError(t *testing.T) {
	s, _ := encrypt.NewStore(tempDir(t), "passphrase")
	_, err := s.Read("nonexistent")
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestStore_Write_FileIsNotPlaintext(t *testing.T) {
	dir := tempDir(t)
	s, _ := encrypt.NewStore(dir, "passphrase")
	data := []byte("secret=hunter2")
	_ = s.Write("cfg", data)

	raw, _ := os.ReadFile(filepath.Join(dir, "cfg"))
	if string(raw) == string(data) {
		t.Error("file should not contain plaintext")
	}
}

func TestStore_Read_WrongPassphrase_ReturnsError(t *testing.T) {
	dir := tempDir(t)
	s1, _ := encrypt.NewStore(dir, "correct")
	_ = s1.Write("key", []byte("value"))

	s2, _ := encrypt.NewStore(dir, "wrong")
	_, err := s2.Read("key")
	if err == nil {
		t.Fatal("expected decryption error with wrong passphrase")
	}
}
