// Package encrypt provides AES-GCM encryption and decryption for snapshot
// data stored on disk, ensuring sensitive configuration content is protected
// at rest.
package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"
)

// ErrInvalidCiphertext is returned when decryption fails due to malformed input.
var ErrInvalidCiphertext = errors.New("encrypt: invalid or corrupted ciphertext")

// Encryptor encrypts and decrypts data using AES-256-GCM.
type Encryptor struct {
	key [32]byte
}

// New creates a new Encryptor derived from the given passphrase.
// The passphrase is hashed with SHA-256 to produce a 32-byte key.
func New(passphrase string) *Encryptor {
	return &Encryptor{key: sha256.Sum256([]byte(passphrase))}
}

// Encrypt encrypts plaintext using AES-256-GCM and returns the
// nonce-prefixed ciphertext. Each call uses a freshly generated nonce.
func (e *Encryptor) Encrypt(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(e.key[:])
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	sealed := gcm.Seal(nonce, nonce, plaintext, nil)
	return sealed, nil
}

// Decrypt decrypts a nonce-prefixed ciphertext produced by Encrypt.
func (e *Encryptor) Decrypt(ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(e.key[:])
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	ns := gcm.NonceSize()
	if len(ciphertext) < ns {
		return nil, ErrInvalidCiphertext
	}
	nonce, data := ciphertext[:ns], ciphertext[ns:]
	plaintext, err := gcm.Open(nil, nonce, data, nil)
	if err != nil {
		return nil, ErrInvalidCiphertext
	}
	return plaintext, nil
}
