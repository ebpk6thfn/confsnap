package encrypt

import (
	"fmt"
	"os"
	"path/filepath"
)

// Store wraps an Encryptor to transparently encrypt files written to disk
// and decrypt them on read.
type Store struct {
	enc *Encryptor
	dir string
}

// NewStore creates a Store that persists encrypted files under dir.
func NewStore(dir, passphrase string) (*Store, error) {
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, fmt.Errorf("encrypt.NewStore: %w", err)
	}
	return &Store{enc: New(passphrase), dir: dir}, nil
}

// Write encrypts data and writes it to a file named by key inside the store
// directory.
func (s *Store) Write(key string, data []byte) error {
	ct, err := s.enc.Encrypt(data)
	if err != nil {
		return fmt.Errorf("encrypt.Store.Write: %w", err)
	}
	path := filepath.Join(s.dir, sanitizeKey(key))
	if err := os.WriteFile(path, ct, 0o600); err != nil {
		return fmt.Errorf("encrypt.Store.Write: %w", err)
	}
	return nil
}

// Read decrypts and returns the data stored under key, or an error if the
// file does not exist or decryption fails.
func (s *Store) Read(key string) ([]byte, error) {
	path := filepath.Join(s.dir, sanitizeKey(key))
	ct, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("encrypt.Store.Read: %w", err)
	}
	plaintext, err := s.enc.Decrypt(ct)
	if err != nil {
		return nil, fmt.Errorf("encrypt.Store.Read: %w", err)
	}
	return plaintext, nil
}

// sanitizeKey replaces path separators so that keys can be used as file names.
func sanitizeKey(key string) string {
	safe := make([]byte, len(key))
	for i := 0; i < len(key); i++ {
		if key[i] == '/' || key[i] == '\\' {
			safe[i] = '_'
		} else {
			safe[i] = key[i]
		}
	}
	return string(safe)
}
